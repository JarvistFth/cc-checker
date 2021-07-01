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

	var queue []*ssa.BasicBlock
	blocks := fn.DomPreorder()
	if len(blocks) == 0{
		return
	}
	log.Infof("bfs:%s ", fn.String())
	queue = append(queue,blocks[0])

	for _,blk := range blocks{
		log.Infof("put fn:%s, block idx:%s\n", blk.Parent().Name(),blk.String())
		listnode := NewWorklistNode(ctx,blk)
		Worklist.PushBack(listnode)
		if blk.Index == 0{

			taintbit := ctx.EntryValue
			pos := 0

			for taintbit != 0{
				if(taintbit & 1) == 1{
					BlockContexts[blk].AddIn(fn.Params[pos])
				}
				taintbit = taintbit >> 1
				pos++
			}
		}
	}

	//for len(queue) != 0{
	//	frontblk := queue[0]
	//	queue = queue[1:]
	//
	//
	//	listnode := NewWorklistNode(ctx,frontblk)
	//	Worklist.PushBack(listnode)
	//
	//	log.Infof("put fn:%s, block idx:%s\n", frontblk.Parent().Name(),frontblk.String())
	//
	//
	//	//set entry node.in
	//
	//	if frontblk.Index == 0{
	//
	//		taintbit := ctx.EntryValue
	//		pos := 0
	//
	//		for taintbit != 0{
	//			if(taintbit & 1) == 1{
	//				BlockContexts[frontblk].AddIn(fn.Params[pos])
	//			}
	//			taintbit = taintbit >> 1
	//			pos++
	//		}
	//	}
	//
	//	queue = append(queue,frontblk.Succs...)
	//}

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
		BlockContexts[listnode.Blk] = NewBlockContext(listnode.Blk)
		if listnode.Blk.Index != 0{

			for _,pred := range listnode.Blk.Preds{

				//union pred.out with me.in
				log.Infof("fn:%s , %s", listnode.Blk.Parent().String(),listnode.Blk.String())

				if BlockContexts[pred] == nil{
					log.Infof("%s",pred.String())
					continue
				}

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