package core

import "golang.org/x/tools/go/ssa"

var BlockContexts map[*ssa.BasicBlock]*BlockContext

func init() {
	BlockContexts = make(map[*ssa.BasicBlock]*BlockContext)
}

type BlockContext struct {
	Blk *ssa.BasicBlock
	In map[ssa.Value]bool
	Out map[ssa.Value]bool
}

func (c *BlockContext) SetIn(in map[ssa.Value]bool)  {
	for k,_ := range in{
		c.In[k] = true
	}
}

func (c *BlockContext) SetOut(out map[ssa.Value]bool)  {
	for k,_ := range out{
		c.Out[k] = true
	}
}

func (c *BlockContext) AddIn(val ssa.Value) {
	c.In[val] = true
}

func (c *BlockContext) AddOut(val ssa.Value) {
	c.Out[val] = true
}

func NewBlockContext(blk *ssa.BasicBlock) *BlockContext {
	return &BlockContext{
		Blk: blk,
		In: map[ssa.Value]bool{},
		Out: map[ssa.Value]bool{},
	}
}

func (c *BlockContext) ExistedIn(val ssa.Value) bool {

	_,ok := c.In[val]
	return ok

}

func (c *BlockContext) ExistedOut(val ssa.Value) bool {

	_,ok := c.Out[val]
	return ok

}




