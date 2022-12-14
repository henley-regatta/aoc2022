package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
)

/*
aoc2022 - Day 14, Part 2
------------------------
The Trapped-In-A-Sand-Cave Path Mapping Problem

Given an input of a "vertical map" reporting lines of rock above you
(one "path" per input line), and a defined input of sand pouring in from
500,0, iterate at 1-unit-of-sand *moving at a time*.
Sand falls down. If it can't move directly down, it can move diagonally.
First it will move down and left. If it can't do that, it'll move down and right.
If all three directions (down, down-left, down-right) are blocked (occupied by rock or sand)
it comes to rest and the next unit of sand descends.

If the floor is at yMax+2, how many units of sand have to fall before the entry point is blocked?

(I was right to grow the width and depth but I need to make the map square...)

*/

// for testing:
//var dataFile = "data/day14test.txt"

// for Stars:
var dataFile = "data/day14input.txt"

type texture int

const (
	v texture = iota
	r
	s
	e
)

func (t texture) String() string {
	switch t {
	case v:
		return "." //void
	case r:
		return "#" //rock
	case s:
		return "o" //sand
	case e:
		return "+" //sand entry point
	}
	return "?" //unknown
}

type coord struct {
	x int
	y int
}

func buildMap(fromFile string, sandPoint coord) ([][]texture, coord) {

	file, err := os.Open(dataFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	s := bufio.NewScanner(file)
	//I'm too thick to work out how to do this in a single pass.
	mDim := [2]coord{}
	mDim[0] = coord{x: math.MaxInt, y: math.MaxInt}
	mDim[1] = coord{x: 0, y: 0}
	rDefs := [][]coord{}
	for s.Scan() {
		rPoints := strings.Split(s.Text(), " -> ")
		rDef := []coord{}
		for _, p := range rPoints {
			c := strings.Split(p, ",")
			pX, _ := strconv.Atoi(c[0])
			pY, _ := strconv.Atoi(c[1])
			if pX < mDim[0].x {
				mDim[0].x = pX
			}
			if pX > mDim[1].x {
				mDim[1].x = pX
			}
			if pY < mDim[0].y {
				mDim[0].y = pY
			}
			if pY > mDim[1].y {
				mDim[1].y = pY
			}
			rDef = append(rDef, coord{x: pX, y: pY})
		}
		rDefs = append(rDefs, rDef)
	}

	//For part 2, the maximum amount of sand is a diagonal
	//pyramid, with the tip at the sand entry point.
	//note: the sand COULD move as far across as down, so scale for that
	//We need to scale AND OFFSET the map to accomodate this:
	fmt.Printf("Unaltered range: %v to %v\n", mDim[0], mDim[1])
	yRange := mDim[1].y + 2 //because of the floor

	//Do we need to move the X offset to account for entrypoint?
	dX := yRange + 1
	fmt.Printf("Y range: %d. Half-X range: %d\n", yRange, dX)
	fmt.Printf("entry point req range: %d to %d\n", sandPoint.x-dX, sandPoint.x+dX)
	if sandPoint.x-dX < mDim[0].x {
		mDim[0].x = (sandPoint.x - dX)
	}
	if sandPoint.x+dX > mDim[1].x {
		mDim[1].x = (sandPoint.x + dX)
	}
	fmt.Printf("Entry Offset range: %v to %v\n", mDim[0], mDim[1])

	mDim[0] = coord{x: mDim[0].x - 1, y: mDim[0].y}
	mDim[1] = coord{x: mDim[1].x + 1, y: mDim[1].y + 2} //because of the floor

	xRange := mDim[1].x - mDim[0].x
	xOffset := mDim[0].x

	fmt.Printf("Offsetting X by %d\n", xOffset)

	//Move the sand entry point. Note we cannot offset Y as we need the space for the
	//sand to fall
	newSandpoint := coord{sandPoint.x - xOffset, sandPoint.y}
	fmt.Printf("Sand entry point moves from %v to %v\n", sandPoint, newSandpoint)

	//We can now initialise the map (to a void...)
	cavemap := [][]texture{}
	for y := 0; y < mDim[1].y; y++ {
		mLine := []texture{}
		for x := 0; x < xRange; x++ {
			mLine = append(mLine, v)
		}
		cavemap = append(cavemap, mLine)
	}
	fmt.Printf("Map size is %d x %d\n", len(cavemap[0]), len(cavemap))
	//Add the sand entry point
	cavemap[newSandpoint.y][newSandpoint.x] = e

	//Now we can iterate over our rock definitions adding them to the map
	for _, r := range rDefs {
		cavemap = addRockToCavemap(cavemap, xOffset, r)
	}

	return cavemap, newSandpoint

}

func addRockToCavemap(cavemap [][]texture, xOffset int, rockdef []coord) [][]texture {
	//the rock definitions are a string of coords and each possible pair
	//denotes a line segment.
	p0 := rockdef[0]
	p0.x = p0.x - xOffset
	for i := 1; i < len(rockdef); i++ {
		p1 := rockdef[i]
		p1.x = p1.x - xOffset
		//draw a line from p0 - p1
		cavemap[p0.y][p0.x] = r
		cavemap[p1.y][p1.x] = r
		//Lines are always horizontal or vertical
		//(pardon my spaghetti I no math today)
		if (p1.x - p0.x) > 0 {
			for x := p0.x; x < p1.x; x++ {
				cavemap[p1.y][x] = r
			}
		} else if (p1.x - p0.x) < 0 {
			for x := p1.x; x < p0.x; x++ {
				cavemap[p1.y][x] = r
			}
		} else if (p1.y - p0.y) > 0 {
			for y := p0.y; y < p1.y; y++ {
				cavemap[y][p0.x] = r
			}
		} else if (p1.y - p0.y) < 0 {
			for y := p1.y; y < p0.y; y++ {
				cavemap[y][p0.x] = r
			}
		} else {
			fmt.Printf("Weird line segment %v - %v\n", p0, p1)
		}
		p0 = p1
	}

	return cavemap
}

func printMap(cavemap [][]texture) {
	for y := 0; y < len(cavemap); y++ {
		caveLine := ""
		for x := 0; x < len(cavemap[y]); x++ {
			caveLine = caveLine + fmt.Sprint(cavemap[y][x])
		}
		fmt.Println(caveLine)
	}
	return
}

func addSand(cavemap [][]texture, sandpoint coord) ([][]texture, bool) {
	sandMoving := true
	sP := sandpoint
	for sandMoving && sP.y < len(cavemap)-1 {
		//Can we move sand down?
		if cavemap[sP.y+1][sP.x] == v {
			sP.y = sP.y + 1
		} else if cavemap[sP.y+1][sP.x-1] == v { //What about diag-left?
			sP.y = sP.y + 1
			sP.x = sP.x - 1
		} else if cavemap[sP.y+1][sP.x+1] == v { //What about diag-right?
			sP.y = sP.y + 1
			sP.x = sP.x + 1
		} else { //No, it's stuck.
			sandMoving = false
		}
	}
	cavemap[sP.y][sP.x] = s
	if sP == sandpoint {
		return cavemap, false
	}
	return cavemap, true
}

func main() {
	entryPoint := coord{x: 500, y: 0}
	caveMap, entryPoint := buildMap(dataFile, entryPoint)
	printMap(caveMap)
	fmt.Println("-----------------------------------------------------")
	sandAdded := true
	sandCounter := 0
	for sandAdded {
		caveMap, sandAdded = addSand(caveMap, entryPoint)
		sandCounter++
	}
	printMap(caveMap)
	fmt.Printf("%d units of sand came to rest before chaos\n", sandCounter)
}
