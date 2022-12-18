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
	"time"
)

/*
aoc2022 - Day 16, Part 2
------------------------
The Elephant Volcano Conundrum

Given an input of valves with associated flow rates AND connectivity to other valves,
with 26 minutes to eruption and a 1-minute delay between moves (open valve OR move to new
location), calculate the maximum "pressure released" (flow-rate * time remaining) possible.

This is the same as part one EXCEPT we have 2 movers (apparently we've taught an elephant
to open valves). Hence losing 4 minutes off the front.

*/

// testing
//var dataFile = "data/day16test.txt"

// stars
var dataFile = "data/day16input.txt"

// How long we have before the explosion
const TIMELIMIT = 26

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
	pos       string
	t         int
	nfv       int
	offValves []string
	history   []string
}

// been bitten by this before
func cloneState(orgState cState) cState {

	offV := []string{}
	for _, v := range orgState.offValves {
		offV = append(offV, v)
	}
	hist := []string{}
	for _, h := range orgState.history {
		hist = append(hist, h)
	}
	return cState{
		pos:       orgState.pos,
		t:         orgState.t,
		nfv:       orgState.nfv,
		offValves: offV,
		history:   hist,
	}
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

func deleteFrom(strList []string, deleteme string) []string {
	deletedList := []string{}
	for _, v := range strList {
		if v != deleteme {
			deletedList = append(deletedList, v)
		}
	}
	return deletedList
}

// pure laziness
func costTo(from, to string) int {
	return vInfo[from].d[to]
}

func bestPathFrom(curr cState) (int, []string) {
	bestNFV := curr.nfv
	bestPath := curr.history
	possMoves := curr.offValves
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
			subState.offValves = deleteFrom(curr.offValves, pm)
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

// works on the global. Is this good practice? Naa.
func genDistances(from string) {
	//Find and stash the costs to each node from here
	dList := map[string]int{}
	for t := range vInfo {
		dList[t] = dijkDist(from, t)
	}
	//stash. terribly complicated in go
	if newv, ok := vInfo[from]; ok {
		newv.d = dList
		vInfo[from] = newv
	}
}

// Generate all possible lists of size pSize from an input set inList
// TODO: Fixme. Can't work out how to
// see: https://www.geeksforgeeks.org/print-subsets-given-size-set/
// but see better: https://cloud.tencent.com/developer/article/1412871
var subsets [][]string

func genStrSliceSubsets(inList []string, pSize int) [][]string {
	subsets = [][]string{}
	cur := []string{}
	for i := 1; i <= pSize; i++ {
		subDFS(inList, i, 0, cur)
	}
	return subsets
}
func subDFS(inList []string, n, s int, cur []string) {
	if len(cur) == n {
		tmp := make([]string, n)
		copy(tmp, cur)
		subsets = append(subsets, tmp)
		return
	}
	for i := s; i < len(inList); i++ {
		cur = append(cur, inList[i])
		subDFS(inList, n, i+1, cur)
		cur = cur[:len(cur)-1]
	}
}

func main() {
	parseInputForValveList(dataFile)

	//special case - need all distances from start node
	genDistances("AA")

	offValves := []string{}
	for f := range vInfo {
		//only need info if it's a productive valve
		if vInfo[f].f > 0 {
			offValves = append(offValves, f)
			genDistances(f)
		}
	}
	fmt.Printf("Working with a set of %d valves, %d of which are non-zero-flow: %s\n", len(vInfo), len(offValves), strings.Join(offValves, ","))

	//At this point I need to generate all subsets of splitting offValves in two
	//This is probably easiest with a call that generates all possible half-sets.
	personSubsets := genStrSliceSubsets(offValves, len(offValves)/2)

	bestScore := 0
	bestpPath := []string{"AA(01)"}
	bestePath := []string{"AA(01)"}

	tStart := time.Now()
	compRequired := len(personSubsets)
	compPerformed := 0
	fmt.Printf("Comparing %d sets of paths\n", compRequired)
	//Now find the best outcome for each of these subset pairs
	for _, pSet := range personSubsets {
		//generate the elephant subset
		eSet := []string{}
		for _, v := range offValves {
			if !strListContains(v, pSet) {
				eSet = append(eSet, v)
			}
		}

		//Find best path for the person
		pState := cState{
			t:         1,
			pos:       "AA",
			nfv:       0,
			offValves: pSet,
			history:   []string{"AA(01)"},
		}
		pScore, pPath := bestPathFrom(pState)
		//Find best path for the elephant
		eState := cloneState(pState)
		eState.offValves = eSet
		eScore, ePath := bestPathFrom(eState)

		//Total score is the sum of the individual scores
		totScore := pScore + eScore
		if totScore > bestScore {
			bestScore = totScore
			bestpPath = pPath
			bestePath = ePath
		}

		compPerformed++
		if compPerformed%100 == 0 {
			pctComplete := compPerformed * 100 / compRequired
			elapsed := time.Since(tStart)
			remaining := elapsed.Seconds() * (100.0 / float64(pctComplete))
			fmt.Printf("Completed %d (%d%%) of %d comparisons; best score yet is %d. Estimate %.1f seconds remaining\n", compPerformed, pctComplete, compRequired, bestScore, remaining)
		}
	}

	fmt.Printf("P took: %s\n", strings.Join(bestpPath, ","))
	fmt.Printf("E took: %s\n", strings.Join(bestePath, ","))
	fmt.Printf("Best future value that can be found is: %d\n", bestScore)
}
