package invokes

import (
	"bytes"
	"cc-checker/ccs/smart-audit-publish/src/core/common"
	"cc-checker/ccs/smart-audit-publish/src/core/contract"
	"cc-checker/ccs/smart-audit-publish/src/core/project"
	"cc-checker/ccs/smart-audit-publish/src/core/record"
	"cc-checker/ccs/smart-audit-publish/src/core/rules"

	"errors"
	"fmt"
	"log"
)

// 添加审计事件
func AddEventMain(args []string, context contract.Context) *contract.Response {
	// 解析审计事件输入参数
	registration, err := project.RegistrationFromString(args, context)
	if err != nil {
		return contract.Error(fmt.Sprint("合规事件登录失败，详细信息：", err))
	}
	// 验证审计事件是否合规
	if err = verify(registration, context); err != nil {
		return contract.Error(fmt.Sprint("合规事件数据验证失败，详细信息：", err))
	}
	// 存储审计事件
	if err = record.StoreItem(registration, context); err != nil {
		return contract.Error(fmt.Sprintf("合规事件%s存储失败，详细信息：%s",
			registration.Key(), err))
	}
	// 存储审计当事人规范对象所对应的存储条数
	if err = record.StoreCount(registration, context); err != nil {
		return contract.Error(fmt.Sprintf("合规事件%s相应的索引值存储失败，详细信息：%s",
			registration.Key(), err))
	}
	log.Println("审计事件录入成功, 审计事件ID:", registration.ID)
	return &contract.Response{Payload: []byte("OK")}
}

// 根据审计当事人ID、项目ID以及规则ID，获取所有审计当事人的审计事件
func QueryEventsMain(args []string, context contract.Context) *contract.Response {
	if len(args) < 2 {
		return contract.Error("查询失败，需要提查询供审计事件对应的当事人ID、项目ID以及规则ID")
	}

	// 根据传入的参数获取eventID
	eventID, err := project.GetEventID(args, context)
	if err != nil {
		return contract.Error(err.Error())
	}

	// 获取第几次录入信息
	index, err := record.GetRecordCount(project.GetRegistrationCountKey(*eventID), context)
	if err != nil {
		return contract.Error(fmt.Sprintf("获取审计事件第几次录入信息出错，详细信息：%s", err.Error()))
	}

	// 获取开始及结束Key，查询满足条件的所有记录
	result, err := getQueryEventResult(eventID, index, context)
	if err != nil {
		return contract.Error(err.Error())
	}
	return &contract.Response{Payload: result}
}

// 获取审计事件查询的最终结果，json格式
func getQueryEventResult(eventID *common.Uint256,
	index uint32, context contract.Context) ([]byte, error) {

	reg := project.Registration{
		AuditeeSpecification: project.AuditeeSpecification{ID: *eventID}}
	startKey := reg.Key()
	reg.Index = index
	endKey := reg.Key()
	resultsIterator, err := context.GetStateByRange(startKey, endKey)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("获取审计事件信息出错，详细信息：%s", err.Error()))
	}
	defer resultsIterator.Close()

	var buffer bytes.Buffer
	bArrayMemberAlreadyWritten := false
	buffer.WriteString(`{"result":[`)

	for resultsIterator.HasNext() {
		//获取迭代器中的每一个值
		_, value, err := resultsIterator.Next()
		if err != nil {
			return nil, errors.New("Fail")
		}
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		//将查询结果放入Buffer中
		buffer.WriteString(string(value))
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString(`]}`)
	log.Printf("Query result: %s", buffer.String())

	return buffer.Bytes(), nil
}

// 验证注册规则
func verify(registration *project.Registration, context contract.Context) error {
	if err := rules.ValidateRules(registration.Rule.ID, registration.Params,
		context); err != nil {
		return fmt.Errorf("合规事件%s规则验证失败，详细信息：%s", registration.ID, err)
	}

	return nil
}
