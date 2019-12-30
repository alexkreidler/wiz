/*
Package gutils implements some basic graph utilities
*/
package gutils

import (
	"gonum.org/v1/gonum/graph"
)

func IterateChildNodes(nodes graph.Nodes, f func(n graph.Node)) {
	nodes.Next()
	for nodes.Len() > 0 {
		f(nodes.Node())
		nodes.Next()
	}
}
