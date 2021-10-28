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

	//shim.NewMockStub("Cc",cc)

	//st.GetState()

	//st.MockInvoke("1",)
	dynamic.CheckNonDitermined(cc,"setWithTime",[]string{"asd","bcd"})
	dynamic.CheckOtherRules(cc, "setWithTime",[]string{"1","2"})

}
