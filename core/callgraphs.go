package core

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