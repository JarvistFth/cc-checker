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