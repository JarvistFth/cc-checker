package invokes

import (
	"cc-checker/ccs/normal/smart-audit-publish/src/core/contract"
	"cc-checker/ccs/normal/smart-audit-publish/src/core/orgnization"
	"cc-checker/ccs/normal/smart-audit-publish/src/core/record"
	"fmt"
	"log"
	"strconv"
	"strings"
)

const (
	maintainerSplitter = ","
)

// 初始化审计合约
func InitMain(context contract.Context) *contract.Response {
	maintainers, err := parseMaintainers(context.GetArgs(), context)
	if err != nil {
		contract.Error(fmt.Sprint("解析合约运维人员出错，详细信息："))
	}
	// 存储审计运维人员信息
	for _, v := range maintainers {
		if err := record.StoreItem(v, context); err != nil {
			return contract.Error(fmt.Sprintf("运维人员'%s'信息存储出错，详细信息：%s",
				v.Name, err.Error()))
		}
	}
	// 存储审计运维人员个数
	err = context.PutState(orgnization.MaintainerCountKey, []byte(strconv.FormatUint(uint64(
		len(maintainers)), 10)))
	if err != nil {
		return contract.Error(fmt.Sprintf("运维人员个数信息存储出错，详细信息：%s",
			err.Error()))
	}

	log.Print("审计运维人员信息初始化成功")
	return &contract.Response{}
}

// 从输入参数中解析包括哪些审计维护人员
func parseMaintainers(args []string, context contract.Context) (
	[]*orgnization.Maintainer, error) {
	var result []*orgnization.Maintainer
	for i, v := range args {
		m, err := parseSingleMaintainer(v, uint32(i), context)
		if err != nil {
			return nil, err
		}

		result = append(result, m)
	}
	return result, nil
}

// 从输入参数中解析单个审计维护人员
func parseSingleMaintainer(arg string, index uint32, context contract.Context) (
	*orgnization.Maintainer, error) {
	subArgs := strings.Split(arg, maintainerSplitter)
	return orgnization.MaintainerFromString(subArgs, index, context)
}
