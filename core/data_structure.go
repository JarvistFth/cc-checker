package core

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

type RwDetails struct {
	parents *ssa.Function
	key		ssa.Value
}

type ReadAfterWriteMap map[ssa.CallInstruction]RwDetails



func (m ReadAfterWriteMap) Put(c ssa.CallInstruction, rw RwDetails) {
	m[c] = rw
}

func (m ReadAfterWriteMap) Delete(c ssa.CallInstruction) {
	delete(m, c)
}

func (m ReadAfterWriteMap) Contains(c ssa.CallInstruction) (rw RwDetails, ok bool) {
	rw,ok = m[c]
	return rw,ok
}

type LatticeTag struct {
	tag    string
	msgSet map[string]bool
}

func (t *LatticeTag) Add(tag string) {
	if t.msgSet == nil{
		t.msgSet = make(map[string]bool)
	}
	if _,ok := t.msgSet[tag];ok{
		return
	}

	t.tag += tag + " | "
	t.msgSet[tag] = true
}

func (t *LatticeTag) Contains(tag string) bool {
	_,ok :=t.msgSet[tag]
	return ok
}

func (t *LatticeTag) Delete(tag string) {
	delete(t.msgSet,tag)
}

func (t *LatticeTag) String() string {
	return t.tag
}
