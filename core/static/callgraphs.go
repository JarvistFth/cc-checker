package static

import (
	"golang.org/x/tools/go/callgraph"
	"golang.org/x/tools/go/pointer"
	"golang.org/x/tools/go/ssa"
)

var CallGraphs *callgraph.Graph

func BuildCallGraph(mainpkg []*ssa.Package, querys ...*ssa.Function) *pointer.Result {
	cfg := &pointer.Config{
		Mains:           mainpkg,
		BuildCallGraph: true,

	}

	//for fn,_ := range ssautil.AllFunctions(prog){
	//	for _,block := range fn.Blocks{
	//		for _,instr := range block.Instrs{
	//			if val,ok := instr.(ssa.Value);ok{
	//				if utils.CanPointToVal(val){
	//					cfg.AddQuery(val)
	//				}
	//				if utils.CanPointToInDirect(val){
	//					cfg.AddIndirectQuery(val)
	//				}
	//			}
	//		}
	//	}
	//}

	//for _, fn := range querys{
	//	for _,block := range fn.Blocks{
	//		for _, instr := range block.Instrs{
	//			if val,ok := instr.(ssa.Value);ok{
	//				if utils.CanPointToVal(val){
	//					cfg.AddQuery(val)
	//				}
	//				if utils.CanPointToInDirect(val){
	//					cfg.AddIndirectQuery(val)
	//				}
	//			}
	//		}
	//	}
	//}

	result,err := pointer.Analyze(cfg)
	if err != nil{
		log.Errorf(err.Error())
		return nil
	}
	CallGraphs = result.CallGraph
	return result
}
