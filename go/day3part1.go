package main

/*
aoc2022 - Day 3, Part 1
-----------------------
The Rucksack Packing Challenge

Each Input line contains a definition of what's in a Rucksack. It's even-numbered and each half line represents a compartment load.
One item (character) is mis-placed per Rucksack - it appears in both compartments. Identify it.
This item can be prioritised (a-z : 1-26, A-Z: 27-52)
Return the SUM of the PRIORITIES of each duplicate element.
*/

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

// for testing:
//var dataFile = "data/day3test.txt"

// for stars:
var dataFile = "data/day3input.txt"

// I know I'm being sold a dummy that the chars are really,
// truly just ints but it makes the code easier to understand
// just slower...
func findConcurrence(str1 string, str2 string) string {
	concurrence := ""
	for _, refitem := range str1 {
		item := string(refitem)
		if strings.Contains(str2, item) && !(strings.Contains(concurrence, item)) {
			concurrence = concurrence + item
		}
	}
	return concurrence
}

// Damn them for not letting us do a simple subtraction on range
func runeToPriority(item rune) int {
	if item >= 'A' && item <= 'Z' {
		return int(item-'A') + 27
	} else if item >= 'a' && item <= 'z' {
		return int(item-'a') + 1
	} else {
		return 0
	}
}

func main() {
	file, err := os.Open(dataFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	s := bufio.NewScanner(file)
	sumPriorities := 0
	for s.Scan() {
		//The rucksack is the whole string:
		rucksack := s.Text()
		//The compartments are each half of the string:
		compA := rucksack[:len(rucksack)/2]
		compB := rucksack[len(rucksack)/2:]
		//sanity check
		if len(compA) != len(compB) {
			log.Fatal("Rucksack was not evenly packed; violation of contract")
		}
		matches := findConcurrence(compA, compB)
		//sanity check
		if len(matches) != 1 {
			log.Fatal("More than one item shared between compartments; violation of contract")
		}
		priority := runeToPriority(rune(matches[0]))
		fmt.Printf("%v(%02d) : %v/%v\n", matches, priority, compA, compB)
		sumPriorities += priority
	}
	fmt.Printf("Final Priority Sum: %d\n", sumPriorities)

}
