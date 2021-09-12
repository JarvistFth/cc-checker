package logger

import "testing"

func TestGetLogger(t *testing.T) {
	log := GetLogger()

	log.Debug("Debug!!")
	log.Warning("test - debug!!")
	log.Info("test - INFO!!")
	log.Error("error")
}
