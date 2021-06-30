package core

import "golang.org/x/tools/go/ssa"

var CallGraphs map[*ssa.Function]*CallNodes

func init() {
	CallGraphs = map[*ssa.Function]*CallNodes{}
}

type CallNodes struct {
	Callers []*CallValueCtx
	Callees []*ValueContext
	Me	*ssa.Function
}

func NewCallNodes(me *ssa.Function) *CallNodes {
	return &CallNodes{Me: me}
}

func (n *CallNodes) AddCallers(caller ...*CallValueCtx) {
	n.Callers = append(n.Callers,caller...)
}

func (n *CallNodes) AddCallees(callee ...*ValueContext) {
	n.Callees = append(n.Callees,callee...)
}

type CallValueCtx struct {
	ctx *ValueContext
	blk *ssa.BasicBlock
}

func NewCallValueCtx(context *ValueContext, blk *ssa.BasicBlock) *CallValueCtx {
	return &CallValueCtx{
		ctx: context,
		blk: blk,
	}
}