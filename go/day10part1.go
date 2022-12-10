package main

/*
aoc2022 - Day 10, Part 1
-----------------------
The Virtual Machine Debugging Exercise.
A CPU has a since register, X, initialised to 1
Two instructions are supported:
	"addx V" (2 cycles to complete). On completion, X = X + V
	"noop" (1 cycle to complete). Does nothing

Given a list of instructions, calculate the "signal strength" (cycle number * X value)
at T = 20 and then every 40 cycles (20, 60, 100, 140, 180, 220).

Report the sum of signal strengths.
*/

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

// Helper to act as array pop
func intPop(list []int) []int {
	if len(list) == 1 {
		log.Fatal("attempting to pop past the last element on the array")
	}
	return list[1:]
}

// Bookkeeping function
func recordSigStrength(cycle, regX, sumSigStrength int, pPoints []int) (int, []int) {
	pVal := pPoints[0] * regX
	fmt.Printf("Cycle = %d, Reg X = %d, sigStrength = %d\n", cycle+1, regX, pVal)
	sumSigStrength += pVal
	if len(pPoints) == 1 { //Last measurement cycle. We can finish now
		fmt.Printf("Reached final measurement point. Cycles=%d, regX now %d. Sum of signal strength = %d\n", cycle, regX, sumSigStrength)
		os.Exit(0)
	}

	return sumSigStrength, intPop(pPoints)
}

// testing (part one expected answer 13140)
//var dataFile = "data/day10test.txt"

// live
var dataFile = "data/day10input.txt"

func main() {
	file, err := os.Open(dataFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	s := bufio.NewScanner(file)

	//Define our probe points (cycle times)
	probePoints := []int{20, 60, 100, 140, 180, 220}

	//the value we'll finally report:
	sumSigStrength := 0

	//Initialise the Virtual Machine
	regX := 1
	cycle := 1

	//Loop. Read and Execute the instructions
	for s.Scan() {
		//DECODE
		ins := strings.Split(s.Text(), " ")
		if len(ins) == 1 && ins[0] == "noop" {
			//It's a no-op. All we do is increment cycle
			cycle++
		} else if len(ins) == 2 && ins[0] == "addx" {
			//V is ins[2]
			parmV, e := strconv.Atoi(ins[1])
			if e != nil {
				log.Fatal("Error, invalid parameter in instruction: " + s.Text())
			}
			//SO the only trick here is we need to check whether cycle+1 is at a
			//probe-point before we complete the instruction in cycle+2
			if cycle+1 == probePoints[0] {
				sumSigStrength, probePoints = recordSigStrength(cycle+1, regX, sumSigStrength, probePoints)
			}
			//Now we're past that, complete the instruction
			cycle += 2
			regX = regX + parmV

		} else {
			log.Fatal("Error, invalid instruction: " + s.Text())
		}

		//On completion of instruction execution, check for whether we're at a probe point:
		if cycle == probePoints[0] {
			sumSigStrength, probePoints = recordSigStrength(cycle, regX, sumSigStrength, probePoints)
		}
	}
	fmt.Printf("Ran off the end of the instruction list. Cycles=%d, regX now %d. Sum of signal strength = %d\n", cycle, regX, sumSigStrength)
}
