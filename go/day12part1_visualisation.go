package main

import (
	"bufio"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"math"
	"os"
)

/*
aoc2022 - Day 12, Part 1
------------------------
The Height-Transfer Mobile Signal Finding Day

WITH VISUALISATION

Stitch the generated files together using:
	ffmpeg -f image2 -framerate 10 -i visualisation/aoc2022day12part1_%04d.png -c:v libvpx-vp9 -pix_fmt yuva420p visualisation/aoc2022_day12part1.webm


Given an input heightmap encoded a(lowest)-z(highest),
navigate from point S to point E in the shortest number of steps possible.
Each (cardinal) move can vary by at most one step in height difference.

This smells a lot like Dijkstra, so let's cut to the chase...

*/

// testing
//var dataFile = "data/day12test.txt"

// stars
var dataFile = "data/day12input.txt"

// There's a lot of these in the code so just to keep us straight
type point [2]int

type destInfo struct {
	cost int
	prev point
}

// IT'S GRATUITOUS PNG GENERATION TIME
func generateVisualisation(pointmap [][]int, visited map[point]destInfo, route []point, pngNdx int) {
	// FOR DEM OUTPUTZ
	pngPath := "visualisation/"
	outPngFile := fmt.Sprintf("%saoc2022day12part1_%04d.png", pngPath, pngNdx)
	fmt.Printf("Outputting route to %s\n", outPngFile)

	//Generate a canvas scaled against pointmap size
	xScale := 6
	yScale := 12
	rangeX := (len(pointmap[0]) - 1) * xScale
	rangeY := (len(pointmap) - 1) * yScale
	topLeft := image.Point{0, 0}
	botRight := image.Point{rangeX, rangeY}
	outCanvas := image.NewRGBA(image.Rectangle{topLeft, botRight})

	//First up plot a "background" as greyscale; lower = darker
	//input heights known to be in range 1-26 but we don't want them too dark/light
	gRange := 200 / 26
	for y := 0; y < len(pointmap); y++ {
		for x := 0; x < len(pointmap[y]); x++ {
			gScale := uint8(225 - (pointmap[y][x]*gRange + 25))
			for oY := 0; oY < yScale; oY++ {
				for oX := 0; oX < xScale; oX++ {
					outCanvas.Set(x*xScale+oX, y*yScale+oY, color.RGBA{gScale, gScale, gScale, 0xff})
				}
			}
		}
	}

	//Find a "scale" for the heightmap based on the relative cost values (where known)
	maxPathCost := 0
	for p := range visited {
		if visited[p].cost > maxPathCost {
			maxPathCost = visited[p].cost
		}
	}
	//We want to scale "cost" across the map from green(minimum) to red(maximum)
	pathScale := maxPathCost / 255
	for p := range visited {
		for oY := 0; oY < yScale; oY++ {
			for oX := 0; oX < xScale; oX++ {
				green := uint8(255 - pathScale*visited[p].cost)
				red := uint8(pathScale * visited[p].cost)
				outCanvas.Set(p[0]*xScale+oX, p[1]*yScale+oY, color.RGBA{red, green, uint8(0), 0xff})
			}
		}
	}

	//Now plot the actual route over the top
	for p := range route {
		for oY := 0; oY < yScale; oY++ {
			for oX := 0; oX < xScale; oX++ {
				outCanvas.Set(route[p][0]*xScale+oX, route[p][1]*yScale+oY, color.RGBA{0x00, 0x00, 0xff, 0xff})
			}
		}
	}

	//Write the file
	f, _ := os.Create(outPngFile)
	png.Encode(f, outCanvas)
}

