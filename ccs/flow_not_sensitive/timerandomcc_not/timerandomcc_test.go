package main

import (
	"cc-checker/core/dynamic"
	"math/rand"
	"testing"
	"time"
)

func genUUID() string {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	bytes := make([]byte, 36)
	for i := 0; i < 36; i++ {
		b := r.Intn(26) + 65
		bytes[i] = byte(b)
	}
	return string(bytes)
}

func TestMockInvoke(t *testing.T) {
	cc := new(SimpleAsset)

	agent := dynamic.NewStubAgent ("cc",cc)

	//var strargs []string
	//
	//args := make([][]byte,len(strargs))
	//
	//for i:=0 ; i<len(strargs); i++{
	//	args[i] = []byte(strargs[i])
	//}
	//
	//stub.MockInit(genUUID(),args)
	//
	//stub.MockInvoke(genUUID(),args)

	agent.MockInvoke("123",[][]byte{[]byte("a")})

}
