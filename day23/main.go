package main

import (
	"bufio"
	"log"
	"os"
	"sort"
	"strings"

	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/simple"
	"gonum.org/v1/gonum/graph/topo"
)

type network struct {
	nodes map[string]map[string]bool
}

func newNetwork() network {
	theMap := make(map[string]map[string]bool)
	return network{theMap}
}

func (n network) add(s string) {
	vals := strings.Split(s, "-")
	node1, node2 := vals[0], vals[1]

	if _, ok := n.nodes[node1]; ok {
		n.nodes[node1][node2] = true
	} else {
		n.nodes[node1] = map[string]bool{node2: true}
	}

	if _, ok := n.nodes[node2]; ok {
		n.nodes[node2][node1] = true
	} else {
		n.nodes[node2] = map[string]bool{node1: true}
	}
}

func (n network) solve1() int {
	trigraphs := map[[3]string]bool{}

	for n1 := range n.nodes {
		if n1[0] != 't' {
			continue
		}

		for n2 := range n.nodes[n1] { // All of the neighbors of n
			for n3 := range n.nodes[n1] {
				if n2 == n3 {
					continue
				}

				if n.nodes[n2][n3] {
					triplet := [3]string{n1, n2, n3}
					sort.Strings(triplet[:])
					trigraphs[triplet] = true
				}
			}
		}

	}

	return len(trigraphs)
}

func (n network) solve2() string {
	nodeIDs := map[string]int{}
	backIDs := map[int]string{}

	for n := range n.nodes {
		nodeIDs[n] = len(nodeIDs)
		backIDs[len(nodeIDs)-1] = n
	}

	g := simple.NewUndirectedGraph()

	for n1, neighbors := range n.nodes {
		for n2 := range neighbors {
			edge := g.NewEdge(simple.Node(nodeIDs[n1]), simple.Node(nodeIDs[n2]))
			g.SetEdge(edge)
		}
	}

	cliques := topo.BronKerbosch(g) // Returns the set of maximal cliques

	party := []graph.Node{}

	for _, c := range cliques {
		if len(c) > len(party) {
			party = c
		}
	}

	names := make([]string, len(party))
	for i, n := range party {
		names[i] = backIDs[int(n.ID())]
	}

	sort.Strings(names)

	return strings.Join(names, ",")
}

func main() {
	fname := "example.txt"

	if len(os.Args) > 1 {
		fname = os.Args[1]
	}

	file, err := os.Open(fname)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	n := newNetwork()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		n.add(line)
	}

	ans1 := n.solve1()
	log.Println("Found", ans1, "sets of three connected computers")

	ans2 := n.solve2()
	log.Println("The party password is", ans2)
}
