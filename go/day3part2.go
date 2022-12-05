package main

/*
aoc2022 - Day 3, Part 2
-----------------------
The Rucksack Packing Challenge

Each Input line contains a definition of what's in a Rucksack.
Each 3 lines determines an Elve Group of rucksacks
One item is shared between all 3 rucksacks. Identify it.
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
// PART2: Also makes a 3-way comparison easier, lol.
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

func threewaymatch(rucksack [3]string) string {
	firstconcurrence := findConcurrence(rucksack[0], rucksack[1])
	match := findConcurrence(firstconcurrence, rucksack[2])
	fmt.Printf("%v is concurrence of %v, %v, %v\n", match, rucksack[0], rucksack[1], rucksack[2])
	if len(match) != 1 {
		log.Fatal("Error determining concurrence, should have had a single match had: " + match)
	}
	return match
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
	var elfGroup [3]string
	ptr := 0
	for s.Scan() {
		//The rucksack is the whole string:
		elfGroup[ptr] = s.Text()
		if ptr == 2 {
			//We have an elf group, process
			matchingElement := threewaymatch(elfGroup)
			elemPriority := runeToPriority(rune(matchingElement[0]))
			sumPriorities += elemPriority
			//reset to start again
			ptr = 0
		} else {
			ptr++
		}
	}
	//Sanity check - we should have JUST FINISHED the final elfgroup and the ptr should be zero
	if ptr != 0 {
		log.Fatal("Input did not contain a complete set of elf backback groups, Contract violation")
	}
	fmt.Printf("Final Priority Sum: %d\n", sumPriorities)

}
