package callgraphs

import (
	"cc-checker/ssautils"
	"cc-checker/utils"
	"golang.org/x/tools/go/callgraph"
	"path/filepath"
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
		nd := result.CallGraph.Nodes[invokef]

		log.Infof("fn:%s",invokef.String())

		//for _,edg := range nd.In{
		//	log.Infof("%s",prog.Fset.Position(edg.Pos()))
		//	log.Infof(edg.Site.String())
		//}

		callgraph.GraphVisitEdges(result.CallGraph, func(edge *callgraph.Edge) error {

			caller := edge.Caller
			callee := edge.Callee

			posCaller := prog.Fset.Position(caller.Func.Pos())
			//posCallee := prog.Fset.Position(callee.Func.Pos())
			//posEdge   := prog.Fset.Position(edge.Pos())
			//fileCaller := fmt.Sprintf("%s:%d", posCaller.Filename, posCaller.Line)
			filenameCaller := filepath.Base(posCaller.Filename)

			// omit synthetic calls
			if isSynthetic(edge) {
				return nil
			}


			// omit std
			if inStd(caller) || inStd(callee) {
				return nil
			}



			//logf("call node: %s -> %s\n %v", caller, callee, string(data))
			log.Infof("call node: %s -> %s (%s -> %s) %v\n", caller.Func.Pkg, callee.Func.Pkg, caller, callee, filenameCaller)


			return nil
		})

		for _, out := range nd.Out{
			log.Infof("%s",prog.Fset.Position(out.Pos()))
			if out.Site != nil{
				log.Infof("out:%s",out.Site.String())
			}else{
				if out.Site.Common().IsInvoke(){
					log.Infof("dynamic call out: %s", out.Site.Common().Method.FullName())
				}
			}
		}


	}else{
		log.Infof("invoke func is nil\n")
	}

}

func TestWhat(t *testing.T) {

}