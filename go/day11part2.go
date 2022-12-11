package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

/*
aoc2022 - Day 11, Part 2
------------------------
Monkey Keep-Away. A fantastically intricate keepy-uppy game of
item value change and position rotation, designed to pose a parsing
challenge as well as an iterative gameplay test of skills.

This time around, the fear level doesn't decrease (no /3) on eval, and
we need to play FOR 10,000 ROUNDZ

Return the multiple of the two highest values as the Monkey Business value.

*/

// These are the only arithmetic operations seen in the input
type Operation int

const (
	add Operation = iota
	mult
	square
)

// helper func to dump operation
func (o Operation) String() string {
	switch o {
	case add:
		return "Add"
	case mult:
		return "Multiply"
	case square:
		return "Squared"
	}
	return "UNDEFINED"
}

// The game is driven by operations performed by processing elements - monkeys
// nb: they're defined as indexes starting from zero so an array of Monke is good enough
//
//	for addressing purposes
type Monke struct {
	iCount  int       //used as part of the final analysis, tracked throughout the iterations
	op      Operation //See Above.
	operand int       //not used for "square" function but, eh, needed otherwise
	divisor int       //the post-operation test.
	iftrue  int       //monkeys to pass to on true, false
	iffalse int
	holding []int //a list of items the monkey is holding. we start with index numbers but these get modified by operation so they're just ints.
}

// I KNOW this isn't save but it is easy
func qAtoi(str string) int {
	v, e := strconv.Atoi(str)
	if e == nil {
		return v
	}
	return 0
}

// This is easy. Honest.
func extractMonkey(lines []string) Monke {
	var monkey Monke
	//starting items is in lines[1]
	items := strings.Split(lines[1], ": ")
	iBuf := strings.Split(items[1], ", ")
	for _, v := range iBuf {
		monkey.holding = append(monkey.holding, qAtoi(v))
	}
	//Operation is in lines[2]
	opbits := strings.Split(lines[2], "= ")
	//now a bit of work required..
	op := strings.Split(opbits[1], " ")
	if op[1] == "*" {
		if op[2] == "old" {
			monkey.op = square
			monkey.operand = 0
		} else {
			monkey.op = mult
			monkey.operand = qAtoi(op[2])
		}
	} else {
		monkey.op = add
		monkey.operand = qAtoi(op[2])
	}

	//Divisor is in lines[3]
	div := strings.Split(lines[3], "divisible by ")
	monkey.divisor = qAtoi(div[1])
	//and the rest is the same...
	truebits := strings.Split(lines[4], "throw to monkey ")
	monkey.iftrue = qAtoi(truebits[1])
	falsebits := strings.Split(lines[5], "throw to monkey ")
	monkey.iffalse = qAtoi(falsebits[1])

	return monkey
}

func parseDataFile(dFile string) []Monke {
	file, err := os.Open(dFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var monkeys []Monke //what we'll build

	//Since input is blank-line delimited it's easiest to build
	//records based on non-blank lines
	s := bufio.NewScanner(file)
	var mBuf []string
	//master loop
	for s.Scan() {
		line := s.Text()
		if len(line) == 0 {
			monkeys = append(monkeys, extractMonkey(mBuf))
			mBuf = nil
			//mBuf[len(mBuf)-1] = nil //blank the struct out
		} else {
			mBuf = append(mBuf, line)
		}
	}
	//we're always left with a record in the end, n'est ce pas?
	monkeys = append(monkeys, extractMonkey(mBuf))
	return monkeys
}

// Although the text presents this as a 2-part operation
// internally it's atomic - happens before anything else
// HOWEVER in part 2 the fear level *quickly* exceeds any reasonable
// value (way outside uint64 range) so we have to deal with overflows.
func changeFearLevel(inFear, lcf, operand int, op Operation) int {
	outFear := inFear
	//First we INCREASE the fear level:
	switch op {
	case square:
		outFear = inFear * inFear
	case mult:
		outFear = inFear * operand
	case add:
		outFear = inFear + operand
	}

	//We only need the remainder of this value against the Lowest Common Factor
	outFear = outFear % lcf

	return outFear
}

// This is probably overkill but it neatens the code
func testPass(fear int, divisor int) bool {
	if fear%divisor == 0 {
		return true
	}
	return false
}

// Master function to take a Monkey's turn
func takeTurn(m, lcf int, monkeys []Monke) []Monke {
	//because monkeys is complex it's copied-by-reference:
	outMonkeys := monkeys
	//copy the current holding list before we start
	var thisHolding []int
	for i := range monkeys[m].holding {
		thisHolding = append(thisHolding, monkeys[m].holding[i])
	}
	//now reset the output list:
	outMonkeys[m].holding = nil
	//FOR EACH ITEM held by the monkey:
	for i := range thisHolding {
		//POP the item off the slice
		item := thisHolding[i]
		//Process the item
		iVal := changeFearLevel(item, lcf, monkeys[m].operand, monkeys[m].op)
		passTo := monkeys[m].iffalse
		if testPass(iVal, monkeys[m].divisor) {
			passTo = monkeys[m].iftrue
		}
		//fmt.Printf("Item %d changes to %d, passes to monkey %d\n", item, iVal, passTo)
		outMonkeys[passTo].holding = append(outMonkeys[passTo].holding, iVal)
		outMonkeys[m].iCount++ //Monkey has finished evaluating this item
	}
	return outMonkeys
}

// testing
//var dataFile = "data/day11test.txt"

// stars
var dataFile = "data/day11input.txt"

func main() {
	monkeys := parseDataFile(dataFile)

	//We can't store "fear level" it grows so large. Instead we can just store the
	//REMAINDER against the lowest common factor of the divisors and use that.
	//But we need to calculate that...
	lcf := 1
	for m := range monkeys {
		lcf = lcf * monkeys[m].divisor
	}
	fmt.Printf("Calculated lowest common factor as %d\n", lcf)

	//WE HAVE REACHED THE STARTING LINE. Time to implement the game algorithm....
	//This is now a 10,000 turn game
	for t := 1; t < 10001; t++ {
		//Iterate over the monkeys taking each one in term
		for m := range monkeys {
			monkeys = takeTurn(m, lcf, monkeys)
		}
		//All we really care about is the number of items evaluated per Monkey
		if t%1000 == 0 || t == 20 {
			fmt.Printf("Round %d\n", t)
			for m := range monkeys {
				fmt.Printf("\tMonkey %d evaluated %d items\n", m, monkeys[m].iCount)
			}
		}

	}
	//At the end of the game we need to work out which monkey evaluated the most
	topevals := [2]int{0, 0}
	for m := range monkeys {
		if monkeys[m].iCount > topevals[0] {
			topevals[1] = topevals[0]
			topevals[0] = monkeys[m].iCount
		} else if monkeys[m].iCount > topevals[1] {
			topevals[1] = monkeys[m].iCount
		}
		fmt.Printf("Monkey %d looked at %d items\n", m, monkeys[m].iCount)
	}
	//Eval score is the product of the top 2:
	fmt.Printf("Monkey Business Score is: %d\n", topevals[0]*topevals[1])
}
