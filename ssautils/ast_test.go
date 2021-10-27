package ssautils

import (
	"go/build"
	"golang.org/x/tools/go/packages"
	"golang.org/x/tools/go/ssa/ssautil"
	"testing"
)

func TestUnhandledError(t *testing.T) {
	path := "../ccs/timerandomcc/"
	ssacfg := &packages.Config{
		Mode: packages.LoadAllSyntax,
		Dir:  path,
		//BuildFlags: []string{"-gcflags=-N -l"} ,
		BuildFlags: build.Default.BuildTags,
	}
	initial, err := packages.Load(ssacfg)
	if err != nil {
		log.Fatalf(err.Error())
	}
	prog, _ := ssautil.AllPackages(initial, 0)
	checkAst(prog.Fset, initial[0].Syntax[0], initial[0].TypesInfo)
}
