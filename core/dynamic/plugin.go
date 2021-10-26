package dynamic

import (
	"cc-checker/logger"
	"flag"
	"os"
	"os/exec"
	"text/template"
)

var log = logger.GetLogger()


var gocode = `// Code generated, DO NOT EDIT.
	
package main
import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"math/rand"
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
	CC := new({{.ChaincodeStruct}})

	stub := shim.NewMockStub("CC",CC)

	var strargs []string
	
	Args := make([][]byte,len(strargs))
	
	for i:=0 ; i<len(strargs); i++{
		Args[i] = []byte(strargs[i])
	}

	stub.MockInit(genUUID(),Args)
	
	stub.MockInvoke(genUUID(),Args)
}


`

func genCode() {
	tp,err := template.New("CC-dynamic-analysis").Parse(gocode)
	dest := "dynamic.go"
	var file *os.File
	if file,err = os.Create(dest); err != nil{
		if !os.IsExist(err){
			panic(err.Error())
		}

		os.RemoveAll(dest)
	}
	vals := map[string]string{
		"ChaincodeStruct": os.Args[2],
	}
	tp.Execute(file,vals)
	file.Close()
}

//go:generate
func main() {


	ccTypeName := flag.String("typename","SimpleAssets","chaincode struct name")
	flag.Parse()
	exec.Command("go","generate","./",*ccTypeName)



}

func MockInvoke() {
	//CC := new(globalcc)
}