package invokes

import (
	contract2 "cc-checker/core/dynamic/ccs/normal/smart-audit-publish/src/core/contract"
	orgnization2 "cc-checker/core/dynamic/ccs/normal/smart-audit-publish/src/core/orgnization"
	record2 "cc-checker/core/dynamic/ccs/normal/smart-audit-publish/src/core/record"
	"fmt"
	"log"
	"strconv"
)

// 注册审计当事人，反回值为审计当时人ID
func RegisterAuditeeMain(args []string, context contract2.Context) *contract2.Response {
	auditee, err := orgnization2.AuditeeFromString(args, context)
	if err != nil {
		return contract2.Error(fmt.Sprint("解析审计当事人失败，详细信息：", err))
	}

	if err = record2.StoreItem(auditee, context); err != nil {
		return contract2.Error(fmt.Sprintf("审计当事人%s存储失败，详细信息：%s", auditee.Key(), err))
	}
	if err = record2.StoreCount(auditee, context); err != nil {
		return contract2.Error(fmt.Sprintf("审计当事人%s相应的索引值存储失败，详细信息：%s",
			auditee.Key(), err))
	}

	log.Println("审计当事人录入成功，当事人ID：", auditee.ID)
	return &contract2.Response{
		Payload: []byte(strconv.FormatUint(uint64(auditee.ID), 32)),
	}
}

// 根据审计当事人ID获取当事人信息
func GetAuditeeMain(args []string, context contract2.Context) *contract2.Response {
	if len(args) == 0 {
		return contract2.Error("查询失败，需要提供审计当事人ID")
	}

	auditeeID, err := strconv.ParseUint(args[0], 10, 32)
	if err != nil {
		return contract2.Error(fmt.Sprintf("解析审计当事人ID出错，详细信息：%s", err.Error()))
	}

	auditee := orgnization2.Auditee{
		Member: &orgnization2.Member{ID: uint32(auditeeID)}}
	auditeeBuf, err := context.GetState(auditee.Key())
	if err != nil {
		return contract2.Error(fmt.Sprintf("获取审计当事人信息出错，详细信息：%s", err.Error()))
	}

	return &contract2.Response{Payload: auditeeBuf}
}
