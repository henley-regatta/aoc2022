package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"regexp"
	"strconv"
	"strings"
)

/*
aoc2022 - Day 16, Part 1
------------------------
The Elephant Volcano Conundrum

Given an input of valves with associated flow rates AND connectivity to other valves,
with 30 minutes to eruption and a 1-minute delay between moves (open valve OR move to new
location), calculate the maximum "pressure released" (flow-rate * time remaining) possible.


This looks a lot like modified-Dijkstra with a variable "cost" calculated as flow-units...
(if per-unit "cost" is calculated at visit-time, I think this works)

*/

// testing
//var dataFile = "data/day16test.txt"

// stars
var dataFile = "data/day16input.txt"

// How long we have before the explosion
const TIMELIMIT = 30

// eh could be worse
type valveInfo struct {
	f int            //flowrate
	c []string       //reference to other valves
	d map[string]int //distances to other valves from here
}

// This is a static, build-once-reference-many structure
// (in the new implementation)
var vInfo = map[string]valveInfo{}

type djik struct {
	c int
	p string
}

type cState struct {
	pos     string
	t       int
	nfv     int
	vState  map[string]bool
	history []string
}

// been bitten by this before
func cloneState(orgState cState) cState {

	vStClone := map[string]bool{}
	for v := range orgState.vState {
		vStClone[v] = orgState.vState[v]
	}
	hist := []string{}
	for _, h := range orgState.history {
		hist = append(hist, h)
	}
	return cState{
		pos:     orgState.pos,
		t:       orgState.t,
		nfv:     orgState.nfv,
		vState:  vStClone,
		history: hist,
	}
}

func fmtvState(vS map[string]bool) string {
	vStr := []string{}
	for v := range vS {
		n := v
		if !vS[v] {
			n = strings.ToLower(n)
		}
		vStr = append(vStr, n)
	}
	return strings.Join(vStr, ",")
}

func grokValveFromLine(tline string) (string, valveInfo) {
	w := strings.Split(tline, " ")
	valveLabel := w[1]
	fParse := regexp.MustCompile("([0-9]+)")
	fRes := fParse.FindString(tline)
	flowRate, err := strconv.Atoi(fRes)
	if err != nil {
		fmt.Println(w)
		fmt.Println(fRes)
		log.Fatal("could not convert to integer (flow rate)")
	}
	vList := []string{}
	_, val, found := strings.Cut(tline, "leads to valve ")
	if found {
		vList = []string{val}
	} else {
		vl := strings.Split(tline, "lead to valves ")
		vList = strings.Split(vl[1], ", ")

	}
	//fmt.Printf("valve: %s flow rate %d connects to %v\n", valveLabel, flowRate, vList)
	return valveLabel, valveInfo{f: flowRate, c: vList}
}

