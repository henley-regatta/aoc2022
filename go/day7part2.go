package main

/*
aoc2022 - Day 7, Part 2
-----------------------
The Bloated Mobile Device

The filesystem is 70000000 bytes
We need 30000000 to run the update.
Find the smallest directory that, if deleted, would leave enough space to run the update.

see part 1 for details of parsing.
*/

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func pushDir(orgDir, newDir string) string {
	if orgDir == "/" {
		return "/" + newDir
	} else {
		return orgDir + "/" + newDir
	}

}

func popDir(orgDir string) (string, bool) {
	parts := strings.Split(orgDir[1:], "/")
	if parts[0] == "" {
		return "/", true
	} else {
		return "/" + strings.Join(parts[:len(parts)-1], "/"), false
	}

}

type FileInfo struct {
	name string
	dir  string
	size int
}

func grokListing(listingFile string) []FileInfo {
	file, err := os.Open(listingFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	s := bufio.NewScanner(file)

	//A 2-pass algorithm. In the first path, build our File info structure
	//var DiskData []FileInfo
	CurrentDirectory := ""
	var fileData []FileInfo
	for s.Scan() {
		lTok := strings.Split(s.Text(), " ")
		switch lTok[0] {
		case "$":
			switch lTok[1] {
			case "cd": //Change CurrentDirectory
				{
					switch lTok[2] {
					case "/":
						CurrentDirectory = "/"
					case "..":
						CurrentDirectory, _ = popDir(CurrentDirectory)
					default:
						CurrentDirectory = pushDir(CurrentDirectory, lTok[2])
					}
				}
			default: //NOP - ignore other commands they have no relevance.
			}
		case "dir": //NOP

		default: //It's a file.
			size, _ := strconv.Atoi(lTok[0])
			fStat := FileInfo{name: lTok[1], dir: CurrentDirectory, size: size}
			fileData = append(fileData, fStat)
		}
	}
	return fileData
}

func crtOrAppendDirSize(dirInfo map[string]int, dirName string, addSize int) map[string]int {
	if size, exists := dirInfo[dirName]; exists {
		dirInfo[dirName] = size + addSize
	} else {
		dirInfo[dirName] = addSize
	}
	return dirInfo
}

func getDirSizes(fileData []FileInfo) map[string]int {
	dirInfo := map[string]int{}
	for _, f := range fileData {
		//Ok so THIS dir definitely exists, so add it:
		dirInfo = crtOrAppendDirSize(dirInfo, f.dir, f.size)
		///But we need to iterate up the tree to the root to add to parent:
		cDir := f.dir
		var root bool
		for {
			if cDir, root = popDir(cDir); root {
				break
			} else {
				dirInfo = crtOrAppendDirSize(dirInfo, cDir, f.size)
			}
		}
	}
	return dirInfo
}

const DISKSIZE = 70000000
const UPDATESIZE = 30000000

// for testing:
//var dataFile = "data/day7test.txt"

// for Stars:
var dataFile = "data/day7input.txt"

func main() {
	fData := grokListing(dataFile)
	dirSizes := getDirSizes(fData)
	//Thankfully the way we constructed dirSizes means just one more linear scan required
	diskFree := DISKSIZE - dirSizes["/"]
	minSpaceToFree := UPDATESIZE - diskFree
	fmt.Printf("Currently have %d bytes free; need additional %d for update\n", diskFree, minSpaceToFree)
	minCandSize := DISKSIZE //we know it can't be bigger than this, right?
	dirToDelete := "/"      // I mean, this is the worst case, right?
	for cDir := range dirSizes {
		if dirSizes[cDir] > minSpaceToFree && dirSizes[cDir] < minCandSize {
			minCandSize = dirSizes[cDir]
			dirToDelete = cDir
		}
	}
	fmt.Printf("Smallest Single Dir to Delete is: %v which has size: %d\n", dirToDelete, minCandSize)
}
