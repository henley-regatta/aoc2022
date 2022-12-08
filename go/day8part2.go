package main

/*
aoc2022 - Day 8, Part 2
-----------------------
Treetop Tree House Shenanigans

Using an x*y single-digit map of tree heights (0 lowest, 9 highest) determine the scenic
viewability of a tree.
Visibility restricted to cardinal directions (N,S,E,W). (so that wasn't the elaboration...)

A tree's scenic score is the product of the sum of trees visible from it to the nearest height match in
the cardinal directions.

Find the highest scenic score for a tree in the grid.
*/

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

// testing
//var dataFile = "data/day8test.txt"

// live
var dataFile = "data/day8input.txt"

func parseFileToNumberGrid(fname string) [][]int {
	file, err := os.Open(fname)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	s := bufio.NewScanner(file)

	//Input is line-oriented.
	var treeMap [][]int
	for s.Scan() {
		//Convert the input line to runes
		rLine := []rune(s.Text())
		//Since we know the input is numeric, just subtract x30 to get the digit value
		var tLine []int
		for i := range rLine {
			tLine = append(tLine, int(rLine[i])-0x30)
		}
		treeMap = append(treeMap, tLine)
	}
	return treeMap
}

// Instead of is-it-visible, calculate the visibility score across
// the cardinal points
// criteria given (cardinal points only)
func calcVisibilityScore(treeGrid [][]int, x, y int) int {
	cHeight := treeGrid[y][x]
	//Check needs to be in each cardinal direction, so:
	visMap := map[string]int{"n": 1, "s": 1, "e": 1, "w": 1}
	for i := (x - 1); i > 0; i-- {
		if cHeight > treeGrid[y][i] {
			visMap["w"]++
		} else {
			break
		}
	}
	for i := x + 1; i < len(treeGrid[y])-1; i++ {
		if cHeight > treeGrid[y][i] {
			visMap["e"]++
		} else {
			break
		}
	}
	for i := y - 1; i > 0; i-- {
		if cHeight > treeGrid[i][x] {
			visMap["n"]++
		} else {
			break
		}
	}
	for i := y + 1; i < len(treeGrid)-1; i++ {
		if cHeight > treeGrid[i][x] {
			visMap["s"]++
		} else {
			break
		}
	}

	visScore := visMap["e"] * visMap["w"] * visMap["n"] * visMap["s"]
	//fmt.Printf("[%d,%d](%v) score=%d (%v)\n", x, y, cHeight, visScore, visMap)
	return visScore
}

func main() {
	treeGrid := parseFileToNumberGrid(dataFile)
	fmt.Printf("Dealing with an %d x %d tree grid\n", len(treeGrid[0]), len(treeGrid))
	bestVisibility := 0
	bestLocation := [2]int{0, 0}
	//Iterate over the grid - ignoring the outermost border which is known to be visible
	//determining the score of each element
	for y := 1; y < len(treeGrid)-1; y++ {
		for x := 1; x < len(treeGrid[y])-1; x++ {
			visScore := calcVisibilityScore(treeGrid, x, y)
			if visScore > bestVisibility {
				bestVisibility = visScore
				bestLocation[0] = x
				bestLocation[1] = y
			}
		}
	}
	fmt.Printf("The tree at %v had the best visibility score of %d\n", bestLocation, bestVisibility)
}
