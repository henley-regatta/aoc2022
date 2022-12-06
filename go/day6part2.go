package main

/*
aoc2022 - Day 6, Part 2
-----------------------
The Randomising Radio Packet Initialiser

Given an input of seemingly-random characters making up a string,
find the first set of 14 distinct/non-repeating characters that mark a message header.
Report the position of the last character scanned.
*/

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

// fixed size of header
const HDRLEN = 14

// for fun:
//var dataFile = "data/day6test.txt"

// for stars:
var dataFile = "data/day6input.txt"

func isUniqueHdr(hdr []rune) bool {
	for i := 0; i < len(hdr); i++ {
		for j := i + 1; j < len(hdr); j++ {
			if hdr[i] == hdr[j] {
				return false
			}
		}
	}
	return true
}

// The meat of the algo:
// Grab the first 4 chars, check for uniqueness.
func findHeader(message string) int {
	//Input sanity check
	if len(message) < HDRLEN {
		log.Fatal("Input message less than required header size. Breach of Contract")
		return -1
	}
	mRune := []rune(message) //convert input message to rune array.
	hdr := mRune[0:HDRLEN]   //Initialise the candidate header to the first 4 chars of the message
	for i := HDRLEN; i < len(mRune); i++ {
		//fmt.Printf("Testing index %d with header %v\n", i, hdr)
		//Is the hdr unique?
		if isUniqueHdr(hdr) {
			return i
		} else {
			//no it's not. Replace nth hdr char with one under pointer and increment
			hdr[i%HDRLEN] = mRune[i]
		}
	}
	//if we get to here, we've not found a header
	return -1
}

func main() {
	file, err := os.Open(dataFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	s := bufio.NewScanner(file)
	for s.Scan() {
		fmt.Println(s.Text())
		pos := findHeader(s.Text())
		fmt.Println(pos)
	}
}
