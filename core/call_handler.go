package core

import (
	summary "cc-checker/core/stdlib"
	"golang.org/x/tools/go/ssa"
)

//func isSink(call ssa.CallInstruction) bool{
//
//	var path,recv,name string
//
//	if call.Common().IsInvoke(){
//		path,recv,name = utils.DecomposeAbstractMethod(call.Common())
//	}else{
//		if call.Common().StaticCallee() != nil{
//			path,recv,name = utils.DecomposeFunction(call.Common().StaticCallee())
//		}
//	}
//	if config.IsSink(path,recv,name){
//		//todo: report sink here
//
//		return true
//	}
//	return false
//}
//
//func isExclude(call ssa.CallInstruction) bool {
//	var path,recv,name string
//
//	if call.Common().IsInvoke(){
//		path,recv,name = utils.DecomposeAbstractMethod(call.Common())
//	}else{
//		if call.Common().StaticCallee() != nil{
//			path,recv,name = utils.DecomposeFunction(call.Common().StaticCallee())
//		}
//	}
//	if config.IsExcluded(path,recv,name){
//		//todo: just return
//		return true
//	}
//	return false
//}

//func isSource(call ssa.CallInstruction) (string,bool){
//	var path,recv,name string
//
//	if call.Common().IsInvoke(){
//		path,recv,name = utils.DecomposeAbstractMethod(call.Common())
//	}else{
//		if call.Common().StaticCallee() != nil{
//			path,recv,name = utils.DecomposeFunction(call.Common().StaticCallee())
//		}
//
//		if tag,ok := config.IsSource(path,recv,name);ok{
//			//todo: report sink here
//			return tag,ok
//		}
//	}
//	return "", false
//}

// return: @params 1: is stdlib call?
//		   @params 2: has taint tag?
func isStdLibCall(callInstr ssa.CallInstruction, entry TaintParams) (bool, bool, string) {
	if builtin, ok := callInstr.Common().Value.(*ssa.Builtin); ok {
		hastaint, tag := checkBuiltIn(callInstr, builtin.Name(), entry)
		return true, hastaint, tag

	}

	summ := summary.For(callInstr)

	if summ == nil {
		log.Debugf("%s ,summ is nil\n", callInstr.String())
		return false, false, ""
	}

	var args []ssa.Value
	// For "invoke" calls, Value is the receiver
	if callInstr.Common().IsInvoke() {
		args = append(args, callInstr.Common().Value)
	}
	args = append(args, callInstr.Common().Args...)

	// Determine whether we need to propagate taint.
	tainted := int64(0)
	var ret string
	for idx, tag := range entry {
		tainted |= 1 << idx
		ret += tag
	}

	if (tainted & summ.IfTainted) == 0 {
		return true, false, ret
	}
	return true, true, ret
}

func checkBuiltIn(callInstr ssa.CallInstruction, builtinName string, entry TaintParams) (bool, string) {
	switch builtinName {
	// The values being appended cannot be tainted.
	case "append":
		// The slice argument needs to be tainted because if its underlying array has
		// enough remaining capacity, the appended values will be written to it.
		// The returned slice is tainted if either the slice argument or the values
		// are tainted, so we need to visit the referrers.
		// Only the first argument (dst) can be tainted. (The src cannot be tainted.)
		var args []ssa.Value
		// For "invoke" calls, Value is the receiver

		args = append(args, callInstr.Common().Args...)

		// Determine whether we need to propagate taint.
		if len(entry) == 0 {
			return false, ""
		}
		var ret string
		for _, tag := range entry {
			ret += " " + tag
		}
		return true, ret

	case "copy":
		var args []ssa.Value
		// For "invoke" calls, Value is the receiver

		args = append(args, callInstr.Common().Args...)
		if len(entry) == 0 {
			return false, ""
		}
		var ret string
		for _, tag := range entry {
			ret += " " + tag
		}
		return true, ret

		// Determine whether we need to propagate taint.
	// The builtin delete(m map[Type]Type1, key Type) func does not propagate taint.
	case "delete":
		return false, ""
	default:
		log.Errorf("unexpected built in func:%s", builtinName)
	}
	return false, ""
}
