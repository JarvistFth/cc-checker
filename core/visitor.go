package core

import (
	"cc-checker/config"
	"cc-checker/utils"
	"golang.org/x/tools/go/callgraph"
	"golang.org/x/tools/go/ssa"
)

type visitor struct {
	lattice  map[ssa.Value]map[string]bool
	seen     map[*callgraph.Node]bool
	sinkArgs map[ssa.Value]bool
}

func NewVisitor() *visitor {
	return &visitor{
		seen:     make(map[*callgraph.Node]bool),
		lattice:  make(map[ssa.Value]map[string]bool),
		sinkArgs: make(map[ssa.Value]bool),
	}
}

func (v *visitor) Visit(node *callgraph.Node) {
	if !v.seen[node] {
		log.Infof("traverse visit: %s", node.String())
		v.seen[node] = true

		//check source
		v.loopFunction(node)

		for _, outputEdge := range node.Out {
			if utils.IsSynthetic(outputEdge) || utils.InStd(outputEdge.Callee) || utils.InFabric(outputEdge.Callee) {
				v.seen[outputEdge.Callee] = true
				continue
			}
			log.Infof("out: %s", outputEdge.Callee.String())

			//根据当前的taint的情况，设定函数入参的lattice情况
			v.Visit(outputEdge.Callee)
		}

	}
}

func (v *visitor) loopFunction(node *callgraph.Node) {
	fn := node.Func
	if fn == nil {
		log.Errorf("fn == nil!")
		return
	}
	for _, block := range fn.Blocks {
		for _, i := range block.Instrs {

			switch instr := i.(type) {
			case ssa.CallInstruction:

				v.checkSource(instr)
				v.checkSink(instr)

			case *ssa.Return:
				v.handleReturnValue(node, instr)
			}
		}
	}
	v.handleSinkDetection()

}

// check source :
// if isSource, put call.value into taintSet
func (v *visitor) checkSource(callInstr ssa.Instruction) {
	if tag, yes := config.IsSource(callInstr.(ssa.CallInstruction)); yes {
		log.Infof("source fn here: %s", prog.Fset.Position(callInstr.Pos()))
		v.taint(callInstr, tag)
	}
}

func (v *visitor) checkSink(callInstr ssa.Instruction) {
	if ok := config.IsSink(callInstr.(ssa.CallInstruction)); ok {
		log.Infof("sink fn here: %s", prog.Fset.Position(callInstr.Pos()))
		for _, arg := range callInstr.(ssa.CallInstruction).Value().Call.Args {
			v.sinkArgs[arg] = true
		}
	}
}

func (v *visitor) taint(i ssa.Instruction, tag string) {
	switch val := i.(type) {
	case ssa.Value:
		v.taintReferrers(i, tag)
		v.taintVal(val, tag)
	case *ssa.Store:
		v.taintReferrers(val, tag)
		v.taintVal(val.Val, tag)
		v.taintVal(val.Addr, tag)
	default:
		return
	}

}

func (v *visitor) taintReferrers(i ssa.Instruction, tag string) {

	if val, ok := i.(ssa.Value); ok {
		if val.Referrers() == nil {
			return
		}

		if v.alreadyTainted(val, tag) {
			return
		}

		for _, r := range *val.Referrers() {
			v.taint(r, tag)
		}
	} else {
		log.Warningf("instr: %s is not a value", i.String())
	}

}

func (v *visitor) taintVal(val ssa.Value, tag string) {

	if v.alreadyTainted(val, tag) {
		return
	}
	log.Debugf("taintval: %s, %s, tag:%s", val.Name(), val.String(), tag)
	if v.lattice[val] == nil {
		v.lattice[val] = make(map[string]bool)
	}
	v.lattice[val][tag] = true
}

func (v *visitor) handleSinkDetection() bool {
	for arg, _ := range v.sinkArgs {
		if m, ok := v.lattice[arg]; ok {

			//todo: report detection
			var sinkTag string
			for tag, _ := range m {
				sinkTag += tag + " "
			}
			log.Warningf("sink here %s sinkTag with tag:%s sinkTag", prog.Fset.Position(arg.Pos()), sinkTag)
			return true
		}
	}
	return false
}

func (v *visitor) handleReturnValue(node *callgraph.Node, retInstr *ssa.Return) {

	ins := node.In

	returnValues := retInstr.Results
	for _, result := range returnValues {
		if tags, found := v.lattice[result]; found {

			for tag, _ := range tags {
				for _, in := range ins {
					callsite := in.Site
					v.taint(callsite, tag)
				}
			}

		}
	}
}

func (v *visitor) alreadyTainted(val ssa.Value, tag string) bool {

	if v.lattice[val] == nil {
		return false
	}

	if _, found := v.lattice[val][tag]; !found {
		return false
	}

	return true

}
