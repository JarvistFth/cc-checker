package static

import (
	"fmt"
	"golang.org/x/tools/go/ssa"
	"sort"
)


type TaintParams map[int]string
type TaintParamsVal map[ssa.Value]string

func NewTaintParams() TaintParams {
	return make(map[int]string)
}

func NewTaintParamsVal() TaintParamsVal {
	return make(map[ssa.Value]string)
}

func (p TaintParams) String() string {
	var idx []int
	for k,_ := range p{
		idx = append(idx,k)
	}
	sort.Ints(idx)
	var ret string
	for _, i := range idx{
		ret += fmt.Sprintf("%d",i)
	}
	return ret
}


func taintArgs(args []*ssa.Parameter ,taint TaintParams) TaintParamsVal {
	ret := NewTaintParamsVal()
	for idx,tag := range taint{
		ret[args[idx]] = tag
	}
	return ret
}
