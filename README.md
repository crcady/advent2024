# 2024 Advent of Code
This repo contains my solutions to the 2024 [Advent of Code](https://adventofcode.com) challenge. All problems are solve using Go.

## Organization
Each day has it's own directory, and each directory has it's own `go.mod` file. For the handful of days that use an external Go module, the `go.sum` file is committed as well.

In accordance with the Advent of Code FAQ page, I have not included the puzzle inputs here.

## How to build and run
All of these work just fine with `go run .`, no need for a fancy Makefile here.

Most of the solutions expect to have the smaller/test input provided in a file named `example.txt`. If the binary that is built is called with no inputs, the example will be solved. Most, but **not all**, of the days will work just fine if you instead run `./dayXX example.txt`, but you shouldn't do this because some will fail. A couple of days have other parameters (like a size) that can't be inferred from the input itself, and these are automatically configured to use the example or real sizes based on whether the solver is passed a filename or not; one or two of them are so simple that the example input is just hard-coded in. Just invoke it with the binary name or `go run .`

To solve an actual instance, you **must** provide the filename to the executable, for example `./dayXX input.txt` or `go run . input.txt`.

### Day 24
Day 24 is the exception to the above. The solution is only half-complete, in the sense that you have to run it multiple times and follow the instructions. I don't want to spoil the challenge, but you'll see what I mean if you get there ;)

## Dependencies
I only used two dependencies for this: a set of Z3 bindings and some of the [Gonum](https://www.gonum.org) packages. I used Z3 to automate doing integer or bitvector arithmetic, and Gonum for some graph algorithms. All of the arithmetic could be worked out by hand, and I mostly just used Dijkstra's algorithm from the Gonum package, so it probably would have been faster for me to implement the algorithms myself than learn to use someone else's package, but I wanted to learn how to interact with the broader Go ecosystem.

### Gonum
Gonum should work just fine without any additional configuration.

### Z3
The Z3 bindings I use don't ship the Z3 source or binary. You need to install Z3 and make sure the Go compiler/linker can find it by setting an environment variable. Look at the day13 readme and the readme in the depencency's Github readme for additional information.