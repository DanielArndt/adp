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

// Display a welcome message
func displayWelcome() {
	fmt.Println("\nWelcome to Arndt's Data Processor!")
	fmt.Println("\n" + COPYRIGHT + "\n")
	fmt.Println("Choose an option")
	displayOptions()
}

// Print out the available options
// The available options are stored below in a slice []opt
func displayOptions() {
	for i := 0; i < len(opt); i++ {
		fmt.Println(i, ":", opt[i].desc)
	}
}

// The option struct is used for holding the options available to the
// user.
type option struct {
	// Holds a description of what this option will do
	desc string
	// Method describing the actions of the particular option.
	do func()
}

// Here we define all the possible options for the user in the initial state
var opt = map[int]option {
0: {"Exit", exit},
1: {"Label Data Set", interactiveLabelDataSet},
2: {"Build training and test set", interactiveBuildTrainAndTestSet},
3: {"Convert formats", interactiveConvert},
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
	debugMsg("Input as int: %d", inputInt)
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
