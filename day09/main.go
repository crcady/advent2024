package main

import (
	"bufio"
	"log"
	"os"
	"strconv"
)

func mustConvert(r rune) uint {
	num, err := strconv.Atoi(string(r))
	if err != nil {
		panic(err)
	}
	return uint(num)
}

type block struct {
	id     uint
	start  uint
	length uint
}

func checksum(fs []block) uint {
	var cs uint = 0
	for _, b := range fs {
		blockSum := b.id * (((b.start + b.length - 1) * (b.start + b.length)) - ((b.start - 1) * (b.start))) / 2
		cs += blockSum
	}
	return cs
}

func compact(f []block, fs []block) []block {
	files := make([]block, len(f))
	copy(files, f)

	freespaces := make([]block, len(fs))
	copy(freespaces, fs)

	compacted := make([]block, 0, len(files))

	//Start off with the first file
	compacted = append(compacted, files[0])
	files = files[1:]

	for len(files) > 0 {
		firstFreeSpace := freespaces[0]
		lastFile := files[len(files)-1]

		if firstFreeSpace.length == lastFile.length { // Easy Case
			// Put the last file into the free space
			compacted = append(compacted, block{lastFile.id, firstFreeSpace.start, lastFile.length})
			files = files[:len(files)-1]
			freespaces = freespaces[1:]

			// Write the next file up after it
			if len(files) > 0 {
				compacted = append(compacted, files[0])
				files = files[1:]
			}

		} else if firstFreeSpace.length < lastFile.length { // Can only fit part of the file
			if len(files) == 1 { //Then we can just shift the whole thing
				compacted = append(compacted, block{lastFile.id, firstFreeSpace.start, lastFile.length})
				files = files[1:]
			} else {
				// Put as much of the last file into the free space as we can
				compacted = append(compacted, block{lastFile.id, firstFreeSpace.start, firstFreeSpace.length})
				files[len(files)-1] = block{lastFile.id, lastFile.start, lastFile.length - firstFreeSpace.length}
				freespaces = freespaces[1:]

				// Write the next file up after it
				compacted = append(compacted, files[0])
				files = files[1:]
			}
		} else { // Can fit the entire file, and have remaining free space
			// Move the file and adjust the free space
			compacted = append(compacted, block{lastFile.id, firstFreeSpace.start, lastFile.length})
			files = files[:len(files)-1]
			freespaces[0] = block{0, firstFreeSpace.start + lastFile.length, firstFreeSpace.length - lastFile.length}
		}
	}

	return compacted
}

type typedBlock struct {
	isFile bool
	block  block
}

func compact2(f []block, fs []block) []block {
	files := make([]block, len(f))
	copy(files, f)

	freespaces := make([]block, len(fs))
	copy(freespaces, fs)

	compacted := make([]block, 0, len(files))

	entries := make([]typedBlock, 0)

	for i := range files {
		entries = append(entries, typedBlock{true, files[i]})
		if len(freespaces) > i {
			entries = append(entries, typedBlock{false, freespaces[i]})
		}
	}

	for len(compacted) < len(files) {
		next := entries[0]

		if next.isFile {
			compacted = append(compacted, next.block)
			entries = entries[1:]
			continue
		}

		l := next.block.length
		movedOne := false

		for i := len(entries) - 1; i > 0; i-- { //Not an off-by-one, dont need to look at zero-th entry
			entry := entries[i]
			if !entry.isFile {
				continue
			}

			if entry.block.length > l {
				continue
			}

			compacted = append(compacted, block{entry.block.id, next.block.start, entry.block.length})
			entries[0].block = block{next.block.id, next.block.start + entry.block.length, next.block.length - entry.block.length}
			entries = append(entries[:i], entries[i+1:]...)
			movedOne = true

			break
		}

		if !movedOne {
			entries = entries[1:]
		}
	}
	return compacted
}

func main() {
	fname := "example.txt"

	if len(os.Args) > 1 {
		fname = os.Args[1]
	}

	log.Printf("Starting processing of %s", fname)

	file, err := os.Open(fname)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	scanner.Scan()
	input := scanner.Text()
	files := make([]block, 0, len(input)/2)
	freespaces := make([]block, 0, len(input)/2)
	var pos uint = 0
	var nextID uint = 0

	for i, r := range input {
		n := mustConvert(r)
		if i%2 != 0 { //free space
			freespaces = append(freespaces, block{0, pos, n})
			pos += n
		} else {
			files = append(files, block{nextID, pos, n})
			pos += n
			nextID++
		}
	}

	compacted := compact(files, freespaces)

	log.Println(checksum(compacted))

	compacted2 := compact2(files, freespaces)

	log.Println(checksum(compacted2))
}
