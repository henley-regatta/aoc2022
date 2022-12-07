package main

/*
aoc2022 - Day 7, Part 1
-----------------------
The Bloated Mobile Device

Calculate the size of all directories, summing-up sub-directories
Find all directories with a total size of AT MOST 100000, and calculate
the sum of their total sizes.

In the output, lines starting with "$" are a command. Lines starting "dir" reference a sub-directory.
Lines starting with a number are a file-size.

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

const DIRSIZE = 100000

// for testing:
//var dataFile = "data/day7test.txt"

// for Stars:
var dataFile = "data/day7input.txt"

func main() {
	fData := grokListing(dataFile)
	dirSizes := getDirSizes(fData)
	//Now we just need to iterate over dirSizes map and extract dirs < sizeThresh
	sumSmolDirs := 0
	for d := range dirSizes {
		if dirSizes[d] <= DIRSIZE {
			fmt.Printf("%v : %06d\n", d, dirSizes[d])
			sumSmolDirs += dirSizes[d]
		}
	}
	fmt.Printf("Sum of Dirs Under %v is: %d\n", DIRSIZE, sumSmolDirs)
}
