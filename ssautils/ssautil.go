package ssautils

import (
	"cc-checker/logger"
	"cc-checker/utils"
	"fmt"
	"go/build"
	"golang.org/x/tools/go/packages"
	"golang.org/x/tools/go/ssa"
	"golang.org/x/tools/go/ssa/ssautil"
	"time"
)

var log = logger.GetLogger()

func BuildSSA(path string) (*ssa.Program, []*ssa.Package, error) {
	startTime := time.Now().UnixNano()
	ssacfg := &packages.Config{
		Mode: packages.LoadAllSyntax,
		Dir:  path,
		//BuildFlags: []string{"-gcflags=-N -l"} ,
		BuildFlags: build.Default.BuildTags,
	}
	initial, err := packages.Load(ssacfg)
	if err != nil {
		log.Fatalf(err.Error())
		return nil, nil, nil
	}
	prog, pkgs := ssautil.AllPackages(initial, 0)

	if len(initial) == 0 {
		panic("no packages info!!")
	}
	//if len(initial[0].Syntax) == 0{
	//	panic("no ast-files info!!")
	//}
	//checkAst(prog.Fset, initial[0].Syntax[0], initial[0].TypesInfo)

	prog.Build()
	//mainpkg, err := MainPackages(pkgs)

	//if err != nil{
	//	log.Fatalf(err.Error())
	//	return prog,nil,nil
	//}
	endTime := time.Now().UnixNano()
	seconds := float64((endTime - startTime) / 1e9)
	ms := float64((endTime - startTime) / 1e6)
	log.Infof("build ssa use totaltime: %f s, %f ms",seconds, ms)
	fmt.Printf("build ssa use totaltime: %f s, %f ms\n",seconds, ms)
	return prog, pkgs, nil
}

func MainPackages(pkgs []*ssa.Package) ([]*ssa.Package, error) {
	var mains []*ssa.Package
	for _, p := range pkgs {
		if p != nil && p.Pkg.Name() == "main" {
			mains = append(mains, p)
		}
	}
	if len(mains) == 0 {
		return nil, fmt.Errorf("no main packages")
	}

	//check third-party pkg here
	imports := mains[0].Pkg.Imports()
	log.Debugf("main pkg path:%s", mains[0].Pkg.Path())
	for _, imp := range imports{
		if !utils.IsStdPkg(imp) && !utils.IsFabricPkg(imp) && !utils.IsInternalPkg(imp, mains[0].Pkg.Path()){
			log.Warningf("chaincode use third-party pkg here: %s", imp.String())
			//os.Stdout.WriteString("chaincode use third-party pkg here: "+imp.String()+"\n")

		}
	}

	return mains, nil
}
