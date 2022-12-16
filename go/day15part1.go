package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"sort"
	"strconv"
)

/*
aoc2022 - Day 15, Part 1
------------------------
The Beacon Exclusion Zone Combinatorial Explosion

Given a list of sensor,beacon input and using Manhattan distances, and a
defined input row, compute the number of points where a putative rescue
beacon could NOT be (given that each sensor only reports it's nearest
beacon, so their are inevitably holes in the coverage).

*/

// for testing:
//var dataFile = "data/day15test.txt"

//const TARGETROW = 10

// for Stars:
var dataFile = "data/day15input.txt"

const TARGETROW = 2000000

type coord [2]int            //convention x,y
type sensBeaconPair [2]coord //convention s,b

func parseInput(fromFile string) []sensBeaconPair {
	file, err := os.Open(dataFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	s := bufio.NewScanner(file)

	datalist := []sensBeaconPair{}
	for s.Scan() {
		//dodgy but workable
		coordRE := regexp.MustCompile("(-*[0-9]+)")
		allcoords := coordRE.FindAllString(s.Text(), 4)
		if len(allcoords) == 4 {
			nc := []int{}
			for i := range allcoords {
				n, err := strconv.Atoi(allcoords[i])
				if err == nil {
					nc = append(nc, n)
				}
			}
			if len(nc) == 4 {
				datalist = append(datalist, sensBeaconPair{coord{nc[0], nc[1]}, coord{nc[2], nc[3]}})
			}

		}
	}

	return datalist
}

func intAbs(x int) int {
	if x < 0 {
		return x * -1
	}
	return x
}
func manhattanDistance(p1, p2 coord) int {
	xD := intAbs(p1[0] - p2[0])
	yD := intAbs(p1[1] - p2[1])
	return xD + yD
}

func main() {
	sensorbeaconpairs := parseInput(dataFile)

	//We now have to work out which of the sensors has "field output" on the target row
	//and which X positions that output covers. This is a piece of piss except for overlaps.
	//Work out the range of X values we need to covers
	coverage := []coord{}
	for _, v := range sensorbeaconpairs {
		mDist := manhattanDistance(v[0], v[1])
		if v[0][1]-mDist < TARGETROW && v[0][1]+mDist > TARGETROW {
			//There is an intersection. Work out where
			yDist := intAbs(v[0][1] - TARGETROW)
			txMin := v[0][0] - (mDist - yDist)
			txMax := v[0][0] + (mDist - yDist)
			coverage = append(coverage, coord{txMin, txMax})
		}
	}
	//Now work out overlaps. This only works single-pass if the input list
	//is sorted by starting point:
	sort.SliceStable(coverage, func(i, j int) bool {
		return coverage[i][0] < coverage[j][0]
	})
	fmt.Printf("Raw Coverage at row=%d : \n%v\n", TARGETROW, coverage)
	cR := []coord{}
	for _, v := range coverage {
		covered := false
		for i := range cR {
			if v[0] >= cR[i][0] && v[1] <= cR[i][1] {
				covered = true //completely covered
			} else if v[0] >= cR[i][0] && v[0] <= cR[i][1] && v[1] > cR[i][1] {
				cR[i][1] = v[1] //left-overlap; extend original right
				covered = true
			} else if v[0] < cR[i][0] && v[1] >= cR[i][0] && v[1] <= cR[i][1] {
				cR[i][0] = v[0] //right-overlap; extend original left
				covered = true
			}
		}
		//Check whether we need to add the range in
		if !covered {
			cR = append(cR, v)
		}
	}
	fmt.Printf("Consolidated coverage at row=%d : \n%v\n", TARGETROW, cR)
	//And the answer we seek is just the sum of the ranges in that consolidated coverage
	covered := 0
	for _, v := range cR {
		covered += (v[1] - v[0])
	}
	fmt.Printf("At row=%d there are %d covered elements\n", TARGETROW, covered)
}
