package core

import (
	"cc-checker/config"
	"cc-checker/logger"
	"cc-checker/ssautils"
	"cc-checker/utils"
	"golang.org/x/tools/go/callgraph"
	"sync"
	"testing"
)

func TestVisitor_Visit(t *testing.T) {
	prog, allpkgs, err := ssautils.BuildSSA("../ccs/timerandomcc/")
	mainpkgs, err := ssautils.MainPackages(allpkgs)
	if err != nil {
		t.Fatalf(err.Error())
	}
	if len(mainpkgs) == 0 {
		log.Infof("%s", "empty mainpkgs")
	} else {
		log.Info(mainpkgs[0].String())
	}

	result := BuildCallGraph(mainpkgs)
	_, invokef := utils.FindInvokeMethod(prog, mainpkgs[0])
	//var putstatefn *ssa.Function
	//fn := mainpkgs[0].Func("set")
	var wg sync.WaitGroup
	defer wg.Wait()
	if invokef != nil {
		nd := result.CallGraph.Nodes[invokef]
		for _, out := range nd.Out {
			go func(out *callgraph.Edge) {
				wg.Add(1)
				v := NewVisitor()
				log.Infof("invoke out:%s", out.Callee.String())
				v.Visit(out.Callee)
				wg.Done()
			}(out)

		}
	} else {
		log.Infof("invoke func is nil\n")
	}

}

func TestVisitor_Visit2(t *testing.T) {
	var err error
	cfg, err = config.ReadConfig()
	if err != nil {
		log.Fatalf(err.Error())
	}
	log.Debug(cfg.String())
	prog, allpkgs, err = ssautils.BuildSSA("../ccs/timerandomcc/")
	mainpkgs, err := ssautils.MainPackages(allpkgs)
	if err != nil {
		t.Fatalf(err.Error())
	}
	if len(mainpkgs) == 0 {
		log.Infof("%s", "empty mainpkgs")
	} else {
		log.Info(mainpkgs[0].String())
	}

	result := BuildCallGraph(mainpkgs)
	_, invokef := utils.FindInvokeMethod(prog, mainpkgs[0])
	if invokef != nil {
		invokef.WriteTo(logger.LogFile)
		v := NewVisitor()
		v.Visit(result.CallGraph.Nodes[invokef])
		//nd := result.CallGraph.Nodes[invokef]
		//for _, out := range nd.Out{
		//	v := NewVisitor()
		//	log.Infof("invoke out:%s", out.Callee.String())
		//	v.Visit(out.Callee)
		//}
	} else {
		log.Infof("invoke func is nil\n")
	}
}
