package core

import (
	"golang.org/x/tools/go/callgraph"
	"golang.org/x/tools/go/pointer"
	"golang.org/x/tools/go/ssa"
)

var CallGraphs *callgraph.Graph

func BuildCallGraph(mainpkg []*ssa.Package) *pointer.Result {
	cfg := &pointer.Config{
		Mains:           mainpkg,
		BuildCallGraph: true,
	}
	result,err := pointer.Analyze(cfg)

	if err != nil{
		log.Errorf(err.Error())
		return nil
	}
	CallGraphs = result.CallGraph
	return result
}

func GetCallerBlock(cg *callgraph.Graph, fn *ssa.Function) []*ssa.BasicBlock {
	in := cg.Nodes[fn].In
	//out := cg.Nodes[fn].Out



	var ret []*ssa.BasicBlock

	for _, i := range in{
		ret = append(ret,i.Site.Block())
		i.Site.Value().Referrers()
	}

	return ret

}