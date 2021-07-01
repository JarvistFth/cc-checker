package core

import "golang.org/x/tools/go/ssa"

var CallGraphs map[*ssa.Function]*CallNodes

func init() {
	CallGraphs = map[*ssa.Function]*CallNodes{}
}

type CallNodes struct {
	Callers []*CallerWithBlock
	Callees []*ValueContext
	Me	*ssa.Function
}

func NewCallNodes(me *ssa.Function, ) *CallNodes {

	node := &CallNodes{Me: me}
	if _,ok := CallGraphs[me]; !ok{
		CallGraphs[me] = node
	}



	return node
}

func (n *CallNodes) AddCallers(caller ...*CallerWithBlock) {
	n.Callers = append(n.Callers,caller...)
}

func (n *CallNodes) AddCallees(callee ...*ValueContext) {
	n.Callees = append(n.Callees,callee...)
}



type CallerWithBlock struct {
	ctx *ValueContext
	blk *ssa.BasicBlock
}

func NewCallValueCtx(context *ValueContext, blk *ssa.BasicBlock) *CallerWithBlock {
	return &CallerWithBlock{
		ctx: context,
		blk: blk,
	}
}

