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
  - `go/day8part2.go` - Took a while to recognise that the direction-of-scan matters when you're looking for _FIRST_ blocking tree not just _ANY_ blocking tree. And the off-by-one-ism of my indexing and scoring means I'm not really satisfied with the code - it works but I'm not 100% sure I've internalised where I've iterated to/from on the checks. Oh well, an answer is an answer....
  - `go/day9part1.go` - Distracted by other things this is a "bitty" solution arrived at by stages. Another mechanistic rather than inspired algorithm but it does reproduce the steps in the explanation.
  - `go/day9part2.go` - OK, so my mechanistic approach to part 1 paid off, especially the functional breakdown, because extending a "knot" to a "rope" meant only changing one data-structure to an array of 10 elements and then adding an iterate-down-the-line call on each move. Again: I'm sure there's more efficient ways to do it, I *know* there's cleverer ways to do it, but this works in respectable time and with decent efficiency...
  - `go/day10part1.go` - Aha, first simple vm of the year. Interesting mechanism but I think I've implemented the cycle-value spying routine in a way that won't help me in Part 2. Spent more time on fixing Go syntax/grammer errors than on the core problem here.
  - `go/day10part2.go` - Got lucky again as Part 2 doesn't extend the signal strength but diverts into a whole other area. Wonderful little problem with some fun fun output. Really enjoyed this one!
  - `go/day11part1.go` - Most of the effort here is in parsing the input, the actual algorithm follows from the instructions quite nicely. And works.
  - `go/day11part2.go` - As my son puts it: Properly _lanternfished_ by this one. Got an answer dead quick but completely missed the fact that fearLevel escalates way past the bound of any reasonable `int` size very very quickly, masked by the fact that _GO silently wraps around_ and does not provide any indicator of overflow. Given code's behaviour I was able to deduce this myself rapidly but was very stuck by how to tackle it. Needed a long talk with a tame (trainee) Mathmatician to convince myself that a) remainders were sufficient provided b) you get the _right_ remainder at the right time (i.e. choosing the modulus needs a global approach not per-monkey). Actual code changes required trivial once the problem is understood.
  - `go/day12part1.go` - So I should get some points for spotting this is a [Dijkstra's Algorithm](https://isaaccomputerscience.org/concepts/dsa_search_dijkstra) problem right off the bat. **BUT** I should lose lots and lots and lots of points for:
    - Completely making a horlics of implementing Dijkstra because I _refuse_ to read the documents
    - Wasting hours on debugging because I _failed_ to spot the specification of height changes was asymmetric
    - Failing to take account of the fact that the map can contain "islands of unreachability" and that the classic end-conditions won't work.
  - `go/day12part2.go` - A seemingly-simple elaboration ("Start from ALL lowest points") that should just require iterating on start point. And yet... this found a termination condition I hadn't accounted for which caused a crash which caused a lot of tracing to find out that "no route found" was something I wasn't testing for. Sort that out, implement a simple sub-route test (cheaper to lookup existing routes to see if a start point was already on it than to calculate the whole route again, but it's a fairly rare condition in the input set so it doesn't save much time). Since there's 2000 more routes to find it's clearly not a _quick_ solution but it's Good Enough (~2 minutes execution) so here it is submitted.
  - `go/day12part1_visualisation.go` and `go/day12part2_visualisation.go` - some days cry out to be visualised so hard that even though you spent so long and should just _move on_, you have to give it a go. At least I learned how to make PNGs in GO from it. Output available [here](https://www.guided-naafi.org/aoc2022/2022/12/12/VisualisationOfAOC2022Day12.html)
  - `go/day13part1.go` - An exercise in believing in booleans when there's explitly defined a tristate of continue. Hours of debugging and analysis revealed eventually that I was using "it passes" as "continue evaluation" on return from recursion so a pass _could_ (And given the input inevitably _did_) result in later override by a fail, aka the wrong answer. Naturally this doesn't come up in the sample data.......
  - `go/day13part2.go` - Aww YIS. "Sort all the packets" was a red herring; it's sufficient only to know which are smaller than the little master, and which are bigger than the big master packet. Which makes the algorithm lovely and easy to do (once you get past your off-by-one indexing problems...)
  - `go/day14part1.go` - There's no tricks or cleverness in this, I just did the mechanical grind through the problem even visualising it (in ASCII, natch) as I went. 
  - `go/day14part2.go` - This is just a resetting of the termination condition IF you've scaled and moved the map appropriately. A feat it took me far longer than it should have to achieve successfully....
  - `go/day14part2_visualisation.go` - The problem was _screaming_ for a visualisation so here's the code to generate one and [here](https://www.guided-naafi.org/aoc2022/2022/12/14/VisualisationOfAOC2022Day14.html) is the output

## NodeJS
  - `nodejs/day1part2.js` - Baa. Same algorithm as Go but twice as long to develop because 
    1. working around Node async behaviour on file read
    2. WTF is the default sort ignoring integer values and sorting by string value for? 
    3. idiot can't work a slice