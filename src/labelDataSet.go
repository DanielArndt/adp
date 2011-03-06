/* 
 * labelDataSet.go
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
    "strconv"
    "strings"
)

// This uses the quick method where maps are used. This is only for sets
// of rules which have no conflicts and have a simple format "if coloumn
// x = value then label is y.

func quickRules(filepath string) map[int]map[int]string {
	// list to return which contains the parsed rules
	debugMsg("Opening file \"" + filepath + "\"")
	// Open the rule file
	dataFile, err := os.Open(filepath, os.O_RDONLY, 0666)
	errCheck(err)
	defer dataFile.Close()
	// Create a buffered reader for the rule file
	dataReader := bufio.NewReader(dataFile)
	// A map which points features to another map which contains possible values 
	// for that feature.
	// map[featureindex->[value->label]]
	featToValMap := map[int]map[int]string{}
	// Read in the contents
	for line, err := dataReader.ReadString('\n'); // read line by line
	err == nil;                  // loop until end of file or error
	line, err = dataReader.ReadString('\n') {
		// Trim newline from end
		line = strings.TrimRight(line, "\n")
		if strings.HasPrefix(line, "#") {
			// Ignore comments
			debugMsg("Skipping line due to comment:", line)
		} else {
			// Split by fields
			fields := strings.Fields(line)
			if len(fields) == 3 {
				// Deal with comma seperated feature indexes
				features := strings.Split(fields[0], ",", -1)
				// Make a struct for each feature index
				for i := 0; i < len(features); i++ {
					debugMsg("Making label rule:")
					debugMsg("if feature [", features[i], "] == [", fields[1], "] {")
					debugMsg("\tlabel =", fields[2])
					debugMsg("}")
					// Read in some values
					featureIndex, err := strconv.Atoi(features[i])
					errCheck(err)
					value, err := strconv.Atoi(fields[1])
					errCheck(err)
					label := fields[2]
					errCheck(err)
					_, exists := featToValMap[featureIndex]
					if exists {
						featToValMap[featureIndex][value] = label
					} else {
						featToValMap[featureIndex] = map[int]string{value: label}
					}
				}
			} else {
				debugMsg("Malformed line: \"" + line + "\"")
				debugMsg("Length:", len(line))
				debugMsg("Err:",err)
			}
		}
	}
	return featToValMap
}

//state 1 - Label a data set
func interactiveLabelDataSet() {
	var (
		featToValMap map[int]map[int]string
		err          os.Error
		inputInt     int
		inputString  string
	)
	// Load in the rules
	fmt.Println("Which rule method would you like to use?")
	fmt.Println("0 : quick rules")
	fmt.Println("1 : extended rules (unfinished)")
	fmt.Print("> ")
	_, err = Scanf("%d", &inputInt)
	errCheck(err)
	switch inputInt {
	case 0:
		featToValMap = quickRules("label.rules")
	}

	// Read out the maps stored for each feature
	for k, v := range featToValMap {
		debugMsg("column:", k, v)
	}
	// Begin labeling the data set
	fmt.Println("Label a data set")
	fmt.Println("Please enter the location of the file which contains the",
		"dataset")
	fmt.Print("file name> ")
	// Receive file name of data set
	_, err = Scanf("%s", &inputString)
	errCheck(err)
	debugMsg("Opening file:", inputString)
	// Open the file for input and create a buffered reader for the file
	dataFile, err := os.Open(inputString, os.O_RDONLY, 0666)
	errCheck(err)
	// We do not need this file after, so close it upon leaving this method
	defer dataFile.Close()
	dataReader := bufio.NewReader(dataFile)
	// Open a file for the labeled training set
	debugMsg("Opening file:", dataFile.Name()+".labeled")
	labeledFile, err := os.Open(
		dataFile.Name()+".labeled",
		os.O_CREATE+os.O_WRONLY+os.O_TRUNC,
		0666)
	errCheck(err)
	debugMsg("Writing to file:", dataFile.Name()+".labeled")
	debugMsg("Labeling... this may take a while")
	// We do not need this file after, so close it upon leaving this method
	defer labeledFile.Close()
	// Create a variable for the line read, and the label assigned
	var line, label string
	// Loop over each line of the file
	for line, err = dataReader.ReadString('\n'); // read line by line
	err == nil;                                  // stop on error or end of file
	line, err = dataReader.ReadString('\n') {
		line = strings.TrimRight(line, "\n")
		// Split the line into it's feature values
		feature := strings.Split(line, ",", -1)
		// FIXME: fix the way we deal with malformed lines
		if len(feature) < 5 {
			debugMsg("Skipping line due to abnormal formation")
			break
		}
		//Find the rule that satisfies the current individual, if any.
		for ruleFeature, ruleValMap := range featToValMap {
			instanceFeatVal, err := strconv.Atoi(feature[ruleFeature])
			errCheck(err)
			// Try to find the corresponding value in the map for the current
			// feature index.
			valLabel, exists := ruleValMap[instanceFeatVal]
			if exists {
				label = valLabel
				break
			} else if label == "" {
				label = "OTHER"
			}
		}
		// Write labeled line to labeled file
		_, err = labeledFile.WriteString(line + "," + label + "\n")
		errCheck(err)
		label = ""
	}
}

