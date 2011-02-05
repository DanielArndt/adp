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
 * http://code.google.com/p/adp
 *  
 */

package main

import (
	"fmt"
	"bufio"
	"log"
	"sort"
	"os"
	"syscall"
	"strings"
	"strconv"
)

// Create some constants
const (
	TAB       = "\t"
	COPYRIGHT = "Copyright (C) 2010 Daniel Arndt\n" +
		"This program comes with ABSOLUTELY NO WARRANTY;\n" +
		"This is free software, and you are welcome to redistribute it\n" +
		"under certain conditions; please see the COPYING file which\n" +
		"you should have received with this program. For more \n" +
		"information, please visit: \n" +
		"http://web.cs.dal.ca/~darndt"
	LOGSLP = 50000
)

// Create some reusable variables
var (
	err         os.Error      // The most recent error
	inputInt    int           // Most recent user input integer
	inputString string        // Most recent user input string
	Stdin       *bufio.Reader // Used by the Scanf function
)

// Replaced the built-in fmt.Scanf with a wrapper on a buffered IO reader.
// This is to stop additional input from being sent to the calling program
// upon exit.
// See: http://code.google.com/p/go/issues/detail?id=1359
func Scanf(f string, v ...interface{}) (int, os.Error) {
	var n int
	for n, err = fmt.Fscanf(Stdin, f, v...); n < 1; {
		n, err = fmt.Fscanf(Stdin, f, v...)
	}
	return n, err
}

func debugMsg(v ...interface{}) {
	fmt.Print("DEBUG: ")
	fmt.Println(v...)
}

// Check if an error has occured
func errCheck(err os.Error) {
	if err != nil {
		log.Fatalln("Error:", err)
	}
}

// Sleep used to allow the log to catch up.
func sleep(i int64) {
	syscall.Sleep(i)
}

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
var opt = map[int]option{
	0: {"Exit", exit},
	1: {"Label Data Set", interactiveLabelDataSet},
	2: {"Build trainind and test set", interactiveBuildTrainAndTestSet},
}


func init() {
	Stdin = bufio.NewReader(os.Stdin)
}

func main() {
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



// state 3 - Edit features of a dataset
func interactiveFeatureEditor() {
	// STEP 1:
	// Receive the name of the data set to work on and open the file
	fmt.Println("\nEdit the feature set")
	fmt.Println("--------------------")
	fmt.Print("file name> ")
	// Receive file name of data file
	_, err = Scanf("%s", &inputString)
	errCheck(err)
	debugMsg("Opening file:", inputString)
	// Open the file for input and create a buffered reader for the file
	dataFile, err := os.Open(inputString, os.O_RDONLY, 0666)
	errCheck(err)
	// We do not need this file after, so close it upon leaving this method
	defer dataFile.Close()

	// STEP 2:
	// Allow the user to input commands, and interpret them

	fmt.Print("command> ")
	// Receive file name of data file
	var cmd, cmdpar string
	_, err = Scanf("%s %s", &cmd, &cmdpar)
	errCheck(err)
	// Split the parameters by comma to get each individual value
	params := strings.Split(cmdpar, ",", -1)
	debugMsg("CMD:", cmd)
	debugMsg("PAR:", cmdpar)
	// List holding the items on which to act upon
	var actList []int
	var splitParams []string
	// Integer form of the parameter
	var intParam, intParamEnd int
	for i := 0; i < len(params); i++ {
		// Deal with each individual parameter
		debugMsg("PAR", i, "::", params[i])
		splitParams = strings.Split(params[i], "-", -1)
		if len(splitParams) == 1 {
			intParam, err = strconv.Atoi(splitParams[0])
			actList = append(actList, intParam)
		} else if len(splitParams) == 2 {
			// The start of the range to act upon
			intParam, err = strconv.Atoi(splitParams[0])
			// The end of the range to act upon
			intParamEnd, err = strconv.Atoi(splitParams[1])
			for j := intParam; j <= intParamEnd; j++ {
				actList = append(actList, j)
			}
		} else {
			log.Fatalln("Error splitting parameters. len(splitParams):",
				len(splitParams))
		}
	}
	sort.SortInts(actList)
	for i := 0; i < len(actList); i++ {
		// Print out each element in the action list
		debugMsg("actList", i, "::", actList[i])
	}
	sleep(LOGSLP)
	dataReader := bufio.NewReader(dataFile)
	var line string
	// Read from file loop
	for line, err = dataReader.ReadString('\n'); // read line by line
	err == nil;                                  // stop on error
	line, err = dataReader.ReadString('\n') {
		feature := strings.Split(line, ",", -1)
		i := 0
		for i < len(feature) {
			//diff := actList[i] - i
			//only output those not in actList
		}
	}
}
