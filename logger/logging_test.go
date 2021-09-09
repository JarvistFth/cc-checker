package logger

import "testing"

func TestGetLogger(t *testing.T) {
	log := GetLogger()


	log.Debug("test - debug!!")
	log.Info("test - INFO!!")
}
