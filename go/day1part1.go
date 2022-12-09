package main

/*
aoc2022 - Day 1, Part 1
-----------------------
Hungry Hungry Elves

Each  Elf enumerates the calorific content of their snacks, 1 snack per line.
Elves are separated in the input by a blank line.

Find the Elf carrying the most Calories. How  many total Calories in that Elf carrying?
*/

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
)

// for testing:
// var dataFile = "data/day1test.txt"
// for Stars:
var dataFile = "data/day1input.txt"

func main() {
	file, err := os.Open(dataFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	s := bufio.NewScanner(file)

	//Given each LINE is either blank (delimiting elves) or
	//numeric, the obvious thing to do here is build a list
	//incrementing the value by calories per line. However,
	//traversal is a lot easier if we just build a Map
	elfMap := map[int]int{}
	ndx := 1

	for s.Scan() {
		w := s.Text()
		n, err := strconv.Atoi(w)
		if err != nil {
			//couldn't convert str to int, must be new elf
			ndx++
		} else {
			elfMap[ndx] += n
		}
	}

	//Work out which one had the most calories
	maxCalsFound := 0
	chonkiestElf := 0
	for elf, calories := range elfMap {
		if calories > maxCalsFound {
			maxCalsFound = calories
			chonkiestElf = elf
		}
	}
	fmt.Println(chonkiestElf)
	fmt.Println(maxCalsFound)
}
