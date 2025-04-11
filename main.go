// Copyright (c) 2025 Grigoriy Efimov
//
// Licensed under the MIT License. See LICENSE file in the project root for details.

package main

import (
	"image"
	"image/color"
	"log"
	"strings"
	"time"

	ink "github.com/CatInBeard/inkview"
)

const defaultFontSize = 14

type TerminalApp struct {
	font               *ink.Font
	inputText          string
	outputText         string
	terminalInputChan  chan string
	terminalOutputChan chan string
	terminalErrorChan  chan string
	fontH              int
	fontW              int
	topTextBoxPosition int
	shouldUpdateScreen bool
}

func (a *TerminalApp) Init() error {
	ink.ClearScreen()

	a.font = ink.OpenFont(ink.DefaultFontMono, a.fontH, true)
	a.font.SetActive(color.RGBA{0, 0, 0, 255})
	a.fontW = ink.CharWidth('a')
	a.topTextBoxPosition = a.fontH

	go func() {
		time.Sleep(5 * time.Second)
		a.shouldUpdateScreen = true
		a.Draw()
		ink.Repaint()
		a.RunCommand("echo \"Welcome to terminal app\"")
	}()

	go term(a.terminalInputChan, a.terminalOutputChan, a.terminalErrorChan)
	go a.HandleTerminalOutput()
	go a.HandleTerminalError()

	ink.SetMessageDelay(time.Second * 5)
	ink.Warningf("Welcome to terminal app", "This application is provided \"as is\" under the MIT license. The source code is available at https://github.com/catInBeard/pb-terminal. Using this terminal emulator application can pose risks to your system and data. Since it emulates a terminal, it can potentially execute commands that may harm your system or compromise your data. You should exercise extreme caution when using this application, especially when executing commands or scripts from untrusted sources. By using this application, you acknowledge that you understand these risks and release the developers from any liability for damages or losses resulting from its use. Proceed with caution and at your own risk.")

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
	a.font.SetActive(color.RGBA{0, 0, 0, 255})

	screenSize := ink.ScreenSize()

	maxCharLength := screenSize.X/a.fontW - 6
	maxLineLength := screenSize.Y/a.fontH - 10

	textLines := strings.Split(a.outputText, "\n")
	var newLines []string
	y := a.topTextBoxPosition

	for _, line := range textLines {
		words := strings.Split(line, " ")
		var newLine string
		for _, word := range words {
			if len(newLine)+len(word)+1 <= maxCharLength {
				newLine += word + " "
			} else {
				newLines = append(newLines, strings.TrimSpace(newLine))
				newLine = word + " "
			}
		}
		newLines = append(newLines, strings.TrimSpace(newLine))
	}

	if len(newLines) > maxLineLength {
		newLines = newLines[len(newLines)-maxLineLength:]
	}

	for _, line := range newLines {
		ink.DrawString(image.Point{a.fontW * 3, y}, line)
		y += a.fontH
	}

	ink.PartialUpdate(image.Rectangle{image.Point{0, 0}, screenSize})

}

func (a *TerminalApp) Key(e ink.KeyEvent) bool {
	if e.State == ink.KeyStateDown {
		switch e.Key {
		case ink.KeyBack:
			if len(a.inputText) > 0 {
				a.inputText = a.inputText[:len(a.inputText)-1]
			}
		case ink.KeyOk:
			a.outputText = a.outputText + "\n" + a.inputText
		}
	}
	return true
}

func (a *TerminalApp) Pointer(e ink.PointerEvent) bool {

	if e.State == ink.PointerDown {
		ink.LoadKeyboard()
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
	a.outputText = a.outputText + "\n$ " + s
	a.terminalInputChan <- s
}

func main() {

	app := &TerminalApp{fontH: defaultFontSize}

	app.terminalInputChan = make(chan string, 5)
	app.terminalOutputChan = make(chan string, 5)
	app.terminalErrorChan = make(chan string, 5)

	if err := ink.Run(app); err != nil {
		log.Fatal(err)
	}

}
