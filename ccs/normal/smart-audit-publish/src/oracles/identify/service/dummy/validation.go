package dummy

import (
	"cc-checker/core/dynamic/ccs/normal/smart-audit-publish/src/core/contract"
	"cc-checker/core/dynamic/ccs/normal/smart-audit-publish/src/oracles/identify/service"
	"crypto/rand"
	"errors"
	"fmt"
	mrand "math/rand"
)

type EntityIdentifyValidation struct {
}

// 模拟物体识别规则验证
func (f *EntityIdentifyValidation) Validate(id contract.ServiceRuleID, args []string) error {
	if len(args) == 0 {
		return errors.New("不允许没有参数的验证")
	}

	// 与人脸识别类似，云从科技的物体识别服务支持base64编码的图片以及由图片生成的特征码。
	//	一个典型的特征码如下所示（注意特征码内容量较大，这里只列出了首尾部分，中间以...形式省略）：
	// g1T4vTuetL2grRVANLFpvwiGEz7YZkbAcWcwv/.../6ulAOO9aYQToAAIA/yZduwg==

	feature, err := f.getEntityFeature(args[0])
	if err != nil {
		return err
	}

	return f.entityCompare(args[0], feature)
}

// 调用物体识别预言机服务中的特征提取接口用以返回图片中相应物体的特征值。
// 这里我们返回一个随机值以模拟该特征值生成过程
func (f *EntityIdentifyValidation) getEntityFeature(entityRaw string) (rtn []byte, err error) {
	rtn = make([]byte, mrand.Int31n(1000))
	if _, err = rand.Read(rtn); err != nil {
		return
	}
	return
}

// 调用物体识别预言机服务中的物体比对接口，并返回包含EntityIdentifyResult中内容的比对结果。
// 这里我们返回一个物体评分在88~100之间的数值，以模拟在以90为界有>80%通过率的情况
func (f *EntityIdentifyValidation) getEntityCompareResult(entityId string,
	rtn []byte) service.EntityIdentifyResult {
	return service.EntityIdentifyResult{
		Result:   0,
		Score:    88 + uint32(mrand.Int31n(12)),
		Sequence: "XXXXX",
	}
}

// 这里对物体识别评分进行评价，评分超过或等于90则验证成功，否则视为其他物体
func (f *EntityIdentifyValidation) entityCompare(entityId string, rtn []byte) error {
	result := f.getEntityCompareResult(entityId, rtn)
	if result.Score < 90 {
		return errors.New(fmt.Sprintf("不属于类别码为%s类型的物体", entityId))
	}
	return nil
}
