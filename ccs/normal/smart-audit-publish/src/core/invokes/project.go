package invokes

import (
	contract2 "cc-checker/ccs/normal/smart-audit-publish/src/core/contract"
	project2 "cc-checker/ccs/normal/smart-audit-publish/src/core/project"
	record2 "cc-checker/ccs/normal/smart-audit-publish/src/core/record"
	"fmt"
	"log"
	"strconv"
)

// 注册项目，返回项目ID
func RegisterProjectMain(args []string, context contract2.Context) *contract2.Response {
	p, err := project2.FromStrings(args, context)
	if err != nil {
		return contract2.Error(fmt.Sprint("解析审计业务失败，详细信息：", err))
	}

	if err = record2.StoreItem(p, context); err != nil {
		return contract2.Error(fmt.Sprintf("审计业务%s存储失败，详细信息：%s", p.Key(), err))
	}

	if err = record2.StoreCount(p, context); err != nil {
		return contract2.Error(fmt.Sprintf("审计业务%s相应的索引值存储失败，详细信息：%s",
			p.Key(), err))
	}

	log.Println("项目录入成功，项目ID：", p.ID)
	return &contract2.Response{
		Payload: []byte(strconv.FormatUint(uint64(p.ID), 32)),
	}
}

// 根据项目ID获取项目信息
func GetProjectMain(args []string, context contract2.Context) *contract2.Response {
	if len(args) == 0 {
		return contract2.Error("查询失败，需要提供项目ID")
	}

	projectID, err := strconv.ParseUint(args[0], 10, 32)
	if err != nil {
		return contract2.Error(fmt.Sprintf("解析项目ID出错，详细信息：%s", err.Error()))
	}

	pj := project2.Project{ID: uint32(projectID)}
	projectBuf, err := context.GetState(pj.Key())
	if err != nil {
		return contract2.Error(fmt.Sprintf("获取项目信息出错，详细信息：%s", err.Error()))
	}
	return &contract2.Response{Payload: projectBuf}
}
