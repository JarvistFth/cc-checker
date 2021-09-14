package core

import (
	"golang.org/x/tools/go/callgraph"
)

type visitor struct {


	seen map[*callgraph.Node]bool
}

func NewVisitor() *visitor {
	return &visitor{seen: make(map[*callgraph.Node]bool)}
}

func (v *visitor) Visit(node *callgraph.Node) {
	if !v.seen[node]{
		v.seen[node] = true

		log.Infof("traverse visit: %s", node.String())
		//check source

		//check sink


		for _,outputEdge := range node.Out{
			v.Visit(outputEdge.Caller)
		}

	}
}
