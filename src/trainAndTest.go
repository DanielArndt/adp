/* 
 * trainAndTest.go
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
    "os"
    "rand"
    "sort"
    "strings"
)

// state 2 - Build and train test set
func interactiveBuildTrainAndTestSet() {
	var (
		err os.Error
		inputString string
		inputInt int
	)
	// STEP 1:
	// Begin building training and test set
	fmt.Println("Building train and test set")
	inputString = promptString("filename", "What file would you like to split?")
	debugMsg("Opening file: %s", inputString)
	// Open the file for reading
	dataFile, err := os.Open(inputString)
	errCheck(err)
	// We do not need this file after, so close it upon leaving this method
	defer dataFile.Close()
	// Create a buffered reader for the file
	dataReader := bufio.NewReader(dataFile)
	var line string

	// STEP 2:
	// Create a map for storing the temporary files

	tempFileMap := map[string]*os.File{}
	countMap := map[string]int{}
	var exists bool       // For checking if element exists
	var tempFile *os.File // Place holder for the temporary file
	// which is to be put in the map.
	for line, err = dataReader.ReadString('\n'); // read line by line
	err == nil;                                  // stop on error
	line, err = dataReader.ReadString('\n') {
		// This is where we split up each label into its own file.
		line = strings.Trim(line, "\n")
		feature := strings.Split(line, ",")
		label := feature[len(feature)-1]
		tempFile, exists = tempFileMap[label]
		countMap[label]++
		if exists {
			// Write to the file
			_, err = tempFile.WriteString(line + "\n")
			errCheck(err)
		} else {
			// Create the file and write the line
			tempFileName := dataFile.Name() + "." + label + ".tmp"
			debugMsg("Creating temporary file: %s", tempFileName)
			tempFile, err := os.OpenFile(
				tempFileName,
				os.O_CREATE+os.O_WRONLY+os.O_TRUNC,
				0666)
			tempFileMap[label] = tempFile
			defer tempFile.Close()
			defer os.Remove(tempFileName)
			_, err = tempFile.WriteString(line + "\n")
			errCheck(err)
		}
	}
	// Close and re-open the files as readable
	debugMsg("Closing all temporary files for writing. Re-opening as read-only.")
	for k, v := range tempFileMap {
		fileName := v.Name()
		v.Close()
		tempFileMap[k], err = os.Open(fileName)
		errCheck(err)
		// We do not need this file after, so close it upon leaving this method
		defer tempFileMap[k].Close()
	}

	// STEP 3: 
	// Receive the number of each label (class) we'd like to add to the training
	// set

	// Hold the amount of each label we'd like in the training set in a map
	trainCountMap := map[string]int{}
	fmt.Println("Please enter the number of each type of label you'd", 
		"like in the training set.")
	// Ask user how much of each label they want and put it in a map 
	// trainCountMap
	for k, v := range countMap {
		inputInt = promptInt(k, "label: %s max: %d", k, v) 
		trainCountMap[k] = inputInt
	}
	debugMsg("Creating: %s", dataFile.Name()+".train")
	// Open a file for writing training data
	trainFile, err := os.OpenFile(
		dataFile.Name()+".train",
		os.O_CREATE+os.O_WRONLY+os.O_TRUNC,
		0666)
	errCheck(err)
	// We do not need this file after, so close it upon leaving this method
	defer trainFile.Close()
	// STEP 4:
	// Read the correct amount of each label in

	for k, v := range trainCountMap {
		debugMsg("label: %s count: %d", k, v)
		dataReader := bufio.NewReader(tempFileMap[k])
		// Open a file for writing testing data
		testFile, err := os.OpenFile(
			dataFile.Name()+"."+k+".test",
			os.O_CREATE+os.O_WRONLY+os.O_TRUNC,
			0666)
		errCheck(err)
		// We do not need this file after, so close it upon leaving this method
		defer testFile.Close()

		if v > 0 {
			// Generate a random permuation
			var randomized sort.IntSlice
			randomized = rand.Perm(countMap[k])
			// use a slice the first /v/ of them
			randomized = randomized[0:v]
			// sort the ints so that as we iterate through each instance we can
			// easily find the next one we need to export
			randomized.Sort()
			// Read through the file, writing the included instances to
			// .train and the others to .test
			lineCount := 0
			if len(randomized) > 0 {
				for line, err = dataReader.ReadString('\n'); // read line by line
				err == nil;                                  // stop on error
				line, err = dataReader.ReadString('\n') {
					if lineCount == randomized[0] {
						_, err = trainFile.WriteString(line)
						errCheck(err)
						if len(randomized) > 1 {
							randomized = randomized[1:len(randomized)]
						} else {
							randomized[0] = -1 // skip the rest
						}
					} else {
						_, err = testFile.WriteString(line)
						errCheck(err)
					}
					lineCount++
				}
			} else {

			}
		} else {
			//TODO: Add a handler for -1
			// None of the label were requested in the training set, so dump to test
			for line, err = dataReader.ReadString('\n'); // read line by line
			err == nil;                                  // stop on error
			line, err = dataReader.ReadString('\n') {
				_, err = testFile.WriteString(line)
				errCheck(err)
			}
		}
		testFile.Close()
	}
	fmt.Println()
}
