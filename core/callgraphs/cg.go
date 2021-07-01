package callgraphs

import (
	"cc-checker/logger"
	"golang.org/x/tools/go/pointer"
	"golang.org/x/tools/go/ssa"
)

var log = logger.GetLoggerWithSTD()

func BuildCallGraph(mainpkg []*ssa.Package) *pointer.Result {
	cfg := &pointer.Config{
			Mains:           mainpkg,
			BuildCallGraph: true,
		}
	result,err := pointer.Analyze(cfg)

	if err != nil{
		log.Errorf(err.Error())
		return nil
	}

	return result
}
