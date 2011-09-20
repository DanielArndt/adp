/* 
 * convert.go
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
	vector "container/vector"
	"fmt"
	"log"
	"os"
	"exp/regexp"
	"strings"
)

type feature struct {
	name     string
	datatype string
}

type header struct {
	name     string
	features vector.Vector
}

var attrRE *regexp.Regexp = regexp.MustCompile("^@attribute\\s+(\\w+)\\s+(\\w+)")

var relRE *regexp.Regexp = regexp.MustCompile("^@relation\\s+(\\S+)")

// For now just skips the header
func readHeader(infile *bufio.Reader) header {
	var hdr header
	for line, err := infile.ReadString('\n'); err == nil && !strings.HasPrefix(strings.ToLower(line), "@data"); line, err = infile.ReadString('\n') {
		line = strings.ToLower(line)
		match := attrRE.FindStringSubmatch(line)
		if match != nil {
			//This line is defining an attribute
			if len(match) < 3 {
				log.Fatalf("Invalid header on line:\n%s", line)
			}
			hdr.features.Push(match[1])
			continue
		}
		match = relRE.FindStringSubmatch(line)
		if match != nil {
			if hdr.name == "" {
				hdr.name = match[1]
			} else {
				log.Fatalf("Double declaration of relation name?\n" +
					"First: %s\n"+
					"Second: %s\n", hdr.name, match[1])
			}
		}
		//TODO: Need to do some parsing of header to create c50 files.
	}
	log.Printf("Read in %d attributes\n", len(hdr.features))
	return hdr
}

func writelineSbbFive(line string,
	datafile *os.File,
	labelfile *os.File) {
	features := strings.Split(line, ",")

	data := features[0 : len(features)-2]
	label := features[len(features)-1]
	fmt.Printf("%s - %s\n", data, label)
}

func interactiveConvert() {
	fmt.Println("Converting data file from ARFF to multiple formats.")
	arffFileName := promptString("arff file",
		"Please enter the path of the ARFF file")
	debugMsg("Opening file: %s", arffFileName)
	// Open the file for input and create a buffered reader for the file
	infileFD, err := os.Open(arffFileName)
	// We do not need this file after, so close it upon leaving this method
	defer infileFD.Close()
	errCheck(err)

	infile := bufio.NewReader(infileFD)
	hdr := readHeader(infile)
	log.Println(hdr)
/*	sbbFiveData, err := os.Create(arffFileName + ".sbb5.data")
	errCheck(err)
	sbbFiveLabels, err := os.Create(arffFileName + ".sbb5.labels")
	errCheck(err)
	for line, err := infile.ReadString('\n'); err == nil; line, err = infile.ReadString('\n') {
		line = strings.TrimRight(line, "\n ")
		if len(line) < 1 {
			continue
		}
		writelineSbbFive(line, sbbFiveData, sbbFiveLabels)
	}*/
}
