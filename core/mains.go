package core

import (
	"cc-checker/config"
	"cc-checker/ssautils"
	"flag"
	"golang.org/x/tools/go/pointer"
	"golang.org/x/tools/go/ssa"
)

var prog *ssa.Program

var invokef *ssa.Function

var cfg *config.Config

var allpkgs []*ssa.Package

var result *pointer.Result

var outputResult map[string]bool


func Main() {
	var err error
	cfg, err = config.ReadConfig()
	if err != nil {
		log.Fatalf(err.Error())
	}
	log.Debug(cfg.String())

	ccsPath := flag.String("path","./","chaincode package path")
	flag.Parse()
	log.Infof("build chaincode path=%s",*ccsPath)
	//prog, allpkgs, err = ssautils.BuildSSA("../ccs/hello/")
	prog, allpkgs, err = ssautils.BuildSSA(*ccsPath)
	mainpkgs, err := ssautils.MainPackages(allpkgs)
	if err != nil {
		panic(err.Error())
	}
	if len(mainpkgs) == 0 {
		log.Infof("%s", "empty mainpkgs")
	} else {
		log.Info(mainpkgs[0].String())
	}
	//_, invokef := utils.FindInvokeMethod(prog, mainpkgs[0])
	//
	//result = BuildCallGraph(mainpkgs,invokef)
	//for v,p := range result.Queries{
	//	log.Debugf("val: %s=%s, ptr:%s", v.Name(),v.String(), p.PointsTo().String())
	//}
	//if invokef != nil {
	//	v := NewVisitor()
	//	v.Visit(result.CallGraph.Nodes[invokef])
	//
	//	//for val,pt := range result.Queries{
	//	//	log.Debugf("value: %s=%s, points to:%s", val.Name(),val.String(), pt.PointsTo().String())
	//	//}
	//	//nd := result.CallGraph.Nodes[invokef]
	//	//for _, out := range nd.Out{
	//	//	v := NewVisitor()
	//	//	log.Infof("invoke out:%s", out.Callee.String())
	//	//	v.Visit(out.Callee)
	//	//}
	//
	//	v.handleSinkDetection()
	//} else {
	//	log.Infof("invoke func is nil\n")
	//}
}
