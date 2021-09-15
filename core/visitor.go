package core

import (
	"cc-checker/config"
	"cc-checker/utils"
	"golang.org/x/tools/go/callgraph"
	"golang.org/x/tools/go/ssa"
)

type visitor struct {

	lattice map[ssa.Value]map[string]bool
	seen map[*callgraph.Node]bool
}

func NewVisitor() *visitor {
	return &visitor{
		seen: make(map[*callgraph.Node]bool),
		lattice: make(map[ssa.Value]map[string]bool),
	}
}

func (v *visitor) Visit(node *callgraph.Node) {
	if !v.seen[node]{
		log.Infof("traverse visit: %s", node.String())
		v.seen[node] = true

		//check source

		v.taintReferrers(node)

		//check sink


		for _,outputEdge := range node.Out{
			if utils.IsSynthetic(outputEdge) || utils.InStd(outputEdge.Callee) || utils.InFabric(outputEdge.Callee){
				v.seen[outputEdge.Callee] = true
				continue
			}
			log.Infof("out: %s", outputEdge.Callee.String())

			//根据当前的taint的情况，设定函数入参的lattice情况
			v.Visit(outputEdge.Callee)
		}

	}
}

func(v *visitor) taintReferrers(node *callgraph.Node) {

	fn := node.Func

	// if is source:
	for _,block := range fn.Blocks{
		for _,instr := range block.Instrs{
			if call,ok := instr.(ssa.CallInstruction);ok{
				tag,yes := config.IsSource(call)
				if yes{
					callValue := call.Value()
					var buf [10]*ssa.Value // avoid alloc in common case
					if callValue.Referrers() != nil{
						for _,referer := range *callValue.Referrers(){
							if refCall,ok := referer.(ssa.CallInstruction); ok{
								ops := refCall.Value().Operands(buf[:])
								for _,op := range ops{
									if op != nil{
										log.Infof("op: %s", (*op).String())
									}
								}

								continue
							}
							if val,ok := referer.(ssa.Value);ok{
								v.taintVal(val,tag)
							}
						}
					}
				}
			}
		}
	}

	// taint referers



}

func (v *visitor) taintVal(val ssa.Value, tag string) {
	if v.lattice[val] == nil{
		v.lattice[val] = make(map[string]bool)
		v.lattice[val][tag] = true
		return
	}

	if _,found := v.lattice[val][tag]; found{
		return
	}else{
		v.lattice[val][tag] = true
	}

}

func (v *visitor) taintCallArgs() {

}
