package main
import (
    //"fmt"
	"testing"
	"github.com/hyperledger/fabric/protos/msp"
    "github.com/golang/protobuf/proto"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/suite"
)

type TS struct {
    suite.Suite
	stub *shim.MockStub
}

const certWithAttrs = `-----BEGIN CERTIFICATE-----
MIIB6TCCAY+gAwIBAgIUHkmY6fRP0ANTvzaBwKCkMZZPUnUwCgYIKoZIzj0EAwIw
GzEZMBcGA1UEAxMQZmFicmljLWNhLXNlcnZlcjAeFw0xNzA5MDgwMzQyMDBaFw0x
ODA5MDgwMzQyMDBaMB4xHDAaBgNVBAMTE015VGVzdFVzZXJXaXRoQXR0cnMwWTAT
BgcqhkjOPQIBBggqhkjOPQMBBwNCAATmB1r3CdWvOOP3opB3DjJnW3CnN8q1ydiR
dzmuA6A2rXKzPIltHvYbbSqISZJubsy8gVL6GYgYXNdu69RzzFF5o4GtMIGqMA4G
A1UdDwEB/wQEAwICBDAMBgNVHRMBAf8EAjAAMB0GA1UdDgQWBBTYKLTAvJJK08OM
VGwIhjMQpo2DrjAfBgNVHSMEGDAWgBTEs/52DeLePPx1+65VhgTwu3/2ATAiBgNV
HREEGzAZghdBbmlscy1NYWNCb29rLVByby5sb2NhbDAmBggqAwQFBgcIAQQaeyJh
dHRycyI6eyJhdHRyMSI6InZhbDEifX0wCgYIKoZIzj0EAwIDSAAwRQIhAPuEqWUp
svTTvBqLR5JeQSctJuz3zaqGRqSs2iW+QB3FAiAIP0mGWKcgSGRMMBvaqaLytBYo
9v3hRt1r8j8vN0pMcg==
-----END CERTIFICATE-----
`

// setup will be run for all tests in the suite
func (suite *TS) SetupTest() {
	suite.stub = shim.NewMockStub("mockStub", new(CC))
	assert.NotNil(suite.T(), suite.stub, "MockStub creation failed")


    // make a creator identity
    //sid1 := &msp.SerializedIdentity{Mspid: "orgMSP",
    //   IdBytes: []byte(certUser1)}

	sid2 := &msp.SerializedIdentity{Mspid: "SampleOrg",
		IdBytes: []byte(certWithAttrs)}

    //b, err := proto.Marshal(sid1)
    b, err := proto.Marshal(sid2)
    assert.Nil(suite.T(), err, "")
    suite.stub.Creator = b

	// call the constructor
	result := suite.stub.MockInit("1",[][]byte{
		[]byte("init"),
		[]byte{}})
	assert.EqualValues(suite.T(), result.Status, shim.OK, "Init is not successful")
}

func (suite *TS) TestInit() { 


}

func (suite *TS) TestGetID() {
	// call put
	result := suite.stub.MockInvoke("1", [][]byte{
		[]byte("testGetID"),
		 })

	assert.EqualValues(suite.T(), shim.OK, result.Status, "Put failed")

}

func TestSuite(t *testing.T) {
    suite.Run(t, new(TS))
}

