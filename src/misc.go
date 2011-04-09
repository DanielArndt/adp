/* 
 * misc.go
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
	"sort"
	"os"
	"strings"
	"strconv"
)

func debugMsg(v... interface{}) {
	if *verbose {
		fmt.Print("DEBUG: ")
		fmt.Println(v...)
	}
}

// Check if an error has occured
func errCheck(err os.Error) {
	if err != nil {
		debugMsg("Error:", err)
		os.Exit(1) //Throw an error.
	}
}

// The interactives struct is used for describing the interactive tasks
// available to the user.
type interactives struct {
	// Holds a description of what this interactive will do
	desc string
	// Method describing the actions of the particular option.
	do func()
}

// state 3 - Edit features of a dataset
func interactiveFeatureEditor() {
	var err os.Error
	var inputString string
	
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
