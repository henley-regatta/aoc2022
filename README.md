# aoc2022 

[2022 Advent Of Code](https://adventofcode/2022) - My solutions

I really, really am not going to take this seriously. This is just a bit of fun. 
This year I'm just going to pick at it in _other_ languages. Not just Python.

## Go
  - `go/day1part1.go` - seemed obvious to dynamically create an int map incrementing "key" (elf index) on a blank line (_i.e._ one that can't be `str.AtoI` converted)
  - `go/day1part2.go` - Let's not dignify this with "efficient" as a characteristic but sorting the previous output by "value" (per-elf calorie count) instead of key seemed obvious, and then taking a slice of the top 3 entries and summing them was the natural conclusion
  - `go/day2part1.go` - Trying to be clever and only enumerating what I saw in the input didn't work, the implication being that I missed some game outcomes. Translating everything to Rock/Paper/Scissors and truth-tabling the game outcomes got me an actual result. Day 2 and already I'm into Voodoo coding....
  - `go/day2part2.go` - Ah the first example of "If you do Part 1 _this_ way, Part 2 becomes easy. Just translate the instuctions and reverse the look-ups and, boomshanka, the desired result.
  - `go/day3part1.go` - A little light substring scanning and range conversion. Clearly part 2 is going to make this much harder...
  - `go/day3part2.go` - AHHAHAAA! a little extra scaffolding to do 3-way instead of 2-way matches by whole lines but I'd basically cracked this in part 1.
  - `go/day4part1.go` - A suspiciously easy answer which (as t'internet notes) 90% of the effort is the input string parsing and only 10% on actual problem logic...
  - `go/day4part2.go` - OK, change 1 comparison condition and we get ANY overlap not just strict overlap. More time renaming variables than writing code here...
  - `go/day5part1.go` - Nothing wrong with a bit of stack-manipulation is there? Spend an inordinate amount of time on string-slicing / rune-converting but that's just me not understanding Go very well. Also had to re-do the stack output in order to see that I'd made an error transcribing my input stack which, er, meant the wrong result. Error in Fingers not Code!
  - `go/day5part2.go` - Scratch the stack-smashing, return to string-slicing. Conceptually an easier problem than Part 1, only complicated by my inability to correctly apply indexes on the first attempt. 
  - `go/day6part1.go` - Yay bit of string-scanning. Opportunity for a bit of modulus arithmetic, truth telling and functional parameterisation too.
  - `go/day6part2.go` - **GET IN** Correct anticipation of the part 2 solution means the actual code-change required is a single digit added to the CONST value. `<smugface>`
  - `go/day7part1.go` - I have a feeling I made too much of a mess of implementing what is basically the `du` command here, AND I think I've done it in a non-scalable/works-for-part-2 way to boot...
  - `go/day7part2.go` - Oh. I got lucky with the part 2. I may not have built the tree-representation efficiently in part 1 but it provided exactly the data I needed to complete part2 as a linear (`O(n)`) scan of part 1 results. Hooray.
  - `go/day8part1.go` - The input parsing is fairly neat. The visibility algorithm is fairly mechanistic though, nothing too clever. Wasted about an hour not understanding how to use `range()` because I should have been using simple `for i:=;i<len(gridsize);i++` loops instead. Lots of debug `fmt.Printf` required to understand where the comparison logic was going wrong.


## NodeJS
  - `nodejs/day1part2.js` - Baa. Same algorithm as Go but twice as long to develop because 
    1. working around Node async behaviour on file read
    2. WTF is the default sort ignoring integer values and sorting by string value for? 
    3. idiot can't work a slice