package core

import (
	sll "github.com/emirpasic/gods/lists/singlylinkedlist"
	"golang.org/x/tools/go/ssa"
)


type WorkList struct {
	list *sll.List
}

var Worklist *WorkList


func init()  {
	NewWorkList()
}

func NewWorkList() *WorkList {
	Worklist = &WorkList{list: sll.New()}
	return Worklist
}

func (l *WorkList) PushBack(vals ...interface{}) {
	l.list.Append(vals)
}

func (l *WorkList) Front() *WorklistNode{
	ret,ok := l.list.Get(0)
	if ok{
		return ret.(*WorklistNode)
	}
	return nil
}

func (l *WorkList) RemoveFront() *WorklistNode {
	ret,ok := l.list.Get(0)
	if ok{
		l.list.Remove(0)
		return ret.(*WorklistNode)
	}
	return nil
}

func (l *WorkList) Empty() bool {
	return l.list.Empty()
}



type WorklistNode struct {

	Blk *ssa.BasicBlock
	Ctx *ValueContext
}

func NewWorklistNode(context *ValueContext, block *ssa.BasicBlock) *WorklistNode{


	return &WorklistNode{

		Blk: block,
		Ctx: context,
	}

}

func UnionOut(in,out map[ssa.Value]bool) {
	for k,v := range out{
		if _,ok := in[k];ok{
			continue
		}
		in[k] = v
	}
}