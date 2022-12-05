package dj_test

import (
	"testing"

	"github.com/go-dj/dj"
	"github.com/stretchr/testify/require"
)

func TestBFS(t *testing.T) {
	// Create a graph.
	graph := map[string][]string{
		"1": {"2", "3"},
		"2": {"4", "5"},
		"3": {"6", "7"},
		"4": {"8", "9"},
		"5": {"10", "11"},
		"6": {"12", "13"},
		"7": {"14", "15"},
	}

	// Create a function that returns the next nodes.
	next := func(v string) []string {
		return graph[v]
	}

	var visited []string

	// Create a function that records the visited nodes.
	fn := func(v string) {
		visited = append(visited, v)
	}

	// Perform a breadth-first search.
	dj.BFS("1", next, fn)

	// Check the visited nodes.
	require.Equal(t, []string{
		"1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11", "12", "13", "14", "15",
	}, visited)
}

func TestDFS(t *testing.T) {
	// Create a graph.
	graph := map[string][]string{
		"1": {"2", "3"},
		"2": {"4", "5"},
		"3": {"6", "7"},
		"4": {"8", "9"},
		"5": {"10", "11"},
		"6": {"12", "13"},
		"7": {"14", "15"},
	}

	// Create a function that returns the next nodes.
	next := func(v string) []string {
		return graph[v]
	}

	var visited []string

	// Create a function that records the visited nodes.
	fn := func(v string) {
		visited = append(visited, v)
	}

	// Perform a depth-first search.
	dj.DFS("1", next, fn)

	// Check the visited nodes.
	require.Equal(t, []string{
		"1", "3", "7", "15", "14", "6", "13", "12", "2", "5", "11", "10", "4", "9", "8",
	}, visited)
}
