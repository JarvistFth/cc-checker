package core

import (
	"cc-checker/logger"
	"cc-checker/utils"
	"golang.org/x/tools/go/ssa"
	"reflect"
)




func InitContext(ctx *ValueContext){
	AddValueContext(ctx)
	bfs(ctx)

}

func bfs(ctx *ValueContext){
	fn := ctx.Method
	fn.WriteTo(logger.LogFile)

	queue := make([]*ssa.BasicBlock,0)
	queue = append(queue,fn.Blocks[0])

	for len(queue) != 0{
		frontblk := queue[0]
		queue = queue[1:]


		listnode := NewWorklistNode(ctx,frontblk)
		Worklist.PushBack(listnode)


		//set entry node.in

		if frontblk.Index == 0{

			taintbit := ctx.EntryValue
			pos := 0

			for taintbit != 0{
				if(taintbit & 1) == 1{
					BlockContexts[frontblk].AddIn(fn.Params[pos])
				}
				taintbit = taintbit >> 1
				pos++
			}
		}

		queue = append(queue,frontblk.Succs...)
	}

}


func StartAnalysis(fn *ssa.Function) {

	valueCtx := &ValueContext{
		Method:     fn,
		Callee:     nil,
		EntryValue: 0,
		ExitValue: 0,
	}
	InitContext(valueCtx)


	for !Worklist.Empty(){
		listnode := Worklist.RemoveFront()

		if listnode.Blk.Index != 0{

			for _,pred := range listnode.Blk.Preds{

				//union pred.out with me.in
				UnionOut(BlockContexts[listnode.Blk].In,BlockContexts[pred].Out)
			}
		}

		oldout := utils.DeepCopy(BlockContexts[listnode.Blk].Out)

		analyzeInBlock(listnode)

		if !reflect.DeepEqual(oldout,BlockContexts[listnode.Blk].Out){
			for _,suc := range listnode.Blk.Succs{
				newlistnode := NewWorklistNode(valueCtx,suc)
				Worklist.PushBack(newlistnode)
			}
		}

	}
}