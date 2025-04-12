// Copyright (c) 2025 Grigoriy Efimov
//
// Licensed under the MIT License. See LICENSE file in the project root for details.

package main

import (
	"log"

	ink "github.com/CatInBeard/inkview"
)

func main() {

	app := &TerminalApp{fontH: defaultFontSize}

	app.terminalInputChan = make(chan string, 5)
	app.terminalOutputChan = make(chan string, 5)
	app.terminalErrorChan = make(chan string, 5)

	if err := ink.Run(app); err != nil {
		log.Fatal(err)
	}

}
