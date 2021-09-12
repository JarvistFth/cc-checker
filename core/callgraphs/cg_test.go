package callgraphs

import (
	"cc-checker/ssautils"
	"cc-checker/utils"
	"golang.org/x/tools/go/callgraph"
	"strings"
	"testing"
)

func TestBuildCallGraph	(t *testing.T) {
	prog, mainpkgs,err := ssautils.BuildSSA("../../ccs/voting/")
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

	//fn := mainpkgs[0].Func("set")
	if invokef != nil{
		//nd := result.CallGraph.Nodes[invokef]

		log.Infof("fn:%s",invokef.String())

		//for _,edg := range nd.In{
		//	log.Infof("%s",prog.Fset.Position(edg.Pos()))
		//	log.Infof(edg.Site.String())
		//}

		callgraph.GraphVisitEdges(result.CallGraph, func(edge *callgraph.Edge) error {
			if edge.Site == nil {
				log.Infof("%s",edge.String())
			}else{
				log.Debugf(edge.Site.String())
				if strings.Contains(edge.Site.String(),"PutState")  {
					log.Infof("putState Callee: %s", edge.Caller.String())
				}
			}
			return nil
		})

		//for _, out := range nd.Out{
		//	log.Infof("%s",prog.Fset.Position(out.Pos()))
		//	log.Infof("out:%s",out.Site.String())
		//}
	}else{
		log.Infof("invoke func is nil\n")
	}

}

func TestWhat(t *testing.T) {

}