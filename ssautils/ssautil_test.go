package ssautils

import "testing"


func TestBuildSSA(t *testing.T) {
	_,mainpkg,err := BuildSSA("../ccs/timerandomcc/")
	if err != nil{
		t.Fatalf(err.Error())
	}
	if len(mainpkg)== 0{
		log.Infof("%s","empty mainpkg")
	}else{
		log.Info(mainpkg[0].String())
	}

}
