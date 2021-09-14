package core

import (
	"cc-checker/ssautils"
	"cc-checker/utils"
	"testing"
)

func TestVisitor_Visit(t *testing.T) {
	prog, allpkgs,err := ssautils.BuildSSA("../../ccs/timerandomcc/")
	mainpkgs,err := ssautils.MainPackages(allpkgs)
	if err != nil{
		t.Fatalf(err.Error())
	}
	if len(mainpkgs)== 0{
		log.Infof("%s","empty mainpkgs")
	}else{
		log.Info(mainpkgs[0].String())
	}

	result := BuildCallGraph(mainpkgs)
	_,invokef := utils.FindInvokeMethod(prog,mainpkgs[0])
	//var putstatefn *ssa.Function
	//fn := mainpkgs[0].Func("set")
	if invokef != nil{
		nd := result.CallGraph.Nodes[invokef]
		v := NewVisitor()
		for _, out := range nd.Out{
			v.Visit(out.Caller)
		}
	}else{
		log.Infof("invoke func is nil\n")
	}
}
