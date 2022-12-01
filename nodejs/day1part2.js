#!/usr/bin/env node

/*
aoc2022 - Day 1, Part 2
-----------------------
Hungry Hungry Elves

Each  Elf enumerates the calorific content of their snacks, 1 snack per line.
Elves are separated in the input by a blank line.

Find the top THREE chonk-carrying Elves. How many calories do they carry all together?
*/
const fs = require('fs');
const readline = require('readline');

// for testing:
//const dataFile = "data/day1test.txt"
// for Stars:
const dataFile = "data/day1input.txt"

var calorieMap = {};
var elfNdx = 0;

async function main() {
    const rl = readline.createInterface({
        input: fs.createReadStream(dataFile),
        crlfDelay: Infinity
    });

    
    
    for await (const line of rl) {
        const text = line.split(' ')[0];
        v = Number(text)
        if(typeof v === "number" && v > 0) {
            if(typeof calorieMap[elfNdx] === "number") { 
                calorieMap[elfNdx] = (calorieMap[elfNdx] + v);
            } else {
                calorieMap[elfNdx] = (v)
            }
        } else {
            elfNdx = elfNdx + 1
        }
    }

    //post-processing, sort calorieMap descending, slice top 3
    //(convert to Array to facilitate; we don't mind losing elfNum)
    const chonkValues = Object.keys(calorieMap).map(function (e) {
        return calorieMap[e]
    })
    
    chonkValues.sort(function (a,b){ return a - b }); //WTF is default sort "convert to UTF-16 and compare strings" for?
    chonkValues.reverse();
    //Get slice of first 3 values and sum
    var sumCals = 0
    chonkValues.slice(0,3).forEach(function (v) {
        sumCals += v;
    })
    
    console.log(sumCals)
}


if (require.main === module) {
    main();
}
