package service

import (
	"cc-checker/ccs/normal/smart-audit-publish/core/contract"
	"errors"
)

// 调用预言机服务的人脸比对接口返回的结果
type FaceCompareResult struct {
	// 人脸比对结果，0表示请求成功，非0会有响应的错误码意义
	Result uint32 `json:"result"`

	// 相似度评分，在0~100之间取值，越大表示越相似
	Score uint32 `json:"score"`

	// 会话序号
	Sequence string `json:"sequence"`
}

type FaceValidation struct {
}

// 人脸识别规则验证
func (f *FaceValidation) Validate(id contract.ServiceRuleID, args []string) error {
	if len(args) < 1 {
		return errors.New("验证规则所需参数不足")
	}

	return f.serviceValidate(id, args[0])
}

// 调用预言机验证人脸是否满足规则
func (f *FaceValidation) serviceValidate(id contract.ServiceRuleID,
	valueExpression string) error {
	// fixme 实际商用时实现人脸识别预言机，然后在这里调用预言机服务
	return nil
}
