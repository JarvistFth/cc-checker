package callgraphs

import (
	"cc-checker/ssautils"
	"cc-checker/utils"
	"go/types"
	"golang.org/x/tools/go/pointer"
	"golang.org/x/tools/go/ssa"
	"strings"
	"testing"
)

func buildPkg() (*ssa.Program, []*ssa.Package,*pointer.Result,*ssa.Function) {
	prog, allpkgs,err := ssautils.BuildSSA("../../ccs/timerandomcc/")
	mainpkgs,err := ssautils.MainPackages(allpkgs)
	if err != nil{
		panic(err.Error())
	}
	if len(mainpkgs)== 0{
		log.Infof("%s","empty mainpkgs")
	}else{
		log.Info(mainpkgs[0].String())
	}
	result := BuildCallGraph(mainpkgs)
	_,invokef := utils.FindInvokeMethod(prog,mainpkgs[0])
	return prog,allpkgs,result,invokef
}

func TestCHACallGraph(t *testing.T) {
	prog,_,_,_ := buildPkg()
	if reflect := prog.ImportedPackage("reflect"); reflect != nil {
		//rV := reflect.Pkg.Scope().Lookup("Value")
		//reflectValueObj := rV
		//reflectValueCall := prog.LookupMethod(rV.Type(), nil, "Call")
		//reflectType := reflect.Pkg.Scope().Lookup("Type").Type().(*types.Named)
		reflectRtypeObj := reflect.Pkg.Scope().Lookup("rtype")
		reflectRtypePtr := types.NewPointer(reflectRtypeObj.Type())
		mset := prog.MethodSets.MethodSet(reflectRtypePtr)
		log.Info(mset.String())
		//for i:=0 ;i<invokef.Signature.Params().Len(); i++ {
		//	param := invokef.Signature.Params().At(0)
		//	if strings.Contains(param.Type().String(), "ChaincodeStubInterface") {
		//		t := param.Type()
		//		fn := utils.FindMethodByType(prog,allpkgs[0],t,"PutState")
		//		//fn := prog.LookupMethod(reflectRtypePtr, call.Method.Pkg(), call.Method.Name())
		//		if fn != nil{
		//			log.Infof("%s",fn.String())
		//		}
		//	}
		//}
	}else{
		log.Warning("reflect is nil")
	}
}

func TestAllPkgMethod(t *testing.T) {
	prog,allpkgs,_,invokef := buildPkg()

	for i:=0 ;i<invokef.Signature.Params().Len(); i++ {
		param := invokef.Signature.Params().At(0)
		if strings.Contains(param.Type().String(), "ChaincodeStubInterface") {
			t := param.Type()
			fn := utils.FindMethodByType(prog,allpkgs[0],t,"PutState")
			if fn != nil{
				log.Infof("%s",fn.String())
			}
		}
	}
}

func TestBuildCallGraph	(t *testing.T) {
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

		//for i:=0 ;i<invokef.Signature.Params().Len(); i++{
		//	param := invokef.Signature.Params().At(0)
		//	if strings.Contains(param.Type().String(),"ChaincodeStubInterface"){
		//		t := param.Type()
		//		mset := prog.MethodSets.MethodSet(t)
		//		fn := prog.LookupMethod(t,param.Pkg(),"PutState")
		//		if fn != nil{
		//			log.Infof("putstate fn: %s", fn.String())
		//		}
		//		log.Infof("mset: %s", mset.String())
		//		//if _,ok := t.(*types.Interface); ok{
		//		//
		//		//}else{
		//		//	log.Warning("not interface type")
		//		//}
		//		//putstatefn = utils.FindMethodByType(prog,mainpkgs[0],ti,"PutState")
		//	}
		//}

		//invokef.WriteTo(os.Stdout)
		//log.Infof("fn:%s",invokef.String())
		//pnode := result.CallGraph.Nodes[putstatefn]

		//for callee,_:= range callees{
		//	log.Infof("putstate callee: %s, %v", callee.Func.String(), prog.Fset.Position(callee.Func.Pos()))
		//}

		//for _,edg := range nd.In{
		//	log.Infof("%s",prog.Fset.Position(edg.Pos()))
		//	log.Infof(edg.Site.String())
		//}
		//callgraph.GraphVisitEdges(result.CallGraph, func(edge *callgraph.Edge) error {
		////
		//	caller := edge.Caller
		//	calli := edge.Site
		//	callee := edge.Callee
		////
		////	//posCaller := prog.Fset.Position(caller.Func.Pos())
		////	//posCallee := prog.Fset.Position(callee.Func.Pos())
		////	//posEdge   := prog.Fset.Position(edge.Pos())
		////	//fileCaller := fmt.Sprintf("%s:%d", posCaller.Filename, posCaller.Line)
		////	//filenameCaller := filepath.Base(posCaller.Filename)
		////
		////	// omit synthetic calls
		//	if isSynthetic(edge) {
		//		return nil
		//	}
		//	if inStd(caller) || inStd(callee) {
		//		return nil
		//	}
		//	if calli != nil{
		//		if calli.Common().IsInvoke(){
		//			interfaceTypeName := calli.Common().Value.Type().String()
		//			methodName := calli.Common().Method.Name()
		//			if strings.Contains(interfaceTypeName, "ChaincodeStubInterface") && strings.Contains(methodName, "PutState"){
		//				log.Infof("dynamic call putState: %s", calli.Common().Method.FullName())
		//				putstatefn = callee.Func
		//
		//			}
		//		}
		//	}
		//	return nil
		//})
		var putstatefn *ssa.Function
		for _, out := range nd.Out{

			//log.Infof("%s",prog.Fset.Position(out.Pos()))
			if out.Site != nil{
				if out.Site.Common().IsInvoke(){
					interfaceTypeName := out.Site.Common().Value.Type().String()
					methodName := out.Site.Common().Method.Name()
					if strings.Contains(interfaceTypeName, "ChaincodeStubInterface") && strings.Contains(methodName, "PutState"){
						call := out.Site.Common()
						reflect := prog.ImportedPackage("reflect")
						if reflect != nil{
							reflectType := reflect.Pkg.Scope().Lookup("Type").Type().(*types.Named)
							if call.Value.Type() == reflectType{
								log.Infof("putstate, value:%s", call.Value.Type().String())
							}
						}
						putstatefn = out.Callee.Func
						log.Infof("dynamic call putState: pkg:%s, name:%s, type:%s", putstatefn.Pkg.Pkg.Name(),putstatefn.Name(), putstatefn.Type().String())
						//prog.LookupMethod(out.Site.Value().Type(),out.Callee.Func.Pkg.Pkg,"PutState")
					}
				}else{
					log.Infof("static call: %s",out.Site.String())
				}
			}else{

			}
		}

		//调用putState的函数
		callerNodes := result.CallGraph.Nodes[putstatefn].In
		//callees := result.CallGraph.Nodes[putstatefn].Out
		var callers []*ssa.Function
		for _,callerNode := range callerNodes {
			caller := callerNode.Caller.Func
			callers = append(callers,caller)
			log.Infof("putState caller fn: %s", caller.String())
		}




	}else{
		log.Infof("invoke func is nil\n")
	}

}

func TestWhat(t *testing.T) {

}