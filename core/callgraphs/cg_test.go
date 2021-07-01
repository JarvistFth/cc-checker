package callgraphs

import (
	"cc-checker/ssautils"
	"testing"
)

func TestBuildCallGraph	(t *testing.T) {
	prog, mainpkgs,err := ssautils.BuildSSA("../../ccs/timerandomcc/")
	if err != nil{
		t.Fatalf(err.Error())
	}
	if len(mainpkgs)== 0{
		log.Infof("%s","empty mainpkgs")
	}else{
		log.Info(mainpkgs[0].String())
	}

	result := BuildCallGraph(mainpkgs)

	fn := mainpkgs[0].Func("set")
	if fn != nil{
		nd := result.CallGraph.Nodes[fn]

		log.Infof("fn:%s",fn.String())

		for _,edg := range nd.In{
			log.Infof("%s",prog.Fset.Position(edg.Pos()))
			log.Infof(edg.Site.Parent().String())
		}

		for _, out := range nd.Out{
			log.Infof("%s",prog.Fset.Position(out.Pos()))
			log.Infof("out:%s",out.Site.Parent().String())
		}
	}

}
