package main

/*
aoc2022 - Day 1, Part 2
-----------------------
Hungry Hungry Elves

Each  Elf enumerates the calorific content of their snacks, 1 snack per line.
Elves are separated in the input by a blank line.

Find the top THREE chonk-carrying Elves. How many calories do they carry all together?
*/

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
)

// for testing:
//var dataFile = "../data/day1test.txt"

// for Stars:
var dataFile = "../data/day1input.txt"

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

	//Populate the top 3 chonkiest carriers
	//(this is easiest to do by sorting the input by VALUE instead of KEY and taking the last 3 entries)
	elfNdx := make([]int, 0, len(elfMap))
	for e := range elfMap {
		elfNdx = append(elfNdx, e)
	}

	//Sort the elfNdx by the number of calories not the index:
	sort.SliceStable(elfNdx, func(i, j int) bool {
		return elfMap[elfNdx[i]] > elfMap[elfNdx[j]]
	})

	//Now we just need a slice of the first 3 elements in this sorted elfNdx
	var chonks []int = elfNdx[0:3]

	totCals := 0
	for _, e := range chonks {
		fmt.Println(e, elfMap[e])
		totCals += elfMap[e]
	}
	fmt.Println(totCals)

}
