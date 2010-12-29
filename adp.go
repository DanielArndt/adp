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
	"rand"
	"os"
	"syscall"
	"strings"
	"strconv"
)

// Create some constants
const (
	TAB       = "\t"
	COPYRIGHT = "Copyright (C) 2010 Daniel Arndt" + "\n" +
		"This program comes with ABSOLUTELY NO WARRANTY;" + "\n" +
		"This is free software, and you are welcome to redistribute it" + "\n" +
		"under certain conditions; please see the COPYING file which" + "\n" +
		"you should have received with this program. For more " + "\n" +
		"information, please visit: " + "\n" +
		"http://web.cs.dal.ca/~darndt"
	LOGSLP = 2
)

// Create some reusable variables for holding information
var (
	err         os.Error // The most recent error
	inputInt    int      // Most recent user input integer
	inputString string   // Most recent user input string
	Stdin       *bufio.Reader
)

// Replaced the built-in fmt.Scanf with a wrapper on a buffered IO reader.
// This is to stop additional input from being sent to the calling program
// upon exit.
// See: http://code.google.com/p/go/issues/detail?id=1359
func Scanf(f string, v ...interface{}) (int, os.Error) {
	var n int
	for n,err = fmt.Fscanf(Stdin, f, v...); n < 1; {
		n,err = fmt.Fscanf(Stdin, f, v...)
	}
	return n,err
}

// Check if an error has occured
func errCheck(err os.Error) {
	if err != nil { 
		log.Exitln("Error:",err) 
	} 
}

// Sleep used to allow the log to catch up.
func sleep(i int64) {
	syscall.Sleep(i)
}

// Display a welcome message
func displayWelcome() {
	fmt.Println("\n" + "Welcome to Arndt's Data Processor!")
	fmt.Println("\n" + COPYRIGHT + "\n") 
	fmt.Println("Choose an option")
	displayOptions()
}

// Print out the available options
// The available options are stored below in a slice []opt
func displayOptions() {
	for i := 0; i < len(opt); i++ {
		fmt.Println(opt[i])
	}
}

// The option struct is used for holding the options available to the
// user.
// TODO: Re-write the options to use a map instead of this struct and the
//       following slice.
type option struct {
	// Holds the id of the option, what the user enters to preform the action
	id int
	// Holds a description of what this option will do
	desc string
	// Method describing the actions of the particular option.
	do   func()
}

// The available options
var opt []option = []option {
	{id: 0, desc: "Exit ADP", do: exit},
	{id: 1, desc: "Label a data set", do: interactiveLabelDataSet},
	{id: 2, desc: "Build training and test set", do: interactiveBuildTrainAndTestSet},
	{id: 3, desc: "Edit the feature set (UNFINISHED)", do: interactiveFeatureEditor},
}

// Loads the file referred to by filepath and parses it into rules used
// to label a data set. [fig 1]
func loadRules(filepath string) map[int]map[int]string {
	// list to return which contains the parsed rules
	log.Println("Opening file \"" + filepath + "\"")
	sleep(LOGSLP)
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
	err == nil;                                   // loop until end of file or error
	line, err = dataReader.ReadString('\n') {
		// Trim newline from end
		line = strings.TrimRight(line, "\n")
		if strings.HasPrefix(line, "#") {
			// Ignore comments
			log.Println("Skipping line due to comment:", line)
			sleep(LOGSLP)
		} else {
			// Split by fields
			field := strings.Fields(line)
			// Deal with comma seperated feature indexes
			feature := strings.Split(field[0], ",", -1)
			// Make a struct for each feature index
			for i := 0; i < len(feature); i++ {
				log.Println("Making label rule:")
				log.Println("if", feature[i], "==", field[1], "{")
				log.Println("\tlabel =", field[2])
				log.Println("}")
				sleep(LOGSLP)
				// Read in some values
				featureIndex, err := strconv.Atoi(feature[i])
				errCheck(err)
				value, err := strconv.Atoi(field[1])
				errCheck(err)
				label := field[2]
				errCheck(err)
				_, exists := featToValMap[featureIndex]
				if exists {
					featToValMap[featureIndex][value] = label
				} else {
					featToValMap[featureIndex] = map[int]string{value:label}
				}
			}
		}
	}
	return featToValMap
}

