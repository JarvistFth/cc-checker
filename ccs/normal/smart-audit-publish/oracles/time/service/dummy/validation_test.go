package dummy

import (
	"testing"
	"time"
)

func TestDummyTimeValidation_Validate(t *testing.T) {
	validation := TimeValidation{}

	assertError(t, func() error {
		// 不允许空参数
		return validation.Validate(0, []string{})
	})

	assertError(t, func() error {
		// 参数不符合规范
		return validation.Validate(0, []string{"2020-05-12"})
	})

	assertError(t, func() error {
		// 早于9点
		return validation.Validate(0, []string{
			timeFromHour(8).Format(layout),
		})
	})

	assertError(t, func() error {
		// 晚于下午6点
		return validation.Validate(0, []string{
			timeFromHour(19).Format(layout),
		})
	})

	// 正常情况
	assertNoError(t, func() error {
		return validation.Validate(0, []string{
			timeFromHour(12).Format(layout),
		})
	})
}

func timeFromHour(hour int) time.Time {
	now := time.Now()
	return time.Date(now.Year(), now.Month(), now.Day(),
		hour, 0, 0, 0, time.UTC)
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
