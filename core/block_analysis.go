package core

import (
	"cc-checker/logger"
	"golang.org/x/tools/go/ssa"
)

var log = logger.GetLoggerWithSTD()




func extractCallee(call ssa.CallInstruction) (  fn *ssa.Function){
	if call.Common().IsInvoke(){
		return nil
	}else{
		if call.Common().StaticCallee() != nil{
			return call.Common().StaticCallee()
		}
	}
	return nil
}

func extractMethodArgs(call ssa.CallInstruction) []ssa.Value{
	var args []ssa.Value
	if call.Common().IsInvoke(){
		args = append(args,call.Common().Value)
	}else{
		args = append(args,call.Common().Args...)
	}

	return args

}

