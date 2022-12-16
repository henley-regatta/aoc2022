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
aoc2022 - Day 15, Part 2
------------------------
The Beacon Exclusion Zone Combinatorial Explosion

Given a list of sensor,beacon input and using Manhattan distances, compute
the only point in a given range where a beacon could NOT be, and
return the "tuning frequency" of those coordinates.

*/

// for testing:
//var dataFile = "data/day15test.txt"

//const TARGETRANGE = 20

// for Stars:
var dataFile = "data/day15input.txt"

const TARGETRANGE = 4000000

type coord [2]int      //convention x,y
type sensData struct { //The beacon is pointless we only need it's manhattan distance from sensor
	s coord
	m int
}

func parseInput(fromFile string) []sensData {
	file, err := os.Open(dataFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	s := bufio.NewScanner(file)

	datalist := []sensData{}
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
				sensPos := coord{nc[0], nc[1]}
				mDist := manhattanDistance(sensPos, coord{nc[2], nc[3]})
				datalist = append(datalist, sensData{s: sensPos, m: mDist})
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

func rowCoverage(row int, sensors []sensData) []coord {
	rC := []coord{}
	for _, v := range sensors {
		if v.s[1]-v.m < row && v.s[1]+v.m > row {
			//There is an intersection. Work out where
			yDist := intAbs(v.s[1] - row)
			txMin := v.s[0] - (v.m - yDist)
			txMax := v.s[0] + (v.m - yDist)
			rC = append(rC, coord{txMin, txMax})
		}
	}
	//SHORTCUT
	if len(rC) == 0 {
		return []coord{}
	}
	//consolidate the coverage. Needs sorted input:
	sort.SliceStable(rC, func(i, j int) bool {
		return rC[i][0] < rC[j][0]
	})
	//now iterate and compute coverage:
	cR := []coord{}
	for _, v := range rC {
		covered := false
		for i := range cR {
			if v[0] >= cR[i][0] && v[1] <= cR[i][1] {
				covered = true //completely covered
			} else if v[0] >= cR[i][0] && v[0] <= (cR[i][1]+1) && v[1] > cR[i][1] {
				cR[i][1] = v[1] //left-overlap; extend original right
				covered = true
			} else if v[0] < cR[i][0] && v[1] >= (cR[i][0]-1) && v[1] <= cR[i][1] {
				cR[i][0] = v[0] //right-overlap; extend original left
				covered = true
			}
		}
		//Check whether we need to add the range in
		if !covered {
			cR = append(cR, v)
		}
	}
	return cR
}

func main() {
	sensordata := parseInput(dataFile)

	var gap coord
	gapList := []coord{}
	//An iterative search....
	for y := 0; y <= TARGETRANGE; y++ {
		coverage := rowCoverage(y, sensordata)
		if len(coverage) != 1 {
			gap = coord{coverage[0][1] + 1, y}
			if coverage[1][0] > gap[0]+1 {
				fmt.Println(y)
				fmt.Println(coverage)
				fmt.Println(gap)
				log.Fatal("Error computing gap, should be 1 big")
			} else {
				//Ideally we'd bang-out here but...
				fmt.Printf("Gap found at: %v\n", gap)
				gapList = append(gapList, gap)
			}
		}
		/*
			if y%1000 == 0 {
				fmt.Printf("%05d rows complete. %05d remaining (%02.0f%% complete)\n", y, TARGETRANGE-y, float32(y)/float32(TARGETRANGE)*100)
			}
		*/
	}

	if len(gapList) != 1 {
		fmt.Println(gapList)
		log.Fatal("Something fucky. Should have found 1 and only 1 gap")
	}

	tuningFreq := gapList[0][0]*4000000 + gapList[0][1]
	fmt.Printf("The Tuning Frequency You Seek is: %d\n", tuningFreq)

}
