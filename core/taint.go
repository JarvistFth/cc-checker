package core
//
//import "golang.org/x/tools/go/ssa"
//
//// TaintTable node with taint-tag
//var TaintTable map[ssa.Value]string
//
//func init() {
//	TaintTable = map[ssa.Value]string{}
//}
//
//
//func(c *CallContext) TaintReferers(val ssa.Value, tag string){
//	ref := val.Referrers()
//	if ref == nil{
//		return
//	}
//
//	for _, instr := range *ref{
//		c.TaintNode(instr.(ssa.Node),tag)
//	}
//}
//
//func(c *CallContext) TaintNode(node ssa.Node, tag string) {
//	val,_ := node.(ssa.Value)
//	//TaintTable[val] = tag
//	switch n := node.(type) {
//	case *ssa.Alloc:
//		// dont need taint
//		return
//
//	case *ssa.Call:
//
//		//check stdlib
//
//		//check builtin
//
//		//not two above, initialize a new function
//		//1. init callctx with entryTaint
//		//2. if callctx has initialzed, hashset me.valueTaintBit = callctx.exitvalue
//		//3. callctx not initialized, init it, analyze it.
//		fn := extractCallee(n)
//		if fn == nil{
//			// invoke mode, don't get in
//			break
//		}
//		args := n.Call.Args
//		taintp := NewTaintParams()
//
//		for i,arg := range args{
//			if tag,ok := c.taintSet[arg];ok{
//				taintp[i] = tag
//			}
//		}
//
//
//		//entryTaint := taintArgs(c.fn.Params,taintp)
//
//
//		var callctx *CallContext
//		var ok bool
//		if callctx,ok = AnalyzedCall(fn,taintp);ok{
//			// has initialzed, hashset its exit value to me.value
//
//			exit := CallContexts[callctx.key].exit
//
//			for _,tag := range exit{
//				c.taintSet[n.Value()] += tag + " "
//			}
//
//
//		}else {
//			callctx = NewCallContext(fn,taintp)
//			callctx.Initialize()
//		}
//
//
//
//
//
//
//
//
//
//	case *ssa.Go:
//		// init new ctx
//
//	case *ssa.Field:
//
//	case *ssa.FieldAddr:
//
//	// Only the Addr (the Value that is being written to) should be visited.
//	case *ssa.Store:
//
//	// Only the Map itself can be tainted by an Update.
//	// The Key can't be tainted.
//	// The Value can propagate taint to the Map, but not receive it.
//	// MapUpdate has no referrers, it is only an Instruction, not a Value.
//	case *ssa.MapUpdate:
//
//
//	case *ssa.Select:
//		//todo support select
//
//	// The only Operand that can be tainted by a Send is the Chan.
//	// The Value can propagate taint to the Chan, but not receive it.
//	// Send has no referrers, it is only an Instruction, not a Value.
//	case *ssa.Send:
//
//
//		// This allows taint to propagate backwards into the sliced value
//		// when the resulting value is tainted
//
//	// These nodes' operands should not be visited, because they can only receive
//	// taint from their operands, not propagate taint to them.
//	case *ssa.BinOp,*ssa.IndexAddr,*ssa.Lookup:
//		var x,y ssa.Value
//		switch nn := n.(type) {
//		case *ssa.BinOp:
//			x = nn.X
//			y = nn.Y
//		case *ssa.IndexAddr:
//			x = nn.X
//			y = nn.Index
//		case *ssa.Lookup:
//			x = nn.X
//			y = nn.Index
//		}
//
//
//
//	case *ssa.ChangeInterface, *ssa.ChangeType, *ssa.Convert, *ssa.Extract, *ssa.Index, *ssa.MakeInterface, *ssa.Next, *ssa.Range, *ssa.Slice, *ssa.TypeAssert:
//
//		var x ssa.Value
//		switch nn := n.(type) {
//		case *ssa.ChangeInterface:
//			x = nn.X
//		case *ssa.ChangeType:
//			x = nn.X
//		case *ssa.Convert:
//			x = nn.X
//		case *ssa.Extract:
//			x = nn.Tuple
//		case *ssa.Index:
//			x = nn.X
//		case *ssa.MakeInterface:
//			x = nn.X
//		case *ssa.Next:
//			x = nn.Iter
//		case *ssa.Range:
//			x = nn.X
//		case *ssa.Slice:
//			x = nn.X
//		case *ssa.TypeAssert:
//			x = nn.X
//		}
//
//
//	// These nodes don't have operands; they are Values, not Instructions.
//
//	case *ssa.UnOp:
//		x := n.X
//
//
//	case *ssa.Return:
//		//flow results to currentctx.exitvalue
//
//
//	// These nodes cannot propagate taint.
//	case  *ssa.DebugRef, *ssa.Defer, *ssa.If, *ssa.Jump, *ssa.MakeClosure, *ssa.Panic, *ssa.RunDefers:
//
//	default:
//		//todo handle phi
//		//log.Errorf("unexpected node received: %T %v; please report this issue\n", n, n)
//
//	}
//}
