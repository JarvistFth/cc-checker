package core

import (
	"fmt"
	"golang.org/x/tools/go/ssa"
)

var Contexts map[string]*ValueContext

func init() {
	Contexts = map[string]*ValueContext{}
}


type ValueContext struct {
	Method *ssa.Function
	Callee *ssa.CallCommon
	EntryValue int64
	ExitValue int64

}

func NewValueContext(fn *ssa.Function, callee *ssa.CallCommon) *ValueContext {
	return &ValueContext{
		Method:     fn,
		Callee: callee,
		EntryValue: 0,
		ExitValue: 0,
	}
}

func AddValueContext(ctx *ValueContext) {
	key := ctx.genKey()
	Contexts[key] = ctx
}

func (c *ValueContext) genKey() string {
	key := c.Method.String()

	key += fmt.Sprintf("%d",c.EntryValue)

	return key
}

func (c *ValueContext) SetEntryValue() {

}

func (c *ValueContext) SetExitValue() {

}

func (c *ValueContext) AlreadyAnalyzed() bool {
	key := c.genKey()
	if _,ok := Contexts[key];ok{
		return true
	}
	return false
}


