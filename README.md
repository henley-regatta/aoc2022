# aoc2022 

[2022 Advent Of Code](https://adventofcode/2022) - My solutions

I really, really am not going to take this seriously. This is just a bit of fun. 
This year I'm just going to pick at it in _other_ languages. Not just Python.

## Go
  - `go/day1part1.go` - seemed obvious to dynamically create an int map incrementing "key" (elf index) on a blank line (_i.e._ one that can't be `str.AtoI` converted)
  - `go/day1part2.go` - Let's not dignify this with "efficient" as a characteristic but sorting the previous output by "value" (per-elf calorie count) instead of key seemed obvious, and then taking a slice of the top 3 entries and summing them was the natural conclusion


## NodeJS
  - `nodejs/day1part2.js` - Baa. Same algorithm as Go but twice as long to develop because 
    1. working around Node async behaviour on file read
    2. WTF is the default sort ignoring integer values and sorting by string value for? 
    3. idiot can't work a slice