package invokes

import (
	"bytes"
	contract2 "cc-checker/ccs/normal/smart-audit-publish/src/core/contract"
	orgnization2 "cc-checker/ccs/normal/smart-audit-publish/src/core/orgnization"
	"fmt"
	"log"
	"strconv"
)

// 查询所有的审计运维人员信息
func GetMaintainersMain(context contract2.Context) *contract2.Response {
	countBuf, err := context.GetState(orgnization2.MaintainerCountKey)
	if err != nil {
		return contract2.Error(fmt.Sprintf("获取审计当事人信息出错，详细信息：%s", err.Error()))
	}
	count, err := strconv.ParseUint(string(countBuf), 10, 32)

	var buffer bytes.Buffer
	bArrayMemberAlreadyWritten := false
	buffer.WriteString(`{"result":[`)

	for i := uint64(0); i < count; i++ {
		maintainer := orgnization2.Maintainer{
			Member: &orgnization2.Member{ID: uint32(i)}}
		mBuf, err := context.GetState(maintainer.Key())
		if err != nil {
			return contract2.Error(fmt.Sprintf(
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

	return &contract2.Response{Payload: buffer.Bytes()}
}
