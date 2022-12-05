package main

/*
aoc2022 - Day 4, Part 1
-----------------------
The Beach Cleanup Optimisation

Given an input of format A-B,C-D representing 2 "distinct" number ranges allocated to 2 elves,
determine how many of these pairs *completely* overlap (i.e. where 1 range is a strict subset of the other.)

Report the number of overlaps found.
*/

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

// for testing:
//var dataFile = "data/day4test.txt"

// for stars:
var dataFile = "data/day4input.txt"

type Pairs struct {
	e1S int
	e1F int
	e2S int
	e2F int
}

func extractRanges(line string) Pairs {
	//line looks like "aaa-bbb,ccc-ddd"
	ranges := strings.Split(line, ",")
	if len(ranges) != 2 {
		log.Fatal("Could not extract 2 ranges from '" + line + "' - contract failure")
	}
	p1 := strings.Split(ranges[0], "-")
	p2 := strings.Split(ranges[1], "-")

	if len(p1) != 2 || len(p2) != 2 {
		log.Fatal("Could not extract numbers from ranges in '" + line + "' - contract failure")
	}
	var parsedPairs Pairs
	parsedPairs.e1S, _ = strconv.Atoi(p1[0])
	parsedPairs.e1F, _ = strconv.Atoi(p1[1])
	parsedPairs.e2S, _ = strconv.Atoi(p2[0])
	parsedPairs.e2F, _ = strconv.Atoi(p2[1])
	return parsedPairs
}

func main() {
	file, err := os.Open(dataFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	s := bufio.NewScanner(file)
	strictOverlaps := 0
	for s.Scan() {
		r := extractRanges(s.Text())
		//e1 can be fully inside e2 or vice-versa:
		if (r.e1S >= r.e2S && r.e1F <= r.e2F) || (r.e2S >= r.e1S && r.e2F <= r.e1F) {
			fmt.Println(r)
			strictOverlaps++
		}
	}
	fmt.Printf("Final Overlaps Sum: %d\n", strictOverlaps)

}
