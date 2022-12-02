package main

/*
aoc2022 - Day 2, Part 2
-----------------------
Rock, Paper, Scissors - The Strategy Game
Elaboration: Given an INPUT move (ABC), and a game DIRECTION (X=LOSE,Y=DRAW,Z=WIN):
	a) Calculate the required OUTPUT move (Rock,Paper,Scissors)
	b) Calculate the final score
A,X = ROCK
B,Y = PAPER
C,Z = SCISSORS

*/

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

// for testing:
//var dataFile = "data/day2test.txt"

// for stars:
var dataFile = "data/day2input.txt"

var lut = map[string]string{
	"A": "ROCK",
	"B": "PAPER",
	"C": "SCISSORS",
	"X": "LOSE",
	"Y": "DRAW",
	"Z": "WIN",
}

var mScores = map[string]int{
	"ROCK":     1,
	"PAPER":    2,
	"SCISSORS": 3,
}

var beats = map[string]string{
	"ROCK":     "SCISSORS",
	"PAPER":    "ROCK",
	"SCISSORS": "PAPER",
}

var loses = map[string]string{
	"ROCK":     "PAPER",
	"PAPER":    "SCISSORS",
	"SCISSORS": "ROCK",
}

func main() {
	file, err := os.Open(dataFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	s := bufio.NewScanner(file)
	//each move is on a line, so this is line by line processing
	roundCtr := 0
	validCtr := 0
	p1Wins := 0
	draws := 0
	p1Losses := 0
	p0SumMoveScore := 0
	p1SumMoveScore := 0
	for s.Scan() {
		hand := strings.Fields(s.Text())
		roundCtr++
		fmt.Printf("%05d:%v ", roundCtr, hand)
		//translate the moves
		p0Move, ok0 := lut[hand[0]]
		p1Result, ok1 := lut[hand[1]]
		if ok0 && ok1 {
			validCtr++
			var p1Move string
			//What is the required p1Move to obtain p1Result?
			if p1Result == "WIN" {
				p1Move = loses[p0Move]
				p1Wins++
			} else if p1Result == "LOSE" {
				p1Move = beats[p0Move]
				p1Losses++
			} else {
				p1Move = p0Move
				draws++
			}
			p0SumMoveScore += mScores[p0Move]
			p1SumMoveScore += mScores[p1Move]
			fmt.Printf("Game %04d: P0: %v(%d) P1 to %v, draws: %v(%d)\n", validCtr, p0Move, mScores[p0Move], p1Result, p1Move, mScores[p1Move])
		} else {
			fmt.Printf("Invalid Game\n")
			os.Exit(1)
		}
	}
	//Calculate the totals:
	sumGames := p1Wins + draws + p1Losses
	fmt.Printf("rounds: %04d games: %04d results: %04d\n", roundCtr, validCtr, sumGames)
	p0Score := p0SumMoveScore + (draws * 3) + (p1Losses * 6)
	p1Score := p1SumMoveScore + (draws * 3) + (p1Wins * 6)
	fmt.Printf("Player 0 wins: %d, Drawn Games: %d, Player 1 wins: %d\n", p1Losses, draws, p1Wins)
	fmt.Printf("Player 0 final score: %d\n", p0Score)
	fmt.Printf("Player 1 final score: %d\n", p1Score)
}
