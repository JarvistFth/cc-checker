package core

import (
	"cc-checker/config"
	summary "cc-checker/core/stdlib"
	"cc-checker/utils"
	"golang.org/x/tools/go/ssa"
)

func reportSink(call *ssa.Call){

	var path,recv,name string

	if call.Common().IsInvoke(){
		path,recv,name = utils.DecomposeAbstractMethod(call.Common())
	}else{
		if call.Common().StaticCallee() != nil{
			path,recv,name = utils.DecomposeFunction(call.Common().StaticCallee())
		}

		if config.IsSink(path,recv,name){
			//todo: report sink here
		}

		if config.IsExcluded(path,recv,name){
			//todo: just return
			return
		}


	}
}

func reportSource(call *ssa.Call) (string,bool){
	var path,recv,name string

	if call.Common().IsInvoke(){
		path,recv,name = utils.DecomposeAbstractMethod(call.Common())
	}else{
		if call.Common().StaticCallee() != nil{
			path,recv,name = utils.DecomposeFunction(call.Common().StaticCallee())
		}

		if tag,ok := config.IsSource(path,recv,name);ok{
			//todo: report sink here
			return tag,ok
		}
	}
	return "", false
}

func checkStdLibCall(callInstr ssa.CallInstruction, blkctx *BlockContext) (bool, bool){
	summ := summary.For(callInstr)

	if summ == nil {
		log.Debugf("%s ,summ is nil\n",callInstr.String())
		return false,false
	}

	var args []ssa.Value
	// For "invoke" calls, Value is the receiver
	if callInstr.Common().IsInvoke() {
		args = append(args, callInstr.Common().Value)
	}
	args = append(args, callInstr.Common().Args...)

	// Determine whether we need to propagate taint.
	tainted := int64(0)
	for i, a := range args {
		if blkctx.ExistedOut(a.(ssa.Value)){
			tainted |= 1 << i
		}

	}

	if (tainted & summ.IfTainted) == 0 {
		return true,false
	}
	blkctx.AddOut(callInstr.Value())
	return true,true


	// Taint call arguments.
	//for _, i := range summ.TaintedArgs {
	//	prop.taint(args[i].(ssa.Node), maxInstrReached, lastBlockVisited, false)
	//}

	// Only actual Call instructions can have Referrers.
	//call, ok := callInstr.(*ssa.Call)
	//if !ok {
	//	return
	//}

	// If there are no referrers, exit early.
	//if call.Referrers() == nil {
	//	return
	//}

	// If the call has a single return value, the return value is the call
	// instruction itself, so if the call's return value is tainted, taint
	// the Referrers.
	//if call.Common().Signature().Results().Len() == 1 {
	//	if len(summ.TaintedRets) > 0 {
	//		//prop.taintReferrers(call, maxInstrReached, lastBlockVisited)
	//		blkctx.AddOut(call.Value())
	//	}
	//	return
	//}
	//
	//indexToExtract := map[int]*ssa.Extract{}
	//for _, r := range *call.Referrers() {
	//	e := r.(*ssa.Extract)
	//	indexToExtract[e.Index] = e
	//}
	//for i := range summ.TaintedRets {
	//	blkctx.AddOut(call.Value())
	//}
}