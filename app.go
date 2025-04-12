// Copyright (c) 2025 Grigoriy Efimov
//
// Licensed under the MIT License. See LICENSE file in the project root for details.

package main

import (
	"image"
	"image/color"
	"strings"
	"time"

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
}

func (a *TerminalApp) Init() error {
	ink.ClearScreen()

	ink.SetKeyboardHandler(a.terminalKeyboardHandler)

	a.font = ink.OpenFont(ink.DefaultFontMono, a.fontH, true)
	a.font.SetActive(color.RGBA{0, 0, 0, 255})
	a.fontW = ink.CharWidth('a')
	a.topTextBoxPosition = a.fontH

	go term(a.terminalInputChan, a.terminalOutputChan, a.terminalErrorChan)
	go a.HandleTerminalOutput()
	go a.HandleTerminalError()

	ink.SetMessageDelay(time.Second * 5)

	ink.Warningf("Welcome to terminal app", "This application is provided \"as is\" under the MIT license. The source code is available at https://github.com/catInBeard/pb-terminal. Using this terminal emulator application can pose risks to your system and data. Since it emulates a terminal, it can potentially execute commands that may harm your system or compromise your data. You should exercise extreme caution when using this application, especially when executing commands or scripts from untrusted sources. By using this application, you acknowledge that you understand these risks and release the developers from any liability for damages or losses resulting from its use. Proceed with caution and at your own risk.")

	go func() {
		time.Sleep(3 * time.Second)
		a.shouldUpdateScreen = true
		a.RunCommand("echo \"Welcome to terminal app!\"")
		a.Draw()
		ink.Repaint()
	}()

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
	a.outputText = a.outputText + "\n$ " + s
	a.terminalInputChan <- s
}

func (a *TerminalApp) terminalKeyboardHandler(text string) {
	a.shouldUpdateScreen = true
	a.RunCommand(text)
}

func (a *TerminalApp) invokeKeybaord(label string) {
	ink.OpenKeyboard(label, 1024)
}