// against everything I was taught, build a Global Variable
func parseInputForValveList(fromFile string) {

	file, err := os.Open(fromFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	//master loop
	s := bufio.NewScanner(file)

	for s.Scan() {
		vLabel, vStats := grokValveFromLine(s.Text())
		vInfo[vLabel] = vStats
	}

	return
}

func strListContains(needle string, haystack []string) bool {
	for _, v := range haystack {
		if needle == v {
			return true
		}
	}
	return false
}

func dijkfindCheapestUnvisited(nlist map[string]djik) (string, bool) {
	minCost := math.MaxInt
	cheapest := ""
	for n := range nlist {
		if nlist[n].c < minCost {
			minCost = nlist[n].c
			cheapest = n
		}
	}
	if minCost < math.MaxInt {
		return cheapest, true
	} else {
		return cheapest, false
	}
}

func dijkDist(from, to string) int {
	if to == from {
		return 0
	} else if strListContains(to, vInfo[from].c) {
		return 1
	}

	//full dijkstra now required
	unvis := map[string]djik{}
	vis := map[string]djik{}
	for v := range vInfo {
		unvis[v] = djik{c: math.MaxInt}
	}
	initial := djik{c: 0}
	unvis[from] = initial

	for len(unvis) > 0 {
		h, ok := dijkfindCheapestUnvisited(unvis)
		if !ok {
			fmt.Printf("No reasonable unvis left\n")
			log.Fatal("Dijkstra failed me")
		} else {
			for _, n := range vInfo[h].c {
				if _, ok := unvis[n]; ok { //only update if still in unvisited
					upd := djik{c: unvis[h].c + 1, p: h}
					unvis[n] = upd
				}
			}
			//fickin' references will kill me
			vis[h] = djik{c: unvis[h].c, p: unvis[h].p}
			delete(unvis, h)
		}
	}
	//if Dijsktra has worked, all we need is the dist
	//for the end node
	if d, ok := vis[to]; ok {
		return d.c
	}
	return math.MaxInt
}

// pure laziness
func costTo(from, to string) int {
	return vInfo[from].d[to]
}

func posMovesFrom(curr cState) []string {
	vMoveList := []string{}
	for v := range curr.vState {
		if !curr.vState[v] {
			vMoveList = append(vMoveList, v)
		}
	}
	return vMoveList
}

func bestPathFrom(curr cState) (int, []string) {
	bestNFV := curr.nfv
	bestPath := curr.history
	possMoves := posMovesFrom(curr)
	//terminal
	if curr.t < TIMELIMIT {
		for _, pm := range possMoves {
			possNFV := curr.nfv
			tNext := curr.t + costTo(curr.pos, pm)
			if tNext < TIMELIMIT {
				//We have time to flip the valve
				possNFV += vInfo[pm].f * (TIMELIMIT - tNext)
				tNext++ //nb: I thought you did this FIRST but the times only work out if you do it last
			}
			//Now update possNFV by recursing through all possible sub-nodes
			subState := cloneState(curr)
			subState.pos = pm
			subState.nfv = possNFV
			subState.t = tNext
			subState.history = append(subState.history, fmt.Sprintf("%s(%02d)", pm, tNext))
			subState.vState[pm] = true
			possPath := []string{}
			possNFV, possPath = bestPathFrom(subState)
			if possNFV > bestNFV {
				bestNFV = possNFV
				bestPath = possPath
			}
		}
	} /*else {
		fmt.Printf("%04d : %s\n", bestNFV, strings.Join(bestPath, ","))
	}*/
	return bestNFV, bestPath
}

func main() {
	parseInputForValveList(dataFile)

	//Work out shortest distances between all valve pairs
	//also stash the number of valves turned on (zero-flow valves)
	//as these are just chaff and help reduce the special cases later
	//and build an initial list of valvestates to help
	numValves := len(vInfo)
	numOnValves := 0
	valveStates := map[string]bool{}
	for f := range vInfo {
		//Find and stash the costs to each node from here
		dList := map[string]int{}
		for t := range vInfo {
			dList[t] = dijkDist(f, t)
		}
		//stash. terribly complicated in go
		if newv, ok := vInfo[f]; ok {
			newv.d = dList
			vInfo[f] = newv
		}
		//Now mark whether this is a productive valve or not
		//and pre-flip it if not
		if vInfo[f].f == 0 {
			valveStates[f] = true
			numOnValves++
		} else {
			valveStates[f] = false
		}
	}
	fmt.Printf("Working with a set of %d valves, %d of which are zero-flow: %s\n", numValves, numOnValves, fmtvState(valveStates))

	initialState := cState{
		t:       1,
		pos:     "AA",
		nfv:     0,
		vState:  valveStates,
		history: []string{"AA(01)"},
	}

	//There's a recursive solution here which focuses only on travel distances (t-incr)
	//and net future value. But it only works if it can reject local maxima (not *just*
	//leaping from highest value to highest value) to find a global instead.
	bestNFV, bestPath := bestPathFrom(initialState)

	fmt.Println(strings.Join(bestPath, ","))
	fmt.Printf("Best future value that can be found is: %d\n", bestNFV)

}
