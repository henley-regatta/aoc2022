package main

/*
aoc2022 - Day 9, Part 2
-----------------------
The Knotty Rope-Bridge Snake Game

Given some bullshit about knots, and a series of instructions to move a "Head" around a grid,
count the number of unique positions a "Tail" visits during that walk.

The "Tail" remains "touching" the Head - i.e. the Tail must be adjacent (horizontally,vertically OR
diagonally) to the Head.

*/

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

// helper to turn tPos into a unique-ish key
func tKey(tpos [2]int) string {
	return strconv.Itoa(tpos[0]) + "-" + strconv.Itoa(tpos[1])
}

// keep forgetting there's no integer Abs function...
func intAbs(x int) int {
	if x < 0 {
		return -1 * x
	}
	return x
}

// Determine adjacency of 2 positions
func adjacent(p1, p2 [2]int) bool {
	//non-adjacent if any coordinate differs by >1
	if intAbs(p2[0]-p1[0]) < 2 && intAbs(p2[1]-p1[1]) < 2 {
		return true
	}
	return false
}

// Actually this is the guts of the problem right here.
// Given a Head movement, work out whether and where we need to
// move Tail:
func checkMoveTPos(hpos, tpos [2]int) [2]int {
	//We only move tpos IF hpos is no longer adjacent.
	//Which we need to work out...
	if adjacent(hpos, tpos) {
		return tpos
	}

	newTpos := tpos
	//We need the position delta to proceed
	xDelta := hpos[0] - tpos[0]
	yDelta := hpos[1] - tpos[1]

	//We need a diagonal (2-step) move IF both Deltas are non-zero
	if xDelta != 0 && yDelta != 0 {
		if intAbs(xDelta) > intAbs(yDelta) {
			//we'll need to later move X but first null Y
			if yDelta > 0 {
				newTpos[1]++
			} else {
				newTpos[1]--
			}
			yDelta = 0
		} else {
			//We'll later move Y but first move X
			if xDelta > 0 {
				newTpos[0]++
			} else {
				newTpos[0]--
			}
			xDelta = 0
		}

	}
	//Now we've handled the diagonal, handle the cardinal
	if xDelta == 0 {
		if yDelta > 0 {
			newTpos[1]++
		} else {
			newTpos[1]--
		}
	} else if yDelta == 0 {
		if xDelta > 0 {
			newTpos[0]++
		} else {
			newTpos[0]--
		}
	}
	return newTpos
}

// testing
//var dataFile = "data/day9test.txt"

// live
var dataFile = "data/day9input.txt"

func main() {
	file, err := os.Open(dataFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	s := bufio.NewScanner(file)

	//We need to know H and T's locations.
	hpos := [2]int{0, 0}
	tpos := [2]int{0, 0}
	//I need to keep track of the number of distinct locations
	//"T" visits. A Set would be ideal, but bodging it with a map
	//with odd keys will do (x-y)
	tLocs := map[string]int{}
	for s.Scan() {
		//Each line contains an instruction consisting of a direction and a repeat factor
		ins := strings.Split(s.Text(), " ")
		if len(ins) == 2 {
			rep, e := strconv.Atoi(ins[1])
			if e == nil {
				//fmt.Printf("Moving %v for %d repeats\n", ins[0], rep)
				for i := 0; i < rep; i++ {
					switch ins[0] {
					case "U":
						{
							hpos[1]++
						}
					case "D":
						{
							hpos[1]--
						}
					case "R":
						{
							hpos[0]++
						}
					case "L":
						{
							hpos[0]--
						}
					default:
						{
							log.Fatal("Invalid direction " + ins[0])
							os.Exit(1)
						}
					}
					//Now we need to decide whether to move T and if so where
					tpos = checkMoveTPos(hpos, tpos)
					if _, ok := tLocs[tKey(tpos)]; ok {
						tLocs[tKey(tpos)]++ //repeat visit; don't care but why not log it
					} else {
						tLocs[tKey(tpos)] = 1 //fresh location, mark it
					}
				}
			}
		}
	}

	fmt.Printf("Head final position: %v Tail final position: %v\n", hpos, tpos)
	//Because of the way we've used a map to track tail position, the answer we
	//seek - the number of unique locations visited by the tail - is just the
	//number of keys in the map
	fmt.Println(len(tLocs))
}
