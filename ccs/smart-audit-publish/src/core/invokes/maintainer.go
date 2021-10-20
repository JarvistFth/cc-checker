package invokes

import (
	"bytes"
	"cc-checker/ccs/smart-audit-publish/src/core/contract"
	"cc-checker/ccs/smart-audit-publish/src/core/orgnization"
	"fmt"
	"log"
	"strconv"
)

// 查询所有的审计运维人员信息
func GetMaintainersMain(context contract.Context) *contract.Response {
	countBuf, err := context.GetState(orgnization.MaintainerCountKey)
	if err != nil {
		return contract.Error(fmt.Sprintf("获取审计当事人信息出错，详细信息：%s", err.Error()))
	}
	count, err := strconv.ParseUint(string(countBuf), 10, 32)

	var buffer bytes.Buffer
	bArrayMemberAlreadyWritten := false
	buffer.WriteString(`{"result":[`)

	for i := uint64(0); i < count; i++ {
		maintainer := orgnization.Maintainer{
			Member: &orgnization.Member{ID: uint32(i)}}
		mBuf, err := context.GetState(maintainer.Key())
		if err != nil {
			return contract.Error(fmt.Sprintf(
				"获取审计维护人员信息出错，详细信息：%s", err.Error()))
		}
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString(string(mBuf))
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString(`]}`)
	log.Printf("Query result: %s", buffer.String())

	return &contract.Response{Payload: buffer.Bytes()}
}
