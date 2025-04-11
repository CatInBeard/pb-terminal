// Copyright (c) 2025 Grigoriy Efimov
//
// Licensed under the MIT License. See LICENSE file in the project root for details.

package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func inputThread(cmd *exec.Cmd, stdin io.WriteCloser, inputChan chan string) {
	for input := range inputChan {
		_, err := stdin.Write([]byte(input + "\n"))
		if err != nil {
			fmt.Println(err)
			return
		}
	}
	stdin.Close()
}

func outputThread(stdout io.ReadCloser, outputChan chan string) {
	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		outputChan <- scanner.Text()
	}
	close(outputChan)
	stdout.Close()
}

func errorThread(stderr io.ReadCloser, errorChan chan string) {
	scanner := bufio.NewScanner(stderr)
	for scanner.Scan() {
		errorChan <- scanner.Text()
	}
	close(errorChan)
	stderr.Close()
}

func term(inputChan, outputChan, errorChan chan string) {
	path, err := getCurrentShellPath()

	if err != nil || path == "" {
		path = "sh"
	}

	for {
		cmd := exec.Command(path)

		stdout, err := cmd.StdoutPipe()
		if err != nil {
			fmt.Println(err)
			time.Sleep(1 * time.Second)
			continue
		}

		stderr, err := cmd.StderrPipe()
		if err != nil {
			fmt.Println(err)
			time.Sleep(1 * time.Second)
			continue
		}

		stdin, err := cmd.StdinPipe()
		if err != nil {
			fmt.Println(err)
			time.Sleep(1 * time.Second)
			continue
		}

		if err := cmd.Start(); err != nil {
			fmt.Println(err)
			time.Sleep(1 * time.Second)
			continue
		}

		go inputThread(cmd, stdin, inputChan)
		go outputThread(stdout, outputChan)
		go errorThread(stderr, errorChan)

		if err := cmd.Wait(); err != nil {
			fmt.Println(err)
			time.Sleep(1 * time.Second)
			continue
		}

		break
	}
}

func getCurrentShellPath() (string, error) {
	uid := os.Getuid()

	file, err := os.Open("/etc/passwd")
	if err != nil {
		return "Can't read /etc", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Split(line, ":")
		if len(fields) != 7 {
			continue
		}
		uidField, err := strconv.ParseInt(fields[2], 10, 64)
		if err != nil {
			continue
		}
		if int(uidField) == uid {
			return fields[6], nil
		}
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}

	return "", fmt.Errorf("Cant get shell path", uid)
}
