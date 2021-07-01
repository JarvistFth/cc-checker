package core

import (
	"cc-checker/logger"
	"golang.org/x/tools/go/ssa"
)

var log = logger.GetLoggerWithSTD()
//type InNode struct {
//	ctx *ValueContext
//	blk *ssa.BasicBlock
//	Values map[ssa.Value]bool
//}
//
//type OutNode struct {
//	ctx *ValueContext
//	blk *ssa.BasicBlock
//	Values map[ssa.Value]bool
//}



func analyzeInBlock(node *WorklistNode){

	blk := node.Blk
	currentCtx := node.Ctx
	blkctx := BlockContexts[blk]
	currentfn := node.Blk.Parent()


	for _,instr := range blk.Instrs{
		//each instr should update me.out
		val,_ := instr.(ssa.Value)
		switch n := instr.(type) {
		case *ssa.Alloc:
			// dont need taint


		case *ssa.Call:
			//check source
			//check sink

			//<cuurentCtx,blk> -> callee's fn




			_,ok := reportSource(n)
			if ok{
				blkctx.AddOut(n.Value())
				break
			}
			reportSink(n)

			isstdcall,_ := checkStdLibCall(n,blkctx)
			if isstdcall{
				break
			}

			fn := extractCallee(n)
			if fn == nil{
				break
			}
			newValueCtx := NewValueContext(fn,n.Common())


			flowLocalToCalleeParams(newValueCtx,BlockContexts[blk],n)

			if existCtx,ok := Contexts[newValueCtx.genKey()];ok{
				newValueCtx.ExitValue = existCtx.ExitValue
			}else{

				// else x' not exist, call initContext(x')
				InitContext(newValueCtx)
			}
			// add callgraph me -> callee



			// if x' exists in contexts
			// exitvalue := x'.exitValue


			// let me.value.tainttag <- flows from exitvalue

			newExitValue := newValueCtx.ExitValue
			for newExitValue != 0{
				if (newExitValue & 1) == 1{
					blkctx.AddOut(n.Value())
				}
				newExitValue = newExitValue >> 1
			}
			// normal flow , update in and out

		case *ssa.Go:
			// init new ctx

		case *ssa.Field:
			if blkctx.ExistedOut(n.X){
				blkctx.AddOut(val)
			}
			//
		case *ssa.FieldAddr:
			if blkctx.ExistedIn(n.X){
				blkctx.AddOut(val)
			}

		// Everything but the actual integer Index should be visited.


		// Everything but the actual integer Index should be visited.

			if blkctx.ExistedIn(n.X) || blkctx.ExistedOut(n.X){
				blkctx.AddOut(val)
			}

		// Only the Addr (the Value that is being written to) should be visited.
		case *ssa.Store:

		// Only the Map itself can be tainted by an Update.
		// The Key can't be tainted.
		// The Value can propagate taint to the Map, but not receive it.
		// MapUpdate has no referrers, it is only an Instruction, not a Value.
		case *ssa.MapUpdate:
			if blkctx.ExistedIn(n.Value) || blkctx.ExistedOut(n.Value){
				blkctx.AddOut(n.Map)
				blkctx.AddOut(n.Key)
			}


		case *ssa.Select:
			//todo support select

		// The only Operand that can be tainted by a Send is the Chan.
		// The Value can propagate taint to the Chan, but not receive it.
		// Send has no referrers, it is only an Instruction, not a Value.
		case *ssa.Send:

			if blkctx.ExistedIn(n.X) || blkctx.ExistedOut(n.X){
				blkctx.AddOut(n.Chan)
			}

			// This allows taint to propagate backwards into the sliced value
			// when the resulting value is tainted

		// These nodes' operands should not be visited, because they can only receive
		// taint from their operands, not propagate taint to them.
		case *ssa.BinOp,*ssa.IndexAddr,*ssa.Lookup:
			var x,y ssa.Value
			switch nn := n.(type) {
			case *ssa.BinOp:
				x = nn.X
				y = nn.Y
			case *ssa.IndexAddr:
				x = nn.X
				y = nn.Index
			case *ssa.Lookup:
				x = nn.X
				y = nn.Index
			}

			if blkctx.ExistedOut(x) || blkctx.ExistedOut(y){
				blkctx.AddOut(val)
			}

		case *ssa.ChangeInterface, *ssa.ChangeType, *ssa.Convert, *ssa.Extract, *ssa.Index, *ssa.MakeInterface, *ssa.Next, *ssa.Range, *ssa.Slice, *ssa.TypeAssert:

			var x ssa.Value
			switch nn := n.(type) {
			case *ssa.ChangeInterface:
				x = nn.X
			case *ssa.ChangeType:
				x = nn.X
			case *ssa.Convert:
				x = nn.X
			case *ssa.Extract:
				x = nn.Tuple
			case *ssa.Index:
				x = nn.X
			case *ssa.MakeInterface:
				x = nn.X
			case *ssa.Next:
				x = nn.Iter
			case *ssa.Range:
				x = nn.X
			case *ssa.Slice:
				x = nn.X
			case *ssa.TypeAssert:
				x = nn.X
			}

			if blkctx.ExistedIn(x) || blkctx.ExistedOut(x){
				blkctx.AddOut(val)
			}
		// These nodes don't have operands; they are Values, not Instructions.

		case *ssa.UnOp:
			x := n.X
			if _,ok := x.(*ssa.Global);ok{
				blkctx.AddOut(val)
			}
			if blkctx.ExistedIn(x) || blkctx.ExistedOut(x){
				blkctx.AddOut(val)
			}

		case *ssa.Return:
			//flow results to currentctx.exitvalue
			tainted := int64(0)
			rets := n.Results
			for i, ret := range rets{
				if blkctx.ExistedOut(ret) {
					tainted |= 1 << i
				}
			}
			currentCtx.ExitValue = tainted

			//for all edges ⟨X′, c⟩ → X in transitions do
			//ADD(worklist, ⟨X′, c⟩)
			//find all call me valuectx with blockctx
			callnodes := CallGraphs.Nodes[currentfn].In//nil here?



			for _,callnode := range callnodes{
				BlockContexts[callnode.Site.Block()].
				Worklist.PushBack(NewWorklistNode(caller.ctx,callnode.Site.Block()))
			}


		// These nodes cannot propagate taint.
		case  *ssa.DebugRef, *ssa.Defer, *ssa.If, *ssa.Jump, *ssa.MakeClosure, *ssa.Panic, *ssa.RunDefers:

		default:
			log.Errorf("unexpected node received: %T %v; please report this issue\n", n, n)

		}


	}
}





func flowLocalToCalleeParams(newCallValueCtx *ValueContext, blkctx *BlockContext, callinstr ssa.CallInstruction) int64 {

	in := blkctx.In

	args := extractMethodArgs(callinstr)

	tainted := int64(0)

	for i,arg := range args{
		if _,ok := in[arg];ok{
			tainted |= 1 << i
		}
	}

	newCallValueCtx.EntryValue = tainted
	return tainted

}


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

