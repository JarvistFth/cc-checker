package core

import (
	"cc-checker/config"
	"cc-checker/ssautils"
	"cc-checker/utils"
	"golang.org/x/tools/go/pointer"
	"golang.org/x/tools/go/ssa"
)

var prog *ssa.Program

var invokef *ssa.Function

var cfg *config.Config

var allpkgs []*ssa.Package

var result *pointer.Result

var outputResult map[string]bool

func Init() {
	cfg, err := config.ReadConfig()
	if err != nil {
		log.Fatalf(err.Error())
	}
	log.Debug(cfg.String())
	var allpkgs []*ssa.Package
	prog, allpkgs, err = ssautils.BuildSSA("../ccs/timerandomcc/")
	mainpkgs, err := ssautils.MainPackages(allpkgs)
	if err != nil {
		panic(err.Error())
	}
	if len(mainpkgs) == 0 {
		log.Warningf("%s", "empty mainpkgs")
	} else {
		log.Infof("mainPkg:%s", mainpkgs[0].String())
	}
	//result := BuildCallGraph(mainpkgs)
	//_,invokef := utils.FindInvokeMethod(prog,mainpkgs[0])
}

func Main() {
	cfg, err := config.ReadConfig()
	if err != nil {
		log.Fatalf(err.Error())
	}
	log.Info(cfg.String())

	var mains []*ssa.Package

	prog, mains, err = ssautils.BuildSSA("../ccs/timerandomcc/")
	if err != nil {
		log.Fatalf(err.Error())
	}

	_, invokef = utils.FindInvokeMethod(prog, mains[0])
	BuildCallGraph(mains)
	StartAnalysis(invokef)
}

func StartAnalysis(fn *ssa.Function) {

}
