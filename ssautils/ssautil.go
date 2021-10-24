package ssautils

import (
	"cc-checker/logger"
	"cc-checker/utils"
	"fmt"
	"go/build"
	"golang.org/x/tools/go/packages"
	"golang.org/x/tools/go/ssa"
	"golang.org/x/tools/go/ssa/ssautil"
	"os"
)

var log = logger.GetLogger()

func BuildSSA(path string) (*ssa.Program, []*ssa.Package, error) {
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
	prog.Build()
	//mainpkg, err := MainPackages(pkgs)

	//if err != nil{
	//	log.Fatalf(err.Error())
	//	return prog,nil,nil
	//}

	return prog, pkgs, nil
}

func MainPackages(pkgs []*ssa.Package) ([]*ssa.Package, error) {
	var mains []*ssa.Package
	for _, p := range pkgs {
		if p != nil && p.Pkg.Name() == "main" && p.Func("main") != nil {
			mains = append(mains, p)
		}
	}
	if len(mains) == 0 {
		return nil, fmt.Errorf("no main packages")
	}

	//check third-party pkg here
	imports := mains[0].Pkg.Imports()
	for _, imp := range imports{
		if !utils.IsStdPkg(imp) && !utils.IsFabricPkg(imp){
			log.Warningf("chaincode use third-party pkg here: %s", imp.String())
			os.Stdout.WriteString("chaincode use third-party pkg here: "+imp.String()+"\n")

		}
	}
	return mains, nil
}
