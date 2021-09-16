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
	sinks map[ssa.CallInstruction]bool
	sinkArgs map[ssa.Value]bool

	returnInstr map[*ssa.Function]*ssa.Return
}

func NewVisitor() *visitor {
	return &visitor{
		seen: make(map[*callgraph.Node]bool),
		lattice: make(map[ssa.Value]map[string]bool),

		
		returnInstr: make(map[*ssa.Function]*ssa.Return),
	}
}

func (v *visitor) Visit(node *callgraph.Node) {
	if !v.seen[node]{
		log.Infof("traverse visit: %s", node.String())
		v.seen[node] = true

		//check source

		v.loopFunction(node.Func)

		v.handleReturnValue()

		v.handleSinkDetection()


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

func (v *visitor) loopFunction(fn *ssa.Function) {
	if fn == nil{
		log.Errorf("fn == nil!")
		return
	}
	for _,block := range fn.Blocks{
		for _,instr := range block.Instrs{
			if ret,ok := instr.(*ssa.Return); ok {
				v.returnInstr[fn] = ret
				continue
			}


			if call,ok := instr.(ssa.CallInstruction); ok	{
				v.checkSource(call)

				v.checkSink(call)
			}
		}
	}

}


// check source :
// if isSource, put call.value into taintSet
func (v *visitor) checkSource(call ssa.CallInstruction) {
	if tag,yes := config.IsSource(call); yes{
		v.taintVal(call.Common().Value,tag)
	}
}

func (v *visitor) checkSink(call ssa.CallInstruction) {
	if ok := config.IsSink(call); ok {
		v.sinks[call] = true
		for _,arg := range call.Value().Call.Args{
			v.sinkArgs[arg] = true
		}
	}
}



func(v *visitor) taintReferrers(node *callgraph.Node) {

	for val,m := range v.lattice{
		if val.Parent() != nil && val.Parent() == node.Func{
			for _,ref := range *val.Referrers(){
				if refval,ok := ref.(ssa.Value);ok{
					for tag,_ := range m{
						v.taintVal(refval,tag)
					}
				}
			}
		}

	}



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

func (v *visitor) handleSinkDetection() bool {
	for arg,_ := range v.sinkArgs{
		if m,ok := v.lattice[arg];ok{

			//report detection
			var sinkTag string
			for tag, _ := range m{
				sinkTag += tag + " "
			}
			log.Warning("sink here %sinkTag with tag:%sinkTag", prog.Fset.Position(arg.Pos()), sinkTag)
			return true
		}
	}
	return false
}

func (v *visitor) handleReturnValue(node *callgraph.Node) {
	 ret,ok := v.returnInstr[node.Func]
	 if ok{
		 returnValues := ret.Results
		 var tags string

		 for _,result := range returnValues{
			 if m,ok := v.lattice[result]; ok{
				 for tag,_ := range m{
					 tags += tag + " "
				 }
			 }
		 }

		 inEdges := node.In

		 for _,inEdge := range inEdges{
			 inEdge.Caller
		 }

	 }
}


