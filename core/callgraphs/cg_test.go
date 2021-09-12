package callgraphs

import (
	"cc-checker/ssautils"
	"cc-checker/utils"
	"golang.org/x/tools/go/callgraph"
	"strings"
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
			if isSynthetic(edge){
				return nil
			}
			if edge.Site == nil {
				//log.Infof("%s",edge.String())
			}else{
				if strings.Contains(edge.Callee.String(),"PutState") || strings.Contains(edge.Caller.String(),"PutState") {
					log.Infof("putState Caller: %s", edge.Caller.String())
					log.Infof("putState Callee: %s", edge.Callee.String())
				}

				if edge.Site != nil && edge.Site.Common().IsInvoke(){
					//if strings.Contains(edge.Site.Common().Method.String(), "PutState"){
					//	log.Infof("dynamic call: %s, caller:%s, callee:%s", edge.Site.Common().Method.FullName(),edge.Caller.String(),edge.Callee.String())
					//}
					log.Infof("dynamic call:%s", edge.Site.Common().Method.String())
				}
				//if edge.Site.Common().StaticCallee() == nil{
				//	log.Infof("dynamic call:%s", edge.Site.String())
				//}
				////log.Debugf(edge.Site.String())
				//if strings.Contains(edge.Callee.String(),"fabric"){
				//	log.Infof("fabric pkg callee:%s", edge.Callee.String())
				//}

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