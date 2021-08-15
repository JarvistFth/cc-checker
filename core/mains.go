package core

import (
	"cc-checker/config"
	"cc-checker/ssautils"
	"cc-checker/utils"
	"golang.org/x/tools/go/ssa"
)

var prog *ssa.Program

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

	_,invokef := utils.FindInvokeMethod(prog,mains[0])

	BuildCallGraph(mains)
	StartAnalysis(invokef)
}



func StartAnalysis(fn *ssa.Function) {
	entryCtx := NewCallContext(fn,TaintParams{})
	entryCtx.Initialize()

}

