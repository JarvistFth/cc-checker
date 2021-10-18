package core

import (
	"cc-checker/config"
	"cc-checker/logger"
	"cc-checker/utils"
	"fmt"
	"go/types"
	"golang.org/x/tools/go/callgraph"
	"golang.org/x/tools/go/ssa"
)

type AddrMap map[ssa.Value]map[ssa.Value]bool

func (m AddrMap) Put(key,val ssa.Value) {
	if m[key] == nil{
		m[key] = make(map[ssa.Value]bool)
	}

	m[key][val] = true
}

func (m AddrMap) Contains(key ssa.Value) bool {
	if m[key] == nil{
		return false
	}

	_,ok := m[key]
	return ok
}

func (m AddrMap) Delete(key ssa.Value)  {
	delete(m,key)
}

type SinkCallArgsMap map[ssa.CallInstruction]map[ssa.Value]bool

func (m SinkCallArgsMap) Put(key ssa.CallInstruction, val ssa.Value) {
	if m[key] == nil{
		m[key] = map[ssa.Value]bool{}
	}

	m[key][val] = true
}

type visitor struct {
	lattice  map[ssa.Value]*LatticeTag
	seen     map[*callgraph.Node]bool
	sinkArgs SinkCallArgsMap
	latticeSigParams map[*ssa.Function][]*LatticeTag
	ptrs AddrMap
}

func NewVisitor() *visitor {
	return &visitor{
		seen:     make(map[*callgraph.Node]bool),
		lattice:  make(map[ssa.Value]*LatticeTag),
		sinkArgs: map[ssa.CallInstruction]map[ssa.Value]bool{},
		latticeSigParams: make(map[*ssa.Function][]*LatticeTag),
		ptrs: map[ssa.Value]map[ssa.Value]bool{},
	}
}

func (v *visitor) Visit(node *callgraph.Node) {
	log.Infof("traverse visit: %s", node.String())
	if v.seen[node]{
		return
	}
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

	//if !v.seen[node] {
		//log.Infof("traverse visit: %s", node.String())
		//v.seen[node] = true
		//
		////check source
		//v.loopFunction(node)
		//
		//for _, outputEdge := range node.Out {
		//	if utils.IsSynthetic(outputEdge) || utils.InStd(outputEdge.Callee) || utils.InFabric(outputEdge.Callee) {
		//		v.seen[outputEdge.Callee] = true
		//		continue
		//	}
		//	log.Infof("out: %s", outputEdge.Callee.String())
		//
		//	//根据当前的taint的情况，设定函数入参的lattice情况
		//	v.Visit(outputEdge.Callee)
		//}

	//}else{
	//	log.Infof("node fn has visited before: %s", node.String())
	//}

}

func (v *visitor) loopFunction(node *callgraph.Node) {
	fn := node.Func
	if fn == nil {
		log.Errorf("fn == nil!")
		return
	}
	fn.WriteTo(logger.LogFile)


	//taint params
	for _,param := range fn.Params{
		if tags,ok := v.lattice[param];ok{
			log.Infof("taint from param, %s", param.String())
			if param.Referrers() == nil{
				continue
			}

			for _,ref := range *param.Referrers(){
				for tag,_ := range tags.msgSet {
					v.taint(ref,tag)
				}
			}
		}
	}

	//loop basic blocks and instructions
	for _, block := range fn.Blocks {
		for _, i := range block.Instrs {

			switch instr := i.(type) {
			case *ssa.UnOp:
				if _,ok := instr.X.(*ssa.Global); ok{
					log.Warningf("unop global variable here %s", prog.Fset.Position(i.Pos()))
					v.taint(i,"use global variable")
				}
			case *ssa.FieldAddr:
				addr := i.(ssa.Value)
				log.Debugf("put FieldAddr %p", addr)
				v.ptrs.Put(i.(ssa.Value),instr.X)
				if tags,ok := v.lattice[addr];ok{
					for tag,_ := range tags.msgSet {
						v.taintPointers(addr,tag)
					}
				}
			case *ssa.IndexAddr:
				//1. tainted -> ptr points to nothing
				//update ptr pointsToSet here
				//
				addr := i.(ssa.Value)
				log.Debugf("put IndexAddr %p", addr)
				v.ptrs.Put(addr,instr.X)
				if tags,ok := v.lattice[addr];ok{
					for tag,_ := range tags.msgSet {
						v.taintPointers(addr,tag)
					}
				}

			case *ssa.Range:
				ms := instr.X
				if _,ok := ms.Type().(*types.Map);ok{
					v.taint(instr,"range query map")
				}
				//2.
			case ssa.CallInstruction:

				v.checkSource(instr)
				v.checkSink(instr)

				//taint val; we need to taint other call signatures
				v.taintCallSigParams(instr)

			case *ssa.Return:
				v.handleReturnValue(node, instr)
			}
		}
	}
}

func (v *visitor) taintCallSigParams(callInstr ssa.CallInstruction) {
	if callInstr.Common().StaticCallee() == nil{
		return
	}
	args := callInstr.Common().Args
	fn := callInstr.Common().StaticCallee()
	params := fn.Params

	for i,arg := range args{
		if tags,ok := v.lattice[arg];ok{
			for tag,_ := range tags.msgSet {
				v.taintVal(params[i],tag)
				log.Infof("fn:%s, taint params:%s=%s", fn.Name(), arg.Name(),arg.String())
			}
		}
	}

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
			v.sinkArgs.Put(callInstr.(ssa.CallInstruction),arg)
			log.Infof("sink call-args:%s, %d", arg.Name(), len(v.sinkArgs))
		}
	}
}

func (v *visitor) taint(i ssa.Instruction, tag string) {
	switch val := i.(type) {
	//todo: stdlib call function taint flow
	case ssa.CallInstruction:
		if _,yes := config.IsSource(val); yes{
			if v.alreadyTaintedWithTag(val.(ssa.Value), tag) {
				return
			}
			v.taintVal(val.(ssa.Value), tag)
			v.taintReferrers(i, tag)
		}else{
			if yes := utils.IsStdCall(val); yes{
				log.Debugf("stdcall: %s", val.String())
				if v.alreadyTaintedWithTag(val.Value().Call.Value, tag) {
					return
				}
				v.taintVal(val.Value().Call.Value, tag)
				v.taintReferrers(i, tag)
				//args := val.Value().Call.Args
				//for _,arg := range args{
				//	if v.alreadyTainted(arg){
				//
				//	}
				//}
				//summ := summary.For(val)
				//if summ.IfTainted != 0{
				//
				//}
			}
		}



	case ssa.Value:
		if v.alreadyTaintedWithTag(val, tag) {
			return
		}
		v.taintVal(val, tag)
		v.taintReferrers(i, tag)
		//v.taintPointers(val,tag)
	case *ssa.Store:
		if v.alreadyTaintedWithTag(val.Addr, tag) {
			return
		}
		v.taintVal(val.Addr, tag)
		v.taintReferrers(val, tag)
		v.taintPointers(val.Addr,tag)
		//v.taintVal(val.Val, tag)

	default:
		return
	}

}

func (v *visitor) taintPointers(addr ssa.Value, tag string)  {
	log.Debugf("taint pointer: %s, %p",addr.Name(), addr)

	for val,_ := range v.ptrs[addr]{
		log.Debugf("addr points to %s=%s", val.Name(),val.String())
		if i,ok := val.(ssa.Instruction);ok{
			log.Debugf("taint pointer: %s=%s", val.Name(),val.String())
			v.taint(i,tag)
		}
	}
}

func (v *visitor) taintReferrers(i ssa.Instruction, tag string) {

	if val, ok := i.(ssa.Value); ok {
		if val.Referrers() == nil {
			return
		}
		for _, r := range *val.Referrers() {
			v.taint(r, tag)
		}
	} else if st, ok := i.(*ssa.Store); ok {
		addr := st.Addr
		log.Warningf("instr: %s is store instr", i.String())
		if addr.Referrers() == nil {
			return
		}

		for _, r := range *addr.Referrers() {
			log.Infof("store addr ref: %s", r.String())
			v.taint(r, tag)
		}
	} else{
		log.Warningf("instr: %s is not a value and store instr", i.String())
	}

}

func (v *visitor) taintVal(val ssa.Value, tag string) {

	if v.alreadyTaintedWithTag(val, tag) {
		return
	}
	log.Debugf("taintval: %s, %s, tag:%s", val.Name(), val.String(), tag)
	if v.lattice[val] == nil {
		v.lattice[val] = new(LatticeTag)
	}
	v.lattice[val].Add(tag)
}

func (v *visitor) handleSinkDetection() bool {
	outputResult = make(map[string]bool)
	log.Debugf("sink arg map len: %d", len(v.sinkArgs))
	for callInstr,m := range v.sinkArgs {
		for arg,_ := range m{
			log.Debugf("sink arg: %s=%s", arg.Name(),arg.String())
			if tags, ok := v.lattice[arg]; ok {
				//todo: report detection
				output := fmt.Sprintf("sink here %s with tag:%s ",prog.Fset.Position(callInstr.Pos()),tags.String())
				outputResult[output] = true
				//return true
			}
		}

	}

	for o,_ := range outputResult{
		log.Warning(o)
	}

	return false
}

func (v *visitor) handleReturnValue(node *callgraph.Node, retInstr *ssa.Return) {

	ins := node.In

	returnValues := retInstr.Results
	for _, result := range returnValues {
		if tags, found := v.lattice[result]; found {
			log.Debugf("lattice return value: %s", result.Name())
			for tag, _ := range tags.msgSet {
				for _, in := range ins {
					callsite := in.Site
					v.taint(callsite, tag)
				}
			}

		}
	}
}

func (v *visitor) alreadyTaintedWithTag(val ssa.Value, tag string) bool {

	if v.lattice[val] == nil {
		return false
	}

	if !v.lattice[val].Contains(tag){
		return false
	}

	return true
}


func (v *visitor) alreadyTainted(val ssa.Value) bool {
	tag,ok := v.lattice[val]

	if tag == nil{
		return false
	}

	return ok
}