package core

import (
	"fmt"
	"golang.org/x/tools/go/ssa"
)

var CallContexts map[string]*CallContext

type CallContext struct {
	key   string
	fn    *ssa.Function
	entry TaintParams
	exit  TaintParamsVal
	taintSet map[ssa.Value]string
}

func AnalyzedCall(fn *ssa.Function, entry TaintParams) (*CallContext,bool) {
	ctx := CallContext{
		key:   genkey(fn,entry),
		fn:    fn,
		entry: entry,
		exit:  NewTaintParamsVal(),
		taintSet: map[ssa.Value]string{},
	}
	ret,ok := CallContexts[ctx.key]
	return ret,ok
}

func NewCallContext(fn *ssa.Function, entry TaintParams) *CallContext {


	ret := &CallContext{
		key:   genkey(fn,entry),
		fn:    fn,
		entry: entry,
		exit:  NewTaintParamsVal(),
	}
	CallContexts[ret.key] = ret
	return ret
}

func genkey(fn *ssa.Function, entry TaintParams) string {
	return fmt.Sprintf("%s, entry:%d",fn.String(), entry.String())
}

func (c *CallContext) SetEntry(entry TaintParams) {
	c.entry = entry
}

func (c *CallContext) SetExit() {

}

func (c *CallContext) Initialize() {
	fn := c.fn

	vals := taintArgs(c.fn.Params,c.entry)

	for val,tag := range vals{
		c.TaintReferers(val,tag)
	}

	for _,blk := range fn.Blocks{
		for _,instr := range blk.Instrs{
			// if is source , taint referrers
			switch n := instr.(type) {
			case *ssa.Call:
				if isExclude(n){
					break
				}
				if tag,yes := isSource(n);yes{
					c.TaintReferers(n.Value(),tag)
				}

				if isSink(n){
					// check its args has tained tag?
					if tag,ok := c.taintSet[n.Value()];ok{
						//todo report sink
						log.Warningf("sink here: %s with tag %s\n",prog.Fset.Position(n.Pos()), tag)
					}
				}

				//check stdlib call
			case *ssa.Return:

			}

		}
	}



}