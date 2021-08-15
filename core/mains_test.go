package core

import (
	"cc-checker/config"
	"cc-checker/logger"
	ssautils "cc-checker/ssautils"
	"cc-checker/utils"
	"golang.org/x/tools/go/ssa"
	"testing"
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

func TestMain(m *testing.M) {
	Main()
}

func initial() *ssa.Function {
	cfg,err := config.ReadConfig("../config/config.yaml")
	if err != nil{
		log.Fatalf(err.Error())
	}
	log.Info(cfg.String())

	prog,mains,err := ssautils.BuildSSA("../ccs/timerandomcc/")
	if err != nil{
		log.Fatalf(err.Error())
	}

	_,invokef := utils.FindInvokeMethod(prog,mains[0])


	return invokef

}
