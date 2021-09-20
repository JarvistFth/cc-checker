package logger

import (
	"testing"
)

var log = GetLogger()

func TestGetLogger(t *testing.T) {

	log.Info("hello")

	log.Debug("ok")

	log.Warning("123")

	log.Error("231")
}
