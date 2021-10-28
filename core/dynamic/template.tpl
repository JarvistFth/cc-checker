package dynamic


func main() {
	cc := new(SimpleAsset)

	var agents []*StubAgent

	var n = 2
	fn := "setWithRand"
	args := []string{"aaa"}
	b := MockParams(fn,args)

	for i:=0; i<n; i++{
		agent := NewStubAgent("cc",cc)
		agents = append(agents,agent)
	}



	for _,agent := range agents{
		agent.MockInvoke(genUUID(),b)
	}



}
