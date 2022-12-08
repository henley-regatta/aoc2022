package main

/*
aoc2022 - Day 8, Part 1
-----------------------
Treetop Tree House Shenanigans

Using an x*y single-digit map of tree heights (0 lowest, 9 highest) determine
the visibility of a tree.
Visibility restricted to cardinal directions (N,S,E,W). (guess the Part2 elaboration?)

A tree is visible if all trees between it and the edge are lower than it. All Edge trees are visible.

Count the number of visible trees in the grid.
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

// Determine visibility of tree at [x][y] using the
// criteria given (cardinal points only)
func isVisible(treeGrid [][]int, x, y int) bool {
	cHeight := treeGrid[y][x]
	//Check needs to be in each cardinal direction, so:
	visMap := map[string]bool{"n": true, "s": true, "e": true, "w": true}
	for i := 0; i < x; i++ {
		if treeGrid[y][i] >= cHeight {
			visMap["w"] = false
			break
		}
	}
	for i := x + 1; i < len(treeGrid[y]); i++ {
		if treeGrid[y][i] >= cHeight {
			visMap["e"] = false
			break
		}
	}
	for i := 0; i < y; i++ {
		if treeGrid[i][x] >= cHeight {
			visMap["n"] = false
			break
		}
	}
	for i := y + 1; i < len(treeGrid); i++ {
		if treeGrid[i][x] >= cHeight {
			visMap["s"] = false
			break
		}
	}

	anyVis := false
	for i := range visMap {
		if visMap[i] == true {
			anyVis = true
			break
		}
	}
	return anyVis
}

func main() {
	treeGrid := parseFileToNumberGrid(dataFile)
	fmt.Printf("Dealing with an %d x %d tree grid\n", len(treeGrid[0]), len(treeGrid))
	//Iterate over the grid - ignoring the outermost border which is known to be visible
	//determining the visibility of each element
	visibleCount := len(treeGrid[0]) //first row is always visible
	for y := 1; y < len(treeGrid)-1; y++ {
		visibleCount += 2 //First and Last trees in row are always visible
		for x := 1; x < len(treeGrid[y])-1; x++ {
			if isVisible(treeGrid, x, y) {
				visibleCount++
			}
		}
	}
	visibleCount += len(treeGrid[len(treeGrid)-1]) //Last row is always visible
	fmt.Printf("There are %d visible trees in the grid\n", visibleCount)
}
