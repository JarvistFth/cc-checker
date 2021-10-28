package main

import (
	utils2 "cc-checker/ccs/normal/studentmanage/utils"
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)

/**
 * @Author: WuNaiChi
 * @Date: 2020/6/4 14:13
 * @Desc:
 */
type STUChaincode struct{}

// todo:META-INF json文件的内容怎么写

// 创建学生信息
func (t *STUChaincode) createStudentsInfo(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()
	// 获取参数
	//logs.Info("StartCreate!!")
	res, err := utils2.GetCreateStudentParam(args[0])
	if err != nil {
		return utils2.Error(err.Error())
	}
	// 查询链上有没有
	_, exit, err := utils2.IfExist(stub, res.AcctId)
	if err != nil {
		return utils2.Error(err.Error())
	}
	// 存在就不创建
	if exit {
		return utils2.Error("don't have the student info")
	}
	// 没有则创建
	chainInStudentInfo := utils2.ChainOfStudentInfo{
		DocType:     utils2.StudentDocType,
		StudentInfo: res,
	}
	// 信息写到链上
	err = utils2.WriteInfoToChain(stub, res.AcctId, chainInStudentInfo)
	if err != nil {
		//logs.Error("failed to create")
		return utils2.Error(err.Error())
	}
	return utils2.SUCCESS(nil)
}

// 更新学生信息
func (t *STUChaincode) updateStudentInfo(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()
	//logs.Info("StartUpdate!!")
	rsp, err := utils2.GetUpdateStudentParam(args[0])
	if err != nil {
		return utils2.Error(err.Error())
	}
	// 查询链上有没有
	stateByte, exit, err := utils2.IfExist(stub, rsp.AcctId)
	if err != nil {
		return utils2.Error(err.Error())
	}
	if !exit {
		return utils2.Error("don't have the student info")
	}
	studentInfo := new(utils2.ChainOfStudentInfo)
	err = json.Unmarshal(stateByte, studentInfo)
	if err != nil {
		return utils2.Error(err.Error())
	}
	studentInfo.StudentInfo.Name = rsp.Name
	studentInfo.StudentInfo.Grade = rsp.Grade
	studentInfo.StudentInfo.Hobby = rsp.Hobby

	err = utils2.WriteInfoToChain(stub, rsp.AcctId, studentInfo)
	if err != nil {
		//logs.Info("FailedStartUpdate!!")
		return utils2.Error(err.Error())
	}
	return utils2.SUCCESS([]byte{})
}

// 查询学生信息列表
func (t *STUChaincode) queryStudentList(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()
	// 获取参数
	//logs.Info("StartQuery!!")
	fmt.Println("StartQuery!!!")
	rsp, err := utils2.GetQueryStudentParam(args[0])
	if err != nil {
		fmt.Println("GetQueryStudentParam!!!", err)
		return utils2.Error(err.Error())
	}
	// 参数校验sdk来做
	// 配置select语句
	fmt.Println("GetQueryStudentParam!!!")
	selectString, err := utils2.GetQueryStuListSelectString(rsp, utils2.StudentDocType)
	if err != nil {
		fmt.Println("GetQueryStuListSelectString!!!", err)
		return utils2.Error(err.Error())
	}
	fmt.Println("GetQueryStuListSelectString!!!", selectString)
	StateQueryIterator, QueryResponseMetadata, err := utils2.GetCounchdbIter(rsp.Bookmark, selectString, rsp.PageNo, rsp.PageSize, stub)
	if err != nil {
		fmt.Println("GetCounchdbIter!!!", err)
		return utils2.Error(err.Error())
	}
	fmt.Println("GetCounchdbIter!!!")
	defer StateQueryIterator.Close() // todo:这个是为了什么
	// todo : 这里为什么是成功
	if QueryResponseMetadata == nil {
		fmt.Println("QueryResponseMetadata is nil!!!")
		return utils2.SUCCESS([]byte{})
	}
	fmt.Println("QueryResponseMetadata!!!")
	var res = utils2.RspQueryStudentInfo{
		Bookmark: QueryResponseMetadata.Bookmark,
		Count:    QueryResponseMetadata.FetchedRecordsCount,
	}
	// 获取迭代器的指针
	for StateQueryIterator.HasNext() {
		next, err := StateQueryIterator.Next()
		if err != nil {
			fmt.Println("StateQueryIterator!!!", err)
			return utils2.Error(err.Error())
		}
		fmt.Println("StateQueryIterator!!!")
		studentInfo := utils2.ChainOfStudentInfo{}
		err = json.Unmarshal(next.Value, &studentInfo)
		if err != nil {
			fmt.Println("Unmarshal!!!", err)
			return utils2.Error(err.Error())
		}
		res.StudentInfo = append(res.StudentInfo, studentInfo.StudentInfo)
	}
	fmt.Println("Unmarshal!!!")
	if len(res.StudentInfo) == 0 {
		return utils2.SUCCESS([]byte{})
	}

	repStudentInfo, err := json.Marshal(res)
	if err != nil {
		//logs.Info("FailedQuery!!")
		fmt.Println("Marshal!!!", err)
		return utils2.Error(err.Error())
	}
	fmt.Println("Marshal!!!")
	return utils2.SUCCESS(repStudentInfo)
}

