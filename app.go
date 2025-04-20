// Copyright (c) 2025 Grigoriy Efimov
//
// Licensed under the MIT License. See LICENSE file in the project root for details.

package main

import (
	_ "embed"
	"image"
	"image/color"
	"os"
	"strings"
	"time"
	"unicode/utf8"

	ink "github.com/CatInBeard/inkview"
)

const defaultFontSize = 14

type TerminalApp struct {
	font               *ink.Font
	outputText         string
	terminalInputChan  chan string
	terminalOutputChan chan string
	terminalErrorChan  chan string
	fontH              int
	fontW              int
	topTextBoxPosition int
	shouldUpdateScreen bool
	lang               string
}

func (a *TerminalApp) Init() error {
	ink.ClearScreen()
	ink.DrawTopPanel()

	a.lang = ink.GetCurrentLang()

	createCustomKeyboardIfNotExists()
	ink.SetKeyboardHandler(a.terminalKeyboardHandler)

	a.font = ink.OpenFont(ink.DefaultFontMono, a.fontH, true)
	a.font.SetActive(color.RGBA{0, 0, 0, 255})
	a.fontW = ink.CharWidth('a') // Work only for monospace font
	a.topTextBoxPosition = a.fontH

	go term(a.terminalInputChan, a.terminalOutputChan, a.terminalErrorChan)
	go a.HandleTerminalOutput()
	go a.HandleTerminalError()

	ink.SetMessageDelay(time.Second * 10)

	ink.Warningf(a.GetTranslation("warning_title"), a.GetTranslation("warning_text"))

	a.shouldUpdateScreen = true
	a.RunCommand("echo \"" + a.GetTranslation("hello_cmd_text") + "\"")
	a.Draw()
	ink.Repaint()

	return nil
}

func (a *TerminalApp) Close() error {
	return nil
}

func (a *TerminalApp) Draw() {

	if !a.shouldUpdateScreen {
		return
	}

	ink.ClearScreen()
	ink.DrawTopPanel()
	a.font.SetActive(color.RGBA{0, 0, 0, 255})

	screenSize := ink.ScreenSize()

	maxCharLength := screenSize.X/a.fontW - 6
	maxLineLength := screenSize.Y/a.fontH - 10

	textLines := strings.Split(a.outputText, "\n")

	y := a.topTextBoxPosition

	if len(textLines) > maxLineLength {
		textLines = textLines[len(textLines)-maxLineLength:]
	}

	for _, line := range textLines {
		splittedLines := splitText(line, maxCharLength)
		for _, line := range splittedLines {
			ink.DrawString(image.Point{a.fontW * 3, y}, line)
			y += a.fontH
		}
	}

	ink.PartialUpdate(image.Rectangle{image.Point{0, 0}, screenSize})

}

func (a *TerminalApp) Key(e ink.KeyEvent) bool {
	return true
}

func (a *TerminalApp) Pointer(e ink.PointerEvent) bool {

	if e.State == ink.PointerDown {
		ink.Repaint()
		a.shouldUpdateScreen = false
		a.invokeKeybaord("Enter your command")
	}
	return true
}

func (a *TerminalApp) Touch(e ink.TouchEvent) bool {
	return true
}

func (a *TerminalApp) Orientation(o ink.Orientation) bool {
	return true
}

func (a *TerminalApp) HandleTerminalOutput() {
	for output := range a.terminalOutputChan {
		a.outputText = a.outputText + "\n" + output
		ink.Repaint()
	}
}

func (a *TerminalApp) HandleTerminalError() {
	for output := range a.terminalErrorChan {
		a.outputText = a.outputText + "\n" + output
		ink.Repaint()
	}
}

func (a *TerminalApp) RunCommand(s string) {
	trimText := strings.TrimSpace(s)
	if trimText == "clear" {
		a.outputText = ""
		ink.Repaint()
		return
	} else if trimText == "exit" {
		ink.Exit()
	}
	a.outputText = a.outputText + "\n$ " + s
	a.terminalInputChan <- s
}

func (a *TerminalApp) terminalKeyboardHandler(text string) {
	a.shouldUpdateScreen = true
	a.RunCommand(text)
}

func (a *TerminalApp) invokeKeybaord(label string) {
	ink.OpenCustomKeyboard(ink.UserKbdPath+"/devkeyboard.kbd", label, 1024)
}

func splitText(inputStr string, maxLen int) []string {
	var result []string
	var tempStr string

	for _, char := range inputStr {
		if char == '\n' {
			if tempStr != "" {
				result = append(result, tempStr)
				tempStr = ""
			}
		} else {
			tempStr += string(char)
			if utf8.RuneCount([]byte(tempStr)) >= maxLen {
				result = append(result, tempStr)
				tempStr = ""
			}
		}
	}

	if tempStr != "" {
		result = append(result, tempStr)
	}

	return result
}

func (a *TerminalApp) GetTranslation(key string) string {
	return GetTranslation(a.lang, key)
}

//go:embed devkeyboard.kbd
var testKbd []byte

func createCustomKeyboardIfNotExists() error {
	filePath := ink.UserKbdPath + "/devkeyboard.kbd"
	_, err := os.Stat(filePath)
	if err != nil {
		err = os.WriteFile(filePath, testKbd, 0644)
		if err != nil {
			return err
		}
	}
	return nil
}
