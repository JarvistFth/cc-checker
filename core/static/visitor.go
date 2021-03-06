package static

import (
	"cc-checker/config"
	"cc-checker/core/common"
	"cc-checker/logger"
	"cc-checker/utils"
	"go/types"
	"golang.org/x/tools/go/callgraph"
	"golang.org/x/tools/go/ssa"
	"sync"
)



type visitor struct {
	lattice          map[ssa.Value]*common.LatticeTag
	seen             map[*callgraph.Node]bool
	sinkArgs         common.SinkCallArgsMap
	rwMaps           common.ReadAfterWriteMap
	latticeSigParams map[*ssa.Function][]*common.LatticeTag
	ptrs             common.AddrMap
}

func NewVisitor() *visitor {
	return &visitor{
		seen:             make(map[*callgraph.Node]bool),
		lattice:          make(map[ssa.Value]*common.LatticeTag),
		sinkArgs:         map[ssa.CallInstruction]map[ssa.Value]bool{},
		rwMaps:           common.ReadAfterWriteMap{},
		latticeSigParams: make(map[*ssa.Function][]*common.LatticeTag),
		ptrs:             map[ssa.Value]map[ssa.Value]bool{},
	}
}

func (v *visitor) Visit(node *callgraph.Node) {
	log.Infof("traverse visit: %s", node.String())
	if v.seen[node]{
		return
	}
	v.seen[node] = true

	//check source
	v.loopFunction(node)
	v.rwMaps = map[*ssa.Function]string{}
	var wg sync.WaitGroup
	for _, outputEdge := range node.Out {
		wg.Add(1)
		go func(outputEdge *callgraph.Edge) {
			defer wg.Done()
			if utils.IsSynthetic(outputEdge) || utils.InStd(outputEdge.Callee) || utils.InFabric(outputEdge.Callee) {
				v.seen[outputEdge.Callee] = true
				return
			}
			log.Infof("fn:%s -> out: %s", node.Func.Name(),outputEdge.Callee.String())

			//根据当前的taint的情况，设定函数入参的lattice情况
			v.taintCallSigParams(outputEdge.Site)
			v.Visit(outputEdge.Callee)
		}(outputEdge)
	}

	//for _, outputEdge := range node.Out {
	//	if utils.IsSynthetic(outputEdge) || utils.InStd(outputEdge.Callee) || utils.InFabric(outputEdge.Callee) {
	//		v.seen[outputEdge.Callee] = true
	//		return
	//	}
	//	log.Infof("fn:%s -> out: %s", node.Func.Name(),outputEdge.Callee.String())
	//
	//	//根据当前的taint的情况，设定函数入参的lattice情况
	//	v.taintCallSigParams(outputEdge.Site)
	//	v.Visit(outputEdge.Callee)
	//	//log.Warning("end visit here")
	//}

	//if !v.seen[node] {
		//log.Infof("traverse visit: %s", node.String())
		//v.seen[node] = true
		//
		////check source
		//v.loopFunction(node)
		//
		//for _, outputEdge := range node.Out {
		//	if utils.IsSynthetic(outputEdge) || utils.InStd(outputEdge.Callee) || utils.InFabric(outputEdge.Callee) {
		//		v.seen[outputEdge.Callee] = true
		//		continue
		//	}
		//	log.Infof("out: %s", outputEdge.Callee.String())
		//
		//	//根据当前的taint的情况，设定函数入参的lattice情况
		//	v.Visit(outputEdge.Callee)
		//}

	//}else{
	//	log.Infof("node fn has visited before: %s", node.String())
	//}

}

