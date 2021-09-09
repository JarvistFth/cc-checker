package ssautils

import (
	"cc-checker/logger"
	"fmt"
	"golang.org/x/tools/go/packages"
	"golang.org/x/tools/go/ssa"
	"golang.org/x/tools/go/ssa/ssautil"
)

var log = logger.GetLogger()

func BuildSSA(path string) (*ssa.Program ,[]*ssa.Package, error){
	ssacfg := &packages.Config{
		Mode: packages.LoadAllSyntax,
		Dir:  path,
		Env:  nil,
		Fset: nil,
	}
	initial, err := packages.Load(ssacfg)
	if err != nil{
		log.Fatalf(err.Error())
		return nil,nil,nil
	}
	prog, pkgs := ssautil.Packages(initial,0)
	prog.Build()
	mainpkg, err := mainPackages(pkgs)

	if err != nil{
		log.Fatalf(err.Error())
		return prog,nil,nil
	}

	return prog,mainpkg,nil
}

func mainPackages(pkgs []*ssa.Package) ([]*ssa.Package, error) {
	var mains []*ssa.Package
	for _, p := range pkgs {
		if p != nil && p.Pkg.Name() == "main" && p.Func("main") != nil {
			mains = append(mains, p)
		}
	}
	if len(mains) == 0 {
		return nil, fmt.Errorf("no main packages")
	}
	return mains, nil
}
