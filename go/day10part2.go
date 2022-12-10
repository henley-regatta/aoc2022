package main

/*
aoc2022 - Day 10, Part 2
-----------------------
The Virtual Machine Debugging Exercise.
A CPU has a since register, X, initialised to 1
Two instructions are supported:
	"addx V" (2 cycles to complete). On completion, X = X + V
	"noop" (1 cycle to complete). Does nothing

X now records the mid-point of a 3-pixel sprite. This is drawn on a CRT
of dimensions x=40, y=6 (zero indexed). x=0 is to the left.
The CRT is drawn 1 pixel per cycle. It takes 40 cycles to draw a line,
40*6 = 240 cycles to refresh the whole screen
IF the sprite position (X-1, X, X+1) matches the CRT horizontal position (0-39)
during a given cycle, THEN the CRT pixel is lit (#). Otherwise it's set dark (.)

Execute the program for a full screen refresh (240 cycles). Display the image
formed. What 8 capital letters appear on the CRT ?

*/

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

// Helper function to turn our output into visible stuff
func dumpDisplay(crt [6][40]bool) {
	for y := 0; y < len(crt); y++ {
		line := ""
		for x := 0; x < len(crt[y]); x++ {
			char := " "
			if crt[y][x] {
				char = "#"
			}
			line = line + char
		}
		fmt.Println(line)
	}
}

// Guts of the output generation. Given regX, cycle count
// and current display, work out whether to light a pixel
func updateDisplay(cycle, regX int, crt [6][40]bool) [6][40]bool {
	//Where's the CRT updating?
	//TODO: chance of off-by-one based on cycle to pixel mapping
	crtY := (cycle - 1) / 40
	crtX := (cycle - 1) % 40

	//Sensecheck. Don't update if crtY out of bounds
	if crtY > 5 {
		return crt
	}
	//Here's the rule. IF regX +/-1 matches crtX, pixel becomes lit
	if crtX == (regX-1) || crtX == regX || crtX == (regX+1) {
		crt[crtY][crtX] = true
	} else {
		crt[crtY][crtX] = false
	}
	return crt
}

// testing (part one expected answer 13140)
// (part two displays an image too complex to embed)
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

	//Define and initialise the CRT display
	crt := [6][40]bool{}
	for y := 0; y < len(crt); y++ {
		for x := 0; x < len(crt[y]); x++ {
			crt[y][x] = false
		}
	}

	//Initialise the Virtual Machine
	regX := 1
	cycle := 1

	//Perform first update
	crt = updateDisplay(cycle, regX, crt)

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
			//Update the display in the mid-cycle:
			crt = updateDisplay(cycle+1, regX, crt)
			//Now we're past that, complete the instruction
			cycle += 2
			regX = regX + parmV

		} else {
			log.Fatal("Error, invalid instruction: " + s.Text())
		}
		//Update CRT on cycle completion
		crt = updateDisplay(cycle, regX, crt)

	}

	fmt.Printf("After %d cycles, Final Image:\n", cycle)
	dumpDisplay(crt)
}
