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
	// STEP 1:
	// Begin building training and test set

	fmt.Println("Building train and test set")
	fmt.Println("What file would you like to split?")
	fmt.Print("file name> ")
	// Receive file name of data file
	_, err = Scanf("%s", &inputString)
	errCheck(err)
	debugMsg("Opening file:", inputString)
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
	var exists bool       // For checking if element exists
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
			tempFileName := dataFile.Name() + "." + label + ".tmp"
			debugMsg("Creating temporary file:", tempFileName)
			tempFile, err := os.Open(
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
	debugMsg("Closing all temporary files for writing. Re-opening as")
	debugMsg("read-only")
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
		debugMsg("inputInt:", inputInt)
		trainCountMap[k] = inputInt
	}
	debugMsg("Creating:", dataFile.Name()+".train")
	debugMsg("Creating:", dataFile.Name()+".test")
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
		debugMsg("label:", k, "count:", v)
		if v > 0 {
			// Generate a random permuation
			rand := rand.Perm(countMap[k])
			// use a slice the first /v/ of them
			rand = rand[0:v]
			// sort the ints so that as we iterate through each instance we can
			// easily find the next one we need to export
			sort.SortInts(rand)
			// Read through the file, writing the included instances to
			// .train and the others to .test
			dataReader := bufio.NewReader(tempFileMap[k])
			lineCount := 0
			if len(rand) > 0 {
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
	}
	fmt.Println()
}
