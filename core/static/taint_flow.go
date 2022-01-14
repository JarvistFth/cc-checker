package static

import (
	"cc-checker/core/common"
	"golang.org/x/tools/go/callgraph"
	"golang.org/x/tools/go/ssa"
)

func (v *visitor) taintCallSigParams(callInstr ssa.CallInstruction) {
	log.Debugf("taint call sig params")
	if callInstr.Common().StaticCallee() == nil{
		return
	}
	args := callInstr.Common().Args
	fn := callInstr.Common().StaticCallee()
	params := fn.Params

	for i,arg := range args{
		if tags,ok := v.lattice[arg];ok{
			for tag,_ := range tags.MsgSet {
				v.taintVal(params[i],tag)
				log.Infof("fn:%s, taint params:%s=%s", fn.Name(), arg.Name(),arg.String())
			}
		}
	}

}

func (v *visitor) taint(i ssa.Instruction, tag string) {
	switch val := i.(type) {
	//todo: stdlib call function taint flow
	//case ssa.CallInstruction:
		//v.taintCallSigParams(val)
		//if _,yes := config.IsSource(val); yes{
		//	if v.alreadyTaintedWithTag(val.(ssa.Value), tag) {
		//		return
		//	}
		//	v.taintVal(val.(ssa.Value), tag)
		//	v.taintReferrers(i, tag)
		//}else{
		//	if yes := utils.IsStdCall(val); yes{
		//		log.Debugf("stdcall: %s", val.String())
		//		if v.alreadyTaintedWithTag(val.Value().Call.Value, tag) {
		//			return
		//		}
		//		v.taintVal(val.Value().Call.Value, tag)
		//		v.taintReferrers(i, tag)
		//		//args := val.Value().Call.Args
		//		//for _,arg := range args{
		//		//	if v.alreadyTainted(arg){
		//		//
		//		//	}
		//		//}
		//		//summ := summary.For(val)
		//		//if summ.IfTainted != 0{
		//		//
		//		//}
		//	}
		//}



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
		//todo: un-taint all val refers without taint tag
		//if !v.alreadyTainted(val.Addr){
		//	delete(v.lattice,val.Addr)
		//}
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
		log.Infof("instr: %s is store instr", i.String())
		if addr.Referrers() == nil {
			return
		}

		for _, r := range *addr.Referrers() {
			log.Infof("store addr ref: %s", r.String())
			v.taint(r, tag)
		}
	} else{
		log.Errorf("instr: %s is not a value and store instr", i.String())
	}

}

func (v *visitor) taintVal(val ssa.Value, tag string) {

	if v.alreadyTaintedWithTag(val, tag) {
		return
	}
	log.Debugf("taintval: %s, %s, tag:%s", val.Name(), val.String(), tag)
	if v.lattice[val] == nil {
		v.lattice[val] = new(common.LatticeTag)
	}
	v.lattice[val].Add(tag)
}

func (v *visitor) handleReturnValue(node *callgraph.Node, retInstr *ssa.Return) {

	ins := node.In

	returnValues := retInstr.Results
	for _, result := range returnValues {
		if tags, found := v.lattice[result]; found {
			log.Debugf("lattice return value: %s", result.Name())
			for tag, _ := range tags.MsgSet {
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

