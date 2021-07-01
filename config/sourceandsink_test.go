package config

import (
	"golang.org/x/tools/go/ssa"
	"testing"
)



type A struct {
	a int
	val ssa.Value
}

func TestReadConfig(t *testing.T) {
	//cfg,err := ReadConfig("../config/config.yaml")
	//if err != nil{
	//	t.Fatalf(err.Error())
	//}
	//log.Info(cfg.String())

	str := "123123"

	switch str {
	case "123123","123456":
		switch str {
		case "123123":
			print("123123\n")
		}
		print("123\n")
	case "123":
		print("123\n")
	}

	print("end\n")





}