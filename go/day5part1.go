package main

/*
aoc2022 - Day 5, Part 1
-----------------------
The Crate Unloading Optimisation Cycle

Given an input of:
a a representation of a "stack" of crates in original format
b a list of instructions of the form MOVE x (blocks FROM s (column TO f (column

Apply the list of instructions to the start stack, with the rule that ONE block is moved at a time
and each stack is, er, a stack (first-in last-out.

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

// The main functions we use are a POP and a PUSH so we need functions
// to do that...
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

func pushStack(inStack []string, col int, block string) (outStack []string) {
	outStack = inStack
	outStack[col] = outStack[col] + block
	return outStack
}

// And this is just a helper/interpreter to iterate over instructions:
func execMoveInstruction(inStack []string, n int, f int, t int) (outStack []string) {
	outStack = inStack
	for i := 0; i < n; i++ {
		crate, tStack := popStack(outStack, f)
		outStack = pushStack(tStack, t, crate)
		//fmt.Println(outStack)
	}
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

	fmt.Printf("Final Stack State:\n%v\n", stack)
	topvals := ""
	var c string
	for i := range stack {
		c, stack = popStack(stack, i)
		topvals = topvals + c
	}
	fmt.Printf("Top Stack Crates after %d moves: %v\n", mvCtr, topvals)
}
