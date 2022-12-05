package main

/*
aoc2022 - Day 5, Part 2
-----------------------
The Crate Unloading Optimisation Cycle

Given an input of:
a a representation of a "stack" of crates in original format
b a list of instructions of the form MOVE x (blocks) FROM s (column) TO f (column)

Apply the list of instructions to the start stack, with the rule that ALL blocks come in a
single move. (I.e. this is identical to part 1 except it ain't quite a stack it's a shift)

Report the TOP blocks from each stack

Unlike other day's input, parsing the start stacks looks like no fun at all. So don't bother.
*/

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func defData(isTesting bool) (fname string, cratez []string) {
	if isTesting {
		fname = "data/day5test.txt"
		cratez := []string{"ZN", "MCD", "P"}
		return fname, cratez
	} else {
		fname = "data/day5input.txt"
		cratez := []string{"NCRTMZP", "DNTSBZ", "MHQRFCTG", "GRZ", "ZNRH", "FHSWPZLD", "WDZRCGM", "SJFLHWZQ", "SQPWN"}
		return fname, cratez
	}

}

func getStackInstruction(line string) (valid bool, moves int, from int, to int) {
	//Line MUST start with "move" to be a valid instruction
	if !strings.Contains(line, "move") {
		return false, 0, 0, 0
	} else {
		//lazy parsing the input is structured enough not to sentinel/check
		wrds := strings.Split(line, " ")
		valid = true
		moves, _ = strconv.Atoi(wrds[1])
		from, _ = strconv.Atoi(wrds[3])
		to, _ = strconv.Atoi(wrds[5])
		//nb our array is 0-indexed but the words are 1-indexed
		return valid, moves, from - 1, to - 1
	}

}

// this is only now used to read out the final state so it's overcomplex
func popStack(inStack []string, col int) (crate string, outStack []string) {
	outStack = inStack
	if len(outStack[col]) == 1 {
		crate = outStack[col]
		outStack[col] = ""
	} else {
		column := []rune(outStack[col])
		crate = string(column[len(column)-1])
		remCol := column[:len(column)-1]
		outStack[col] = string(remCol)
	}
	return crate, outStack
}

// since part 2 is all about SHIFTING a block instead of stack-manip, this func
// now does it all.
func execMoveInstruction(inStack []string, n int, f int, t int) (outStack []string) {
	outStack = inStack
	var moveCrates string
	//Strip the last "n" blocks from fromCol
	if len(outStack[f]) == n {
		moveCrates = outStack[f]
		outStack[f] = ""
	} else {
		moveCrates = outStack[f][len(outStack[f])-(n):]
		outStack[f] = outStack[f][:len(outStack[f])-(n)]
	}
	//that's the source tackled, now append to dest:
	outStack[t] = outStack[t] + moveCrates
	return outStack
}

func printStack(stack []string) {
	for i := range stack {
		fmt.Printf("%1d: %02d %v\n", i+1, len(stack[i]), stack[i])
	}
}

// for testing:
var runTest = false

func main() {
	dataFile, stack := defData(runTest)
	printStack(stack)
	file, err := os.Open(dataFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	s := bufio.NewScanner(file)
	mvCtr := 0
	for s.Scan() {
		y, n, f, t := getStackInstruction(s.Text())
		if y { //only on valid lines
			mvCtr++
			fmt.Printf("Moving %d blocks from stack %d to stack %d\n", n, f+1, t+1)
			stack = execMoveInstruction(stack, n, f, t)
			//printStack(stack)
		}
	}

	fmt.Println("Final Stack State")
	printStack(stack)
	topvals := ""
	var c string
	for i := range stack {
		c, stack = popStack(stack, i)
		topvals = topvals + c
	}
	fmt.Printf("Top Stack Crates after %d moves: %v\n", mvCtr, topvals)
}
