package invokes

import (
	"bytes"
	"cc-checker/core/dynamic/ccs/normal/smart-audit-publish/src/core/common"
	contract2 "cc-checker/core/dynamic/ccs/normal/smart-audit-publish/src/core/contract"
	project2 "cc-checker/core/dynamic/ccs/normal/smart-audit-publish/src/core/project"
	record2 "cc-checker/core/dynamic/ccs/normal/smart-audit-publish/src/core/record"
	"cc-checker/core/dynamic/ccs/normal/smart-audit-publish/src/core/rules"
	"errors"
	"fmt"
	"log"
)

// 添加审计事件
func AddEventMain(args []string, context contract2.Context) *contract2.Response {
	// 解析审计事件输入参数
	registration, err := project2.RegistrationFromString(args, context)
	if err != nil {
		return contract2.Error(fmt.Sprint("合规事件登录失败，详细信息：", err))
	}
	// 验证审计事件是否合规
	if err = verify(registration, context); err != nil {
		return contract2.Error(fmt.Sprint("合规事件数据验证失败，详细信息：", err))
	}
	// 存储审计事件
	if err = record2.StoreItem(registration, context); err != nil {
		return contract2.Error(fmt.Sprintf("合规事件%s存储失败，详细信息：%s",
			registration.Key(), err))
	}
	// 存储审计当事人规范对象所对应的存储条数
	if err = record2.StoreCount(registration, context); err != nil {
		return contract2.Error(fmt.Sprintf("合规事件%s相应的索引值存储失败，详细信息：%s",
			registration.Key(), err))
	}
	log.Println("审计事件录入成功, 审计事件ID:", registration.ID)
	return &contract2.Response{Payload: []byte("OK")}
}

// 根据审计当事人ID、项目ID以及规则ID，获取所有审计当事人的审计事件
func QueryEventsMain(args []string, context contract2.Context) *contract2.Response {
	if len(args) < 2 {
		return contract2.Error("查询失败，需要提查询供审计事件对应的当事人ID、项目ID以及规则ID")
	}

	// 根据传入的参数获取eventID
	eventID, err := project2.GetEventID(args, context)
	if err != nil {
		return contract2.Error(err.Error())
	}

	// 获取第几次录入信息
	index, err := record2.GetRecordCount(project2.GetRegistrationCountKey(*eventID), context)
	if err != nil {
		return contract2.Error(fmt.Sprintf("获取审计事件第几次录入信息出错，详细信息：%s", err.Error()))
	}

	// 获取开始及结束Key，查询满足条件的所有记录
	result, err := getQueryEventResult(eventID, index, context)
	if err != nil {
		return contract2.Error(err.Error())
	}
	return &contract2.Response{Payload: result}
}

// 获取审计事件查询的最终结果，json格式
func getQueryEventResult(eventID *common.Uint256,
	index uint32, context contract2.Context) ([]byte, error) {

	reg := project2.Registration{
		AuditeeSpecification: project2.AuditeeSpecification{ID: *eventID}}
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
func verify(registration *project2.Registration, context contract2.Context) error {
	if err := rules.ValidateRules(registration.Rule.ID, registration.Params,
		context); err != nil {
		return fmt.Errorf("合规事件%s规则验证失败，详细信息：%s", registration.ID, err)
	}

	return nil
}