// Trick One is parsing the input. And I say "trick" because
// S and E are hiding the underlying heights of their squares.
func parseMap(fromFile string) (point, point, [][]int) {
	file, err := os.Open(fromFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var start point
	var end point
	var heightMap [][]int
	//master loop
	s := bufio.NewScanner(file)
	lCount := 0
	for s.Scan() {
		//Now a char-by-char scan, using Runes, because that makes int conversion trivial
		inLine := []rune(s.Text())
		//this would be a 1-liner if we didn't have to check for start and end
		var mapLine []int
		for i, v := range inLine {
			switch v {
			case 'S':
				{
					start[0] = i
					start[1] = lCount
					mapLine = append(mapLine, int('a')-0x60) //GIVEN IN SPECIFICATION DUMBKOPF
				}
			case 'E':
				{
					end[0] = i
					end[1] = lCount
					mapLine = append(mapLine, int('z')-0x60) //GIVEN IN SPECIFICATION DUMBKOPF
				}
			default:
				{
					mapLine = append(mapLine, int(v)-0x60)
				}
			}
		}

		lCount++
		heightMap = append(heightMap, mapLine)
	}
	return start, end, heightMap

}

// prettifier func
func dumpMap(printMap [][]int) {
	for y := 0; y < len(printMap); y++ {
		outLine := ""
		for x := 0; x < len(printMap[y]); x++ {
			if printMap[y][x] < 0 {
				outLine = outLine + "-"
			} else if printMap[y][x] > 26 {
				outLine = outLine + "+"
			} else {
				outLine = outLine + string(rune(printMap[y][x]+0x60))
			}
		}
		fmt.Println(outLine)
	}
}

// quicky to check heights are reachable between two (assumed adjacent) points
// NB: invalid to CLIMB >1 step, can DESCEND any height change....
func costToReach(from, to point, pointmap [][]int) int {
	fHeight := pointmap[from[1]][from[0]]
	tHeight := pointmap[to[1]][to[0]]

	diff := tHeight - fHeight
	if diff <= 1 {
		return 1
	}
	return math.MaxInt //unreachable

}

// Find nodes reachable from a given point with the cost to reach
// (updating the cost map on return)
func cardinalPoints(fromPoint point, pointmap [][]int) []point {
	var cands []point
	//movement is cardinal only. There are at most 4 moves possible:
	if fromPoint[0] < len(pointmap[fromPoint[1]])-1 {
		cands = append(cands, point{fromPoint[0] + 1, fromPoint[1]})
	}
	if fromPoint[0] > 0 {
		cands = append(cands, point{fromPoint[0] - 1, fromPoint[1]})
	}
	if fromPoint[1] < len(pointmap)-1 {
		cands = append(cands, point{fromPoint[0], fromPoint[1] + 1})
	}
	if fromPoint[1] > 0 {
		cands = append(cands, point{fromPoint[0], fromPoint[1] - 1})
	}
	return cands

}

func findCheapestUnvisited(unvisited map[point]destInfo) (point, bool) {
	//O(n) scan works ok here:
	minCost := math.MaxInt
	minPoint := point{math.MaxInt, math.MaxInt}
	for k := range unvisited {
		//TODO: Should eval be <= not < ?
		if unvisited[k].cost < minCost {
			minCost = unvisited[k].cost
			minPoint = k
		}
	}
	if minCost == math.MaxInt {
		return minPoint, false
	}
	return minPoint, true
}

func genRouteSoFar(visited map[point]destInfo, start, end point) []point {

	//Now our best route is the back-track from the end point to the start
	trackback := []point{end}
	for trackback[len(trackback)-1] != start {
		prev := visited[trackback[len(trackback)-1]].prev
		trackback = append(trackback, prev)
	}
	//note this route is backwards end->start, reverse it before returning
	bestRoute := []point{}
	for _, v := range trackback {
		bestRoute = append([]point{v}, bestRoute...)
	}
	return bestRoute
}

func dijkstra_route(start, end point, pointmap [][]int) ([]point, bool) {

	//Initialise the unvisited list (note: cost from start is 0)
	unvisited := map[point]destInfo{}
	for y := 0; y < len(pointmap); y++ {
		for x := 0; x < len(pointmap[y]); x++ {
			pos := point{x, y}
			if pos == start {
				unvisited[pos] = destInfo{cost: 0, prev: point{math.MaxInt, math.MaxInt}}
			} else {
				unvisited[pos] = destInfo{cost: math.MaxInt, prev: point{math.MaxInt, math.MaxInt}}
			}
		}
	}
	//And the visited list:
	visited := map[point]destInfo{}

	pngNdx := 0

	//REPEAT UNTIL WE EVALUATE EVERYTHING:
	for len(unvisited) > 0 {
		//Get cheapest unvisited node
		evalPoint, foundCheapest := findCheapestUnvisited(unvisited)
		if !foundCheapest {
			//suspect we have islands of unreferenceability
			//fmt.Printf("Aborting end scan with %d unvisited to check\n", len(unvisited))
			break
		}
		//find reachable unvisited nodes
		candNext := cardinalPoints(evalPoint, pointmap)
		//Filter first by checking whether it's on the unvisited list
		fList := []point{}
		for _, cand := range candNext {
			//Filter by previously-visited
			if _, ok := unvisited[cand]; ok {
				fList = append(fList, cand)
			}
		}
		for _, cand := range fList {
			//Cost to visit MAY be prohibitive (invalid height GAIN)
			dCost := costToReach(evalPoint, cand, pointmap)
			if dCost < math.MaxInt {
				tCost := dCost + unvisited[evalPoint].cost
				if tCost < unvisited[cand].cost {
					unvisited[cand] = destInfo{cost: tCost, prev: evalPoint}
				}
			}
		}

		//Move current point to the visited list
		visited[evalPoint] = unvisited[evalPoint]
		delete(unvisited, evalPoint)

		//Generate an output frame every X evals
		if len(visited)%50 == 0 {
			cRoute := genRouteSoFar(visited, start, evalPoint)
			generateVisualisation(pointmap, visited, cRoute, pngNdx)
			pngNdx++
		}

	}

	//Sanity check required. Visited must contain both the "end" and the "start"
	//points
	_, isendvisited := visited[end]
	_, isstartvisited := visited[start]
	if !isendvisited || !isstartvisited {
		fmt.Printf("No route to %v\n", end)
		return []point{}, false
	}

	return genRouteSoFar(visited, start, end), true
}

func main() {
	start, end, heightmap := parseMap(dataFile)
	fmt.Printf("Routing from %+v to %+v...\n", start, end)
	route, routeok := dijkstra_route(start, end, heightmap)
	if routeok {
		fmt.Println(route)
		fmt.Printf("Shortest path steps: %d\n", len(route)-1)
	} else {
		log.Fatal("Oh god the only startpoint available doesn't route to the finish")
	}

}