// 查询学生详情
func (t *STUChaincode) queryStudentInfo(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()
	// 获取参数
	//logs.Info("StartQueryInfo!!")
	fmt.Println("StartQueryInfo!!!")
	rsp, err := utils2.GetQueryStudentInfoParam(args[0])
	if err != nil {
		fmt.Println("GetQueryStudentInfoParam!!!", err)
		return utils2.Error(err.Error())
	}
	// 参数校验sdk来做
	// 使用k-v查询
	stateByte, _, err := utils2.IfExist(stub, rsp.AcctId)
	if err != nil {
		fmt.Println("IfExist!!!", err)
		return utils2.Error(err.Error())
	}

	if stateByte == nil {
		return utils2.SUCCESS([]byte{})
	}

	var studentInfo = utils2.ChainOfStudentInfo{}
	err = json.Unmarshal(stateByte, &studentInfo)
	if err != nil {
		fmt.Println("Unmarshal!!!", err)
		return utils2.Error(err.Error())
	}

	tmp := utils2.StudentInfo{
		AcctId: studentInfo.StudentInfo.AcctId,
		Name:   studentInfo.StudentInfo.Name,
		Sex:    studentInfo.StudentInfo.Sex,
		Grade:  studentInfo.StudentInfo.Grade,
		Hobby:  studentInfo.StudentInfo.Hobby,
	}

	repStudentInfo, err := json.Marshal(tmp)
	if err != nil {
		//logs.Info("FailedQueryInfo!!")
		fmt.Println("Marshal!!!", err)
		return utils2.Error(err.Error())
	}
	return utils2.SUCCESS(repStudentInfo)
}

// 删除学生信息
func (t *STUChaincode) deleteStudentInfo(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()
	// 获取参数
	//logs.Info("StartDeleteInfo!!")
	rsp, err := utils2.GetDeleteStudentInfoParam(args[0])
	if err != nil {
		return utils2.Error(err.Error())
	}
	stateByte, exist, err := utils2.IfExist(stub, rsp.AcctId)
	if err != nil || !exist {
		return utils2.Error(err.Error())
	}
	if stateByte == nil {
		return utils2.SUCCESS([]byte{})
	}
	err = stub.DelState(rsp.AcctId)
	if err != nil {
		//logs.Info("FailedDeleteInfo!!")
		return utils2.Error(err.Error())
	}
	return utils2.SUCCESS([]byte{})
}

func (t *STUChaincode) Init(stub shim.ChaincodeStubInterface) peer.Response {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()
	role := &utils2.AdminRole{
		RoleId:    utils2.RoleAdmin,
		Name:      "管理员",
		Authority: utils2.AdminAuthority,
	}
	err := utils2.WriteInfoToChain(stub, role.RoleId, role)
	if err != nil {
		return utils2.Error(err.Error())
	}

	return utils2.SUCCESS(nil)
}
func (t *STUChaincode) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	//panic恢复
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()

	fmt.Println("Invoke!!!!!")
	//检查调用的函数是否存在
	function, args := stub.GetFunctionAndParameters()

	fmt.Println("func, args: ", function, args)
	//校验入参长度
	if len(args) != 1 {
		fmt.Println("func, args: ", args, len(args))
		return utils2.Error(utils2.InputParaError)
	}

	switch function {
	//角色
	case utils2.AuthorityFuncMap[utils2.AuthorityCreateMember]: // 上链学生信息（Admin）
		return t.createStudentsInfo(stub, args)
	case utils2.AuthorityFuncMap[utils2.AuthorityUpdateMember]: // 更新学生信息（Admin）
		return t.updateStudentInfo(stub, args)
	case utils2.AuthorityFuncMap[utils2.AuthorityGetMemberList]: // 查新学生信息（Admin）
		return t.queryStudentList(stub, args)
	case utils2.AuthorityFuncMap[utils2.AuthorityGetMemberInfo]: // 查询学生信息列表（Admin）
		return t.queryStudentInfo(stub, args)
	case utils2.AuthorityFuncMap[utils2.AuthorityDeleteMember]: // 删除学生信息（Admin）
		return t.deleteStudentInfo(stub, args)

	default:
		return utils2.Error(utils2.InternalError)
	}

}

func main() {
	err := shim.Start(new(STUChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}
