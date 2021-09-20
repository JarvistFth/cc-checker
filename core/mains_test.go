package core

import (
	"cc-checker/config"
	"cc-checker/logger"
	ssautils "cc-checker/ssautils"
	"cc-checker/utils"
	"golang.org/x/tools/go/ssa"
)

func init() {
	logger.GetLogger()
}

//func TestInitContext(t *testing.T) {
//
//}
//
//func TestStartAnalysis(t *testing.T) {
//	invokef := initial()
//	StartAnalysis(invokef)
//
//}

func initial() *ssa.Function {
	cfg, err := config.ReadConfig()
	if err != nil {
		log.Fatalf(err.Error())
	}
	log.Info(cfg.String())

	prog, mains, err := ssautils.BuildSSA("../ccs/timerandomcc/")
	if err != nil {
		log.Fatalf(err.Error())
	}

	_, invokef := utils.FindInvokeMethod(prog, mains[0])

	return invokef

}
