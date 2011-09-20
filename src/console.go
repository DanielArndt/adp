/* 
 * console.go
 * 
 * Copyright (C) 2010 Daniel Arndt
 * 
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 *
 * For more information please visit my website at:
 * http://web.cs.dal.ca/~darndt
 *
 * Or the code's repository:
 *
 * http://github.com/danielarndt/adp
 *  
 */

package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

// Create some variables
var (
	Stdin *bufio.Reader // Used by our Scanf function
)

const (
	ppfx = ">>>" // The prompt prefix. This is used at the start of any line
                 // prompting for the user to enter something
)

/*
 * Replaced the built-in fmt.Scanf with a wrapper on a buffered IO reader.
 * This is to stop additional input from being sent to the calling program
 * upon exit.
 * See: http://code.google.com/p/go/issues/detail?id=1359
 */
func Scanf(format string, v ...interface{}) (int, os.Error) {
	var n int
	var err os.Error

	for n, err = fmt.Fscanf(Stdin, format, v...); n < 1; {
		n, err = fmt.Fscanf(Stdin, format, v...)
	}
	return n, err
}

/*
 * Args:
 *   prompt - What to display on the input line, short description of the input
 *            expected. ie. "filename"
 *   format - The format of the displayed message in same style as fmt.Printf
 *   a ...  - Any additional inputs needed as defined by format.
 */
func promptString(prompt string, format string, a ...interface{}) string {
	var (
		err      os.Error
		toReturn string
	)
	if format != "" {
		fmt.Printf(format, a...)
	}
	fmt.Printf("\n%s%s> ", ppfx, prompt)
	_, err = Scanf("%s", &toReturn)
	errCheck(err)
	return toReturn
}

/*
 * Args:
 *   prompt - What to display on the input line, short description of the input
 *            expected. ie. "filename"
 *   format - The format of the displayed message in same style as fmt.Printf
 *   a ...  - Any additional inputs needed as defined by format.
 */
func promptInt(prompt string, format string, a ...interface{}) int {
	var (
		err      os.Error
		toReturn int
	)
	fmt.Printf(format, a...)
	fmt.Printf("\n%s%s> ", ppfx, prompt)
	_, err = Scanf("%d", &toReturn)
	errCheck(err)
	return toReturn
}

func debugMsg(format string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, "DEBUG: ")
	fmt.Fprintf(os.Stderr, format, a...)
	fmt.Fprintln(os.Stderr)
}

// Check if an error has occured
func errCheck(err os.Error) {
	if err != nil {
		log.Fatalln("Error:", err)
	}
}
