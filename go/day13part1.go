package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
)

/*
aoc2022 - Day 13, Part 1
-----------------------
Elven Distress Signal Packet Decoding - Recursive List Comprehension Edition

Given a list of pairs of brackets (possibly nested), identify
how many pairs of packets are in the right order.

*/

// for testing:
//var dataFile = "data/day13test.txt"

//var dataFile = "data/day13test2.txt"

// for Stars:
var dataFile = "data/day13input.txt"

type IorL int

const (
	n IorL = iota
	l
)

type intOrList struct {
	t IorL
	v int
	l []intOrList
}

func (t IorL) String() string {
	switch t {
	case n:
		return "Number"
	case l:
		return "List"
	}
	return "Undefined"
}

type triState int

const (
	unknown triState = iota
	fail
	pass
)

func (s triState) String() string {
	switch s {
	case unknown:
		return "Continue"
	case fail:
		return "UNORDERED"
	case pass:
		return "ORDERED"
	}
	return "Undefined"
}

// I need this mainly to check I've decoded the thing correctly
func sprintList(list []intOrList) string {
	outstring := "["
	for len(list) > 0 {
		e := list[0]
		list = list[1:]
		if e.t == n {
			outstring = outstring + strconv.Itoa(e.v) + ","
		} else {
			outstring = outstring + sprintList(e.l) + ","
		}
	}
	if outstring[len(outstring)-1] == ',' {
		outstring = string(outstring[0 : len(outstring)-1])
	}
	return outstring + "]"
}

func unpackList(lStr string) (string, []intOrList) {
	oList := []intOrList{}
	//I need to parse char-by-char BUT I also need to account for recursion
	nBuf := ""
	for len(lStr) > 0 {
		c := rune(lStr[0])
		lStr = lStr[1:]
		switch c {
		case ',':
			{
				if len(nBuf) > 0 {
					nVal, _ := strconv.Atoi(nBuf)
					oList = append(oList, intOrList{t: n, v: nVal})
				}
				nBuf = ""
			}
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			nBuf += string(c)
		case '[':
			{
				remStr, subList := unpackList(lStr)
				oList = append(oList, intOrList{t: l, l: subList})
				lStr = remStr
			}
		case ']':
			{
				if len(nBuf) > 0 {
					nVal, _ := strconv.Atoi(nBuf)
					oList = append(oList, intOrList{t: n, v: nVal})
				}
				return lStr, oList
			}
		default:
			log.Fatal("Unrecognised input char: " + string(c))
		}
	}
	return lStr, oList
}

func isListPairOrdered(left, right []intOrList) triState {
	//Element by Element analysis required
	for len(left) > 0 && len(right) > 0 {
		l1 := left[0]
		r1 := right[0]
		left = left[1:]
		right = right[1:]
		//Evaluation cases:
		fmt.Printf("\tCMP: %v vs. %v\n\t%v\n\t%v\n", l1.t, r1.t, sprintList([]intOrList{l1}), sprintList([]intOrList{r1}))
		if l1.t == n && r1.t == n {
			fmt.Printf("\t\tVALCOMP: %d < %d ?\n", l1.v, r1.v)
			if l1.v < r1.v {
				return pass
			} else if l1.v > r1.v {
				return fail
			}
			fmt.Printf("\t\tVALCOMP CONTINUE\n")
			//deliberate: if equal continue
		}
		if l1.t == l && r1.t == l {
			fmt.Printf("\t\tRECURSE LIST COMPARE\t%v\tvs.\t%v\n", sprintList(l1.l), sprintList(r1.l))
			sListVal := isListPairOrdered(l1.l, r1.l)
			if sListVal != unknown {
				return sListVal
			}
			//OTHERWISE continue
		}
		if l1.t == n && r1.t == l {
			fmt.Printf("\t\tCOERCE LEFT TO LIST AND RECURSE\n")
			lAsList := intOrList{t: l, l: []intOrList{l1}}
			sListVal := isListPairOrdered([]intOrList{lAsList}, []intOrList{r1})
			if sListVal != unknown {
				return sListVal
			}
		}
		if l1.t == l && r1.t == n {
			fmt.Printf("\t\tCOERCE RIGHT TO LIST AND RECURSE\n")
			rAsList := intOrList{t: l, l: []intOrList{r1}}
			sListVal := isListPairOrdered([]intOrList{l1}, []intOrList{rAsList})
			if sListVal != unknown {
				return sListVal
			}
		}

	}
	//We've run out of stuff to evaluate. BUT
	//do we have elements left in either list?
	if len(left) > 0 {
		fmt.Printf("\t\tRIGHT RAN OUT OF ITEMS FIRST\n")
		return fail //specification: left should run out of elements before right
	} else if len(right) > 0 {
		fmt.Printf("\t\tLEFT RAN OUT OF ITEMS FIRST\n")
		return pass
	} else {
		//both must be empty; this is undefined but t'internetz sez:
		return unknown
	}

}

func isCorrectlyOrdered(pStr []string) bool {

	_, left := unpackList(pStr[0])
	_, right := unpackList(pStr[1])

	left = left[0].l
	right = right[0].l

	//validate our parsing is correct
	if pStr[0] != sprintList(left) {
		fmt.Printf("%s\n%s\n", pStr[0], sprintList(left))
		log.Fatal("left didn't parse")
	}

	if pStr[1] != sprintList(right) {
		fmt.Printf("%s\n%s\n", pStr[1], sprintList(right))
		log.Fatal("right didn't parse")
	}

	fmt.Printf("EVALUATE:\n\t%s\n\t%s\n", sprintList(left), sprintList(right))
	isCorrectlyOrdered := isListPairOrdered(left, right)
	fmt.Printf("EVALUATION ISORDERED: %v\n", isCorrectlyOrdered)

	//NOW FOR THE EVALUATION....
	if isCorrectlyOrdered == pass {
		return true
	}
	return false
}

func main() {
	file, err := os.Open(dataFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	s := bufio.NewScanner(file)

	//We compare 2 lines at a time (2 packets)
	var pPair []string
	ndxCount := 1
	sumOrderedCount := 0
	listCorrectly := []int{}
	for s.Scan() {
		p := s.Text()
		if len(p) == 0 {
			//DO THE PACKET PARSE THING
			fmt.Printf("-- PACKET %04d COMPARISON ---------------------------------------------------------\n", ndxCount)
			if isCorrectlyOrdered(pPair) {
				sumOrderedCount += ndxCount
				listCorrectly = append(listCorrectly, ndxCount)
			}
			//Go again
			pPair = nil
			ndxCount++
		} else {
			pPair = append(pPair, p)
		}
	}
	//There's always one left over...
	fmt.Printf("-- PACKET %04d COMPARISON ---------------------------------------------------------\n", ndxCount)
	if isCorrectlyOrdered(pPair) {
		sumOrderedCount += ndxCount
		listCorrectly = append(listCorrectly, ndxCount)
	}

	fmt.Println(listCorrectly)
	fmt.Printf("After %d packet pairs, Sum of Indicies of Correctly Ordered Pairs: %d\n", ndxCount, sumOrderedCount)
}
