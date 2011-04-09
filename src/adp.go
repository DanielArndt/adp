/* 
 * adp.go
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
	"flag"
	"fmt"
	"bufio"
	"os"
)

// Create some constants
const (
	TAB       = "\t"
	COPYRIGHT = "Copyright (C) 2010 Daniel Arndt\n" +
		"This program comes with ABSOLUTELY NO WARRANTY; " +
		"This is free software, and you are welcome to redistribute it " +
		"under certain conditions; please see the COPYING file which " +
		"you should have received with this program. For more " +
		"information, please visit: \n" +
		"http://web.cs.dal.ca/~darndt"
)

// Create some variables
var (
	Stdin       *bufio.Reader // Used by our Scanf function
	verbose = flag.Bool("v", false, "Print extra, detailed output");
)


// Replaced the built-in fmt.Scanf with a wrapper on a buffered IO reader.
// This is to stop additional input from being sent to the calling program
// upon exit.
// See: http://code.google.com/p/go/issues/detail?id=1359
func Scanf(f string, v ...interface{}) (int, os.Error) {
	var n int
	var err os.Error
	for n, err = fmt.Fscanf(Stdin, f, v...); n < 1; {
		n, err = fmt.Fscanf(Stdin, f, v...)
	}
	return n, err
}


// Display a welcome message
func displayWelcome() {
	fmt.Println("\nWelcome to Arndt's Data Processor!")
	fmt.Println("\n" + COPYRIGHT + "\n")
	fmt.Println("Choose an option")
	displayOptions()
}

// Print out the available tasks. The available tasks are stored below
// in a slice []opt
func displayOptions() {
	for i := 0; i < len(opt); i++ {
		fmt.Println(i, ":", opt[i].desc)
	}
}

// Here we define all the possible tasks for the user in the initial
// state.
var opt = map[int]interactives{
	0: {"Exit the program", exit},
	1: {"Label Data Set", interactiveLabelDataSet},
	2: {"Build training and test set", interactiveBuildTrainAndTestSet},
}


func init() {
	Stdin = bufio.NewReader(os.Stdin)
}

func main() {
	var err os.Error
	var	inputInt int
	
	// MAIN ------------------------- 
	displayWelcome()
	fmt.Printf("> ")
	// Read in an int as the state to go to.
	_, err = Scanf("%d", &inputInt)
	errCheck(err)
	debugMsg("Input as int:", inputInt)
	if inputInt >= 0 && inputInt < len(opt) {
		//Execute the method in the opt struct's "do" field
		opt[inputInt].do()
	} else {
		fmt.Println("Invalid input")
		fmt.Println()
	}
	// ------------------------------
}

//state 0
func exit() {
	fmt.Println("Exiting")
	os.Exit(0)
}
