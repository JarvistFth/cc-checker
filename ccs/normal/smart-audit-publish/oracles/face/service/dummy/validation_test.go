package dummy

import (
	"crypto/rand"
	"log"
	"testing"
)

func TestFaceValidation_Validate(t *testing.T) {
	validation := FaceValidation{}

	success := 0
	fail := 0
	for i := 0; i < 10; i++ {
		if err := validation.Validate(0, []string{base64Img(),
			"xxxx"}); err != nil {
			fail++
		} else {
			success++
		}
	}

	if success == 0 || fail == 0 {
		t.Error("所有验证都成功或失败")
		t.Fail()
	}
	log.Printf("成功次数：%d, 失败次数: %d", success, fail)
}

func base64Img() string {
	buf := make([]byte, 20)
	rand.Read(buf)
	return string(buf)
}
