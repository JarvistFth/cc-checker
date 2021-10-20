package dummy

import (
	"testing"
)

func TestDummyTimeValidation_Validate(t *testing.T) {
	validation := LocationValidation{}

	assertError(t, func() error {
		// 不允许空参数
		return validation.Validate(0, []string{})
	})

	assertError(t, func() error {
		// 参数数量必须>=2
		return validation.Validate(0, []string{"39.9"})
	})

	assertError(t, func() error {
		// 参数不符合规范
		return validation.Validate(0, []string{"39.9", "xxx"})
	})

	assertError(t, func() error {
		// 经纬度（0,0）处
		return validation.Validate(0, []string{"0", "0"})
	})

	assertError(t, func() error {
		// 超出范围不远处
		return validation.Validate(0, []string{"39.91", "116.29"})
	})

	assertNoError(t, func() error {
		// 验证范围中心点位置
		return validation.Validate(0, []string{"39.9", "116.3"})
	})

	assertNoError(t, func() error {
		// 验证范围之内
		return validation.Validate(0, []string{"39.901", "116.299"})
	})
}

func assertError(t *testing.T, fn func() error) {
	if err := fn(); err == nil {
		t.Error("方法应该抛出异常")
		t.Fail()
	}
}

func assertNoError(t *testing.T, fn func() error) {
	if err := fn(); err != nil {
		t.Error("方法不应该抛出异常")
		t.Fail()
	}
}