func (v *visitor) loopFunction(node *callgraph.Node) {
	fn := node.Func
	if fn == nil {
		log.Errorf("fn == nil!")
		return
	}
	fn.WriteTo(logger.LogFile)


	//taint params
	for _,param := range fn.Params{
		if tags,ok := v.lattice[param];ok{
			log.Infof("taint from param, %s", param.String())
			if param.Referrers() == nil{
				continue
			}

			for _,ref := range *param.Referrers(){
				for tag,_ := range tags.MsgSet {
					v.taint(ref,tag)
				}
			}
		}
	}

	//loop basic blocks and instructions
	for _, block := range fn.Blocks {
		for _, i := range block.Instrs {

			switch instr := i.(type) {
			case *ssa.UnOp:
				if g,ok := instr.X.(*ssa.Global); ok{
					//log.Warningf("unop global variable here %s", prog.Fset.Position(i.Pos()))
					v.taint(i,"use global variable " + g.Name() + " ")
				}
			case *ssa.FieldAddr:
				addr := i.(ssa.Value)
				log.Debugf("put FieldAddr %p", addr)
				v.ptrs.Put(i.(ssa.Value),instr.X)
				if tags,ok := v.lattice[addr];ok{
					for tag,_ := range tags.MsgSet {
						v.taintPointers(addr,tag)
					}
				}
			case *ssa.IndexAddr:
				//1. tainted -> ptr points to nothing
				//update ptr pointsToSet here
				//
				addr := i.(ssa.Value)
				log.Debugf("put IndexAddr %p", addr)
				v.ptrs.Put(addr,instr.X)
				if tags,ok := v.lattice[addr];ok{
					for tag,_ := range tags.MsgSet {
						v.taintPointers(addr,tag)
					}
				}

			case *ssa.Range:
				ms := instr.X
				if _,ok := ms.Type().(*types.Map);ok{
					v.taint(instr,"range query map")
				}
				//2.
			case ssa.CallInstruction:

				v.checkSource(instr)
				v.checkSink(instr)


				v.checkReadAfterWrite(instr)
				v.checkRangeQueryAndCrossChannel(instr)

				//taint val; we need to taint other call signatures
				v.taintCallSigParams(instr)

			case *ssa.Return:
				v.handleReturnValue(node, instr)
			}
		}
	}
}


// check source :
// if isSource, put call.value into taintSet
func (v *visitor) checkSource(callInstr ssa.Instruction) {
	if tag, yes := config.IsSource(callInstr.(ssa.CallInstruction)); yes {
		log.Infof("source fn here: %s", prog.Fset.Position(callInstr.Pos()))
		v.taint(callInstr, tag)
	}
}

func (v *visitor) checkSink(callInstr ssa.Instruction) {
	if ok := config.IsSink(callInstr.(ssa.CallInstruction)); ok {
		log.Infof("sink fn here: %s", prog.Fset.Position(callInstr.Pos()))
		for _, arg := range callInstr.(ssa.CallInstruction).Value().Call.Args {
			v.sinkArgs.Put(callInstr.(ssa.CallInstruction),arg)
			log.Infof("sink call-args:%s, %d", arg.Name(), len(v.sinkArgs))
		}
	}
}

func (v *visitor) checkReadAfterWrite(callInstr ssa.CallInstruction) {

	if ok := config.IsCCRead(callInstr); ok{
		//todo: update read function as sink-check
		keyName := callInstr.Common().Args[0].String()
		log.Infof("chaincode reads stub here: %s, key:%s", prog.Fset.Position(callInstr.Pos()), keyName)
		if key,ok := v.rwMaps.Contains(callInstr.Parent()); ok{
			if key == keyName{
				log.Warningf("read after write here:%s", prog.Fset.Position(callInstr.Pos()))
			}
			//os.Stdout.WriteString(fmt.Sprintf("read after write here:%s", prog.Fset.Position(callInstr.Pos())))
		}
	}

	if ok := config.IsCCWrite(callInstr); ok{
		refs := callInstr.Common().Args[0].Referrers()
		if refs != nil{
			for _, instr := range *refs{
				log.Debugf("taint read after write:%s", instr.String())
				v.taint(instr,"read after write")
			}
		}
		keyName := callInstr.Common().Args[0].String()
		log.Infof("chaincode writes stub here: %s, parent:%s, key:%s", prog.Fset.Position(callInstr.Pos()), callInstr.Parent(), keyName)
		v.rwMaps.Put(callInstr.Parent(), keyName)
	}
}

func (v *visitor) checkRangeQueryAndCrossChannel(callInstr ssa.CallInstruction) {
	if ok := config.IsRangeQueryCall(callInstr); ok{
		//log.Warningf("range query photon reads here: %s, key:%s", prog.Fset.Position(callInstr.Pos()),callInstr.Common().Args[0].String())
		//os.Stdout.WriteString(fmt.Sprintf("range query photon reads here: %s, key:%s", prog.Fset.Position(callInstr.Pos()),callInstr.Common().Args[0].String()))
		v.taint(callInstr,"phantom reads")
	}

	if ok := config.IsCrossChannelCall(callInstr); ok{
		log.Warningf("cross channel invoke here: %s, key:%s", prog.Fset.Position(callInstr.Pos()), callInstr.Common().Args[0].String())
		//os.Stdout.WriteString(fmt.Sprintf("cross channel invoke here: %s, key:%s", prog.Fset.Position(callInstr.Pos()), callInstr.Common().Args[0].String()))
	}
}


