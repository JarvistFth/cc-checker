package core

import (
	"cc-checker/utils"
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
		log.Infof("traverse visit: %s", node.String())
		v.seen[node] = true

		//check source

		//check sink


		for _,outputEdge := range node.Out{
			if utils.IsSynthetic(outputEdge) || utils.InStd(outputEdge.Callee) || utils.InFabric(outputEdge.Callee){
				v.seen[outputEdge.Callee] = true
				continue
			}
			log.Infof("out: %s", outputEdge.Callee.String())
			v.Visit(outputEdge.Callee)
		}

	}
}