func main() {
	// MAIN ------------------------- 
	displayWelcome()
	fmt.Printf("> ")
	// Read int, throw away count
	_, err = Scanf("%d", &inputInt)
	errCheck(err)
	// Write the input read
	log.Printf("Input as int: %d\n", inputInt)
	sleep(LOGSLP)
	if inputInt >= 0 && inputInt < len(opt) {
		//Execute the method in the structs "do" field
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

//state 1 - Label a data set
func interactiveLabelDataSet() {
	// Load in the rules
	featToValMap := loadRules("label.rules")
	// Read out the maps stored for each feature
	for k,v := range featToValMap {
		log.Println("key:", k, "val:", v)
	}
	sleep(LOGSLP)
	// Begin labeling the data set
	fmt.Println("Label a data set")
	fmt.Println("Please enter the location of the file which contains the",
		"dataset")
	fmt.Print("file name> ")
	// Receive file name of data set
	_, err = Scanf("%s", &inputString)
	errCheck(err)
	log.Println("Opening file:", inputString)
	sleep(LOGSLP)
	// Open the file for input and create a buffered reader for the file
	dataFile, err := os.Open(inputString, os.O_RDONLY, 0666)
	errCheck(err)
	// We do not need this file after, so close it upon leaving this method
	defer dataFile.Close()
	dataReader := bufio.NewReader(dataFile)
	// Open a file for the labeled training set
	log.Println("Opening file:", dataFile.Name() + ".labeled")
	sleep(LOGSLP)
	labeledFile, err := os.Open(
		dataFile.Name()+".labeled",
		os.O_CREATE+os.O_WRONLY+os.O_TRUNC,
		0666)
	errCheck(err)
	log.Println("Writing to file:", dataFile.Name() + ".labeled")
	log.Println("This may take a while")
	sleep(LOGSLP)
	// We do not need this file after, so close it upon leaving this method
	defer labeledFile.Close()
	// Create a variable for the line read, and the label assigned
	var line, label string
	// Loop over each line of the file
	for line, err = dataReader.ReadString('\n'); // read line by line
	err == nil;                                  // stop on error
	line, err = dataReader.ReadString('\n') {
		line = strings.TrimRight(line, "\n")
		// Split the line into it's feature values
		feature := strings.Split(line, ",", -1)
		// FIXME: fix the way we deal with malformed lines
		if len(feature) < 5 {
			log.Println("Skipping line")
			break
		}
		//Find the rule that satisfies the current individual, if any.
		for ruleFeature,ruleValMap := range featToValMap {
			instanceFeatVal, err := strconv.Atoi(feature[ruleFeature])
			errCheck(err)
			// Try to find the corresponding value in the map for the current
			// feature index.
			valLabel, exists := ruleValMap[instanceFeatVal]
			if exists {
				label = valLabel; break
			} else if (label == "") {
				label = "OTHER"
			}
		}
		_, err = labeledFile.WriteString(line + "," + label + "\n")
		label = ""
		sleep(LOGSLP)
		// Print labeled line to labeled file
		errCheck(err)
	}
}

// state 2 - Build and train test set
func interactiveBuildTrainAndTestSet() {
	// STEP 1:
	// Begin building training and test set

	fmt.Println("Building train and test set")
	fmt.Println("What file would you like to split?")
	fmt.Print("file name> ")
	// Receive file name of data file
	_, err = Scanf("%s", &inputString)
	errCheck(err)
	log.Println("Opening file:", inputString)
	sleep(LOGSLP)
	// Open the file for reading
	dataFile, err := os.Open(inputString, os.O_RDONLY, 0666)
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
	var exists bool // For checking if element exists
	var tempFile *os.File // Place holder for the temporary file
	                      // which is to be put in the map.
	for line, err = dataReader.ReadString('\n'); // read line by line
	err == nil;                                  // stop on error
	line, err = dataReader.ReadString('\n') {
		// This is where we split up each label into its own file.
		line = strings.Trim(line, "\n")
		feature := strings.Split(line, ",", -1)
		label := feature[len(feature)-1]
		tempFile, exists = tempFileMap[label]
		countMap[label]++
		if exists {
			// Write to the file
			_, err = tempFile.WriteString(line + "\n")
			errCheck(err)
		} else {
			// Create the file and write the line
			log.Println("Creating temporary file:", dataFile.Name()+"."+label+".tmp") 
			tempFile, err := os.Open(
				dataFile.Name()+"."+label+".tmp",
				os.O_CREATE+os.O_WRONLY+os.O_TRUNC,
				0666)
			tempFileMap[label] = tempFile
			defer tempFile.Close()
			_, err = tempFile.WriteString(line + "\n")
			errCheck(err)
		}
	}
	// Close and re-open the files as readable
	log.Println("Closing all temporary files for writing. Re-opening as")
	log.Println("read-only")
	sleep(1)
	for k, v := range tempFileMap {
		fileName := v.Name()
		v.Close()
		tempFileMap[k], err = os.Open(fileName, os.O_RDONLY, 0666)
		errCheck(err)
		// We do not need this file after, so close it upon leaving this method
		defer tempFileMap[k].Close()
	}

	// STEP 3: 
	// Receive the number of each label (class) we'd like to add to the training
	// set

	// Hold the amount of each label we'd like in the training set in a map
	trainCountMap := map[string]int{}
	fmt.Println("Please enter the number of each type of label you'd")
	fmt.Println("like in the training set. Enter -1 for no bias")
	// Ask user how much of each label they want and put it in a map trainCountMap
	for k, v := range countMap {
		fmt.Println("label:", k, "max:", v)
		fmt.Printf("> ")
		Scanf("%d", &inputInt)
		log.Println("inputInt:", inputInt)
		sleep(LOGSLP)
		trainCountMap[k] = inputInt
	}
	log.Println("Creating:", dataFile.Name()+".train")
	log.Println("Creating:", dataFile.Name()+".test")
	sleep(LOGSLP)
	// Open a file for writing training data
	trainFile, err := os.Open(
		dataFile.Name()+".train",
		os.O_CREATE+os.O_WRONLY+os.O_TRUNC,
		0666)
	errCheck(err)
	// We do not need this file after, so close it upon leaving this method
	defer trainFile.Close()
	// Open a file for writing testing data
	testFile, err := os.Open(
		dataFile.Name()+".test",
		os.O_CREATE+os.O_WRONLY+os.O_TRUNC,
		0666)
	errCheck(err)
	// We do not need this file after, so close it upon leaving this method
	defer testFile.Close()

	// STEP 4:
	// Read the correct amount of each label in

	for k, v := range trainCountMap {
		log.Println("label:", k, "count:", v)
		sleep(1)
		// Generate a random permuation
		rand := rand.Perm(countMap[k])
		// cut and use the first v of them
		rand = rand[0:v]
		// sort the ints
		sort.SortInts(rand)
		log.Println("rand:", rand)
		// Read through the file, writing the included instances to
		// .train and the others to .test
		dataReader := bufio.NewReader(tempFileMap[k])
		lineCount := 0
		if len(rand) != 0 {
			for line, err = dataReader.ReadString('\n'); // read line by line
			err == nil;                                  // stop on error
			line, err = dataReader.ReadString('\n') {
				if lineCount == rand[0] {
					_, err = trainFile.WriteString(line)
					errCheck(err)
					if len(rand) > 1 {
						rand = rand[1:len(rand)]
					} else {
						rand[0] = -1 // skip the rest
					}
				} else {
					_, err = testFile.WriteString(line)
					errCheck(err)
				}
				lineCount++
			}
		} else {
			// None of the label were requested in the training set, so dump to test
			for line, err = dataReader.ReadString('\n'); // read line by line
			err == nil;                                  // stop on error
			line, err = dataReader.ReadString('\n') {
				_, err = testFile.WriteString(line)
				errCheck(err)
			}
		}
	}

	// STEP 5:
	// Cleanup: Remove the tempFiles

	for _, v := range tempFileMap {
		fileName := v.Name()
		log.Println("Deleting temporary file:", fileName)
		v.Close()
		os.Remove(fileName)
	}
	fmt.Println()
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
	log.Println("Opening file:", inputString)
	sleep(LOGSLP)
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
	log.Println("CMD:", cmd)
	log.Println("PAR:", cmdpar)
	// List holding the items on which to act upon
	var actList []int
	var splitParams []string
	// Integer form of the parameter
	var intParam, intParamEnd int
	for i := 0; i < len(params); i++ {
		// Deal with each individual parameter
		log.Println("PAR", i, "::", params[i])
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
			log.Exitln("Error splitting parameters. len(splitParams):", 
				len(splitParams))
		}
	}
	sleep(LOGSLP)
	sort.SortInts(actList)
	for i := 0; i < len(actList); i++ {
		// Print out each element in the action list
		log.Println("actList", i, "::", actList[i])
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
		for i<len(feature) {
			//diff := actList[i] - i
			//only output those not in actList
		}
	}
}

func init() {
	Stdin = bufio.NewReader(os.Stdin)
}