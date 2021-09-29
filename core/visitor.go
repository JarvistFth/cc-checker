package core

import (
	"cc-checker/config"
	"cc-checker/logger"
	"cc-checker/utils"
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

type visitor struct {
	lattice  map[ssa.Value]*LatticeTag
	seen     map[*callgraph.Node]bool
	sinkArgs map[ssa.Value]bool
	latticeSigParams map[*ssa.Function][]*LatticeTag
	ptrs AddrMap
}

func NewVisitor() *visitor {
	return &visitor{
		seen:     make(map[*callgraph.Node]bool),
		lattice:  make(map[ssa.Value]*LatticeTag),
		sinkArgs: make(map[ssa.Value]bool),
		latticeSigParams: make(map[*ssa.Function][]*LatticeTag),
		ptrs: map[ssa.Value]map[ssa.Value]bool{},
	}
}

func (v *visitor) Visit(node *callgraph.Node) {
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

	for _,param := range fn.Params{
		if tags,ok := v.lattice[param];ok{
			if param.Referrers() == nil{
				continue
			}

			for _,ref := range *param.Referrers(){
				for tag,_ := range tags.hashset{
					v.taint(ref,tag)
				}
			}
		}
	}

	for _, block := range fn.Blocks {
		for _, i := range block.Instrs {

			switch instr := i.(type) {
			case *ssa.UnOp:
				if _,ok := instr.X.(*ssa.Global); ok{
					log.Warningf("unop global variable here %s", prog.Fset.Position(i.Pos()))
					v.taint(i,"use global variable")
				}
			case *ssa.FieldAddr:
				log.Debugf("put FieldAddr %p", i.(ssa.Value))
				v.ptrs.Put(i.(ssa.Value),instr.X)
			case *ssa.IndexAddr:
				//1. tainted -> ptr points to nothing
				//update ptr pointsToSet here
				//
				addr := i.(ssa.Value)
				log.Debugf("put IndexAddr %p", addr)
				v.ptrs.Put(addr,instr.X)
				if tags,ok := v.lattice[addr];ok{
					for tag,_ := range tags.hashset{
						v.taintPointers(addr,tag)
					}
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
			for tag,_ := range tags.hashset{
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
			log.Infof("sink call-args:%s", arg.Name())
			v.sinkArgs[arg] = true
		}
	}
}

func (v *visitor) taint(i ssa.Instruction, tag string) {
	switch val := i.(type) {
	case ssa.Value:
		if v.alreadyTainted(val, tag) {
			return
		}
		v.taintVal(val, tag)
		v.taintReferrers(i, tag)
		//v.taintPointers(val,tag)
	case *ssa.Store:
		if v.alreadyTainted(val.Addr, tag) {
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

	if v.alreadyTainted(val, tag) {
		return
	}
	log.Debugf("taintval: %s, %s, tag:%s", val.Name(), val.String(), tag)
	if v.lattice[val] == nil {
		v.lattice[val] = new(LatticeTag)
	}
	v.lattice[val].Add(tag)
}

func (v *visitor) handleSinkDetection() bool {
	for arg, _ := range v.sinkArgs {
		if tags, ok := v.lattice[arg]; ok {
			//todo: report detection
			log.Warningf("sink here %s sinkTag with tag:%s sinkTag", prog.Fset.Position(arg.Pos()), tags.String())
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

			for tag, _ := range tags.hashset {
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

	if !v.lattice[val].Contains(tag){
		return false
	}

	return true

}

func (v visitor) taintPtr() {

}
