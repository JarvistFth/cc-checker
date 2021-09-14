package core

import (
	"cc-checker/config"
	"cc-checker/ssautils"
	"cc-checker/utils"
	"golang.org/x/tools/go/callgraph"
	"golang.org/x/tools/go/ssa"
	"golang.org/x/tools/go/ssa/ssautil"
)

var prog *ssa.Program

var invokef *ssa.Function

var cfg *callgraph.Graph


func Main(){
	cfg,err := config.ReadConfig("../config/config.yaml")
	if err != nil{
		log.Fatalf(err.Error())
	}
	log.Info(cfg.String())

	var mains []*ssa.Package

	prog,mains,err = ssautils.BuildSSA("../ccs/timerandomcc/")
	if err != nil{
		log.Fatalf(err.Error())
	}

	_,invokef = utils.FindInvokeMethod(prog,mains[0])
	BuildCallGraph(mains)
	StartAnalysis(invokef)
}



func StartAnalysis(fn *ssa.Function) {

	if fn == nil{

		return
	}

	outputEdges := cfg.Nodes[fn].Out

	var outputNodes []*callgraph.Node

	for _,outputEdge := range outputEdges{
		//check source
		outputNodes = append(outputNodes,outputEdge.Caller)

	}

	ssautil.AllFunctions()



	for len(outputNodes) > 0{
		front := outputNodes[0]
		outputNodes = outputNodes[1:]

		//
	}







}

