package main

/*
aoc2022 - Day 9, Part 2
-----------------------
The Knotty Rope-Bridge Snake Game - Wag the Tail edition

Given some bullshit about knots, and a series of instructions to move a "Head" around a grid,
count the number of unique positions a "Tail" visits during that walk.

The elaboration: the Tail is now 9 elements long. Each tracks the position of the one ahead of it.
Same movement rules apply: Tail[n] must be adjacent to Tail[n-1] or Head.

Report on the number of unique spots the END tail visits.

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

// testing (for part 2 should produce answer = 1, tail never moves)
//var dataFile = "data/day9test.txt"

// testing 2 (for part 2 should produce answer = 36)
//var dataFile = "data/day9test2.txt"

// live
var dataFile = "data/day9input.txt"

// big boi - expects p1 = 8129855 p2 = 7750850
//var dataFile = "data/input-2022-09-bb-100000.txt"

func main() {
	file, err := os.Open(dataFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	s := bufio.NewScanner(file)

	//Head and Tail now a 10-element tuple
	rope := [10][2]int{}
	for i := 0; i < 10; i++ {
		rope[i][0] = 0
		rope[i][1] = 0
	}
	//Luckily we still only need to track the tail positions:
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
							rope[0][1]++
						}
					case "D":
						{
							rope[0][1]--
						}
					case "R":
						{
							rope[0][0]++
						}
					case "L":
						{
							rope[0][0]--
						}
					default:
						{
							log.Fatal("Invalid direction " + ins[0])
							os.Exit(1)
						}
					}
					//It's iterate-through-an-array time!
					for i := 1; i < len(rope); i++ {
						rope[i] = checkMoveTPos(rope[i-1], rope[i])
					}
					//And now, the end is near, we can decide how to update
					//the tail position marker:

					tK := tKey(rope[9])
					if _, ok := tLocs[tK]; ok {
						tLocs[tK]++ //repeat visit; don't care but why not log it
					} else {
						tLocs[tK] = 1 //fresh location, mark it
					}
				}
			}
		}
	}

	fmt.Printf("Head final position: %v Tail final position: %v\n", rope[0], rope[9])
	//Because of the way we've used a map to track tail position, the answer we
	//seek - the number of unique locations visited by the tail - is just the
	//number of keys in the map
	fmt.Println(len(tLocs))
}
