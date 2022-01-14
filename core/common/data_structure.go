package common

import (
	"cc-checker/logger"
	"golang.org/x/tools/go/ssa"
)

var log = logger.GetLogger()

type AddrMap map[ssa.Value]map[ssa.Value]bool

func (m AddrMap) Put(key,val ssa.Value) {
	if m[key] == nil{
		m[key] = make(map[ssa.Value]bool)
	}

	m[key][val] = true
}

func (m AddrMap) Contains(key ssa.Value) bool {
	if m[key] == nil{
		return false
	}

	_,ok := m[key]
	return ok
}

func (m AddrMap) Delete(key ssa.Value)  {
	delete(m,key)
}

type SinkCallArgsMap map[ssa.CallInstruction]map[ssa.Value]bool

func (m SinkCallArgsMap) Put(key ssa.CallInstruction, val ssa.Value) {
	if m[key] == nil{
		m[key] = map[ssa.Value]bool{}
	}

	m[key][val] = true
}


type ReadAfterWriteMap map[*ssa.Function] string



func (m ReadAfterWriteMap) Put(parent *ssa.Function, key string) {
	m[parent] = key

}

func (m ReadAfterWriteMap) Delete(parent *ssa.Function) {
	delete(m, parent)
}

func (m ReadAfterWriteMap) Contains(parent *ssa.Function) (key string, ok bool) {
	key,ok = m[parent]
	return key,ok
}

type LatticeTag struct {
	Tag    string
	MsgSet map[string]bool
}

func (t *LatticeTag) Add(tag string) {
	if t.MsgSet == nil{
		t.MsgSet = make(map[string]bool)
	}
	if _,ok := t.MsgSet[tag];ok{
		return
	}

	t.Tag += tag + " | "
	t.MsgSet[tag] = true
}

func (t *LatticeTag) Contains(tag string) bool {
	_,ok :=t.MsgSet[tag]
	return ok
}

func (t *LatticeTag) Delete(tag string) {
	delete(t.MsgSet,tag)
}

func (t *LatticeTag) String() string {
	return t.Tag
}
