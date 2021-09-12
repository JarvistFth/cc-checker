package logger

import (
	"encoding/json"
	"github.com/JarvistFth/go-logging"
	"io/ioutil"
	"os"
	"time"
)

var Logger *logging.Logger
var LogFile *os.File
var debugLogfile *os.File
var debugstdout = false

type LoggerConfig struct {
	Outputs []struct {
		Output  string `json:"output"`
		Level   string `json:"level"`
		Pattern string `json:"pattern"`
		Name    string `json:"name,omitempty"`
	} `json:"outputs"`
}


//
func init() {
	cfg := readConfig()
	setupLogger(cfg)
}

func readConfig() LoggerConfig {
	content,err := ioutil.ReadFile("./conf.json")
	if err != nil{
		panic("read logger json config error!")
	}
	var cfg LoggerConfig
	err = json.Unmarshal(content,&cfg)
	if err != nil{
		panic("unmarshal json config error")
	}
	return cfg
}

func setupLogger(cfg LoggerConfig){
	Logger = logging.MustGetLogger("main")
	var backends []logging.Backend
	for _,output := range cfg.Outputs{
		var logBackends *logging.LogBackend
		var format logging.Formatter
		switch output.Output {
		case "console":
			logBackends = logging.NewLogBackend(os.Stdout,"",0)
			format = logging.MustStringFormatter(output.Pattern)
		case "file":
			LogFile,_ := os.OpenFile(output.Name+"-"+time.Now().Format("2006-01-02-15:04:05")+".log",os.O_APPEND|os.O_WRONLY|os.O_CREATE,0666)
			logBackends = logging.NewLogBackend(LogFile,"",0)
			format = logging.MustStringFormatter(output.Pattern)
		}

		backend := logging.AddModuleLevel(logging.NewBackendFormatter(logBackends,format))
		backend.SetLevel(mapLevel(output.Level),"INFO")
		backends = append(backends,backend)
		//logging.SetBackend(backend)
	}
	logging.SetBackend(backends...)
}

func mapLevel(level string) logging.Level {
	switch level {
	case "DEBUG":
		return logging.DEBUG
	case "INFO":
		return logging.NOTICE
	case "WARN":
		return logging.WARN
	case "NOTICE":
		return logging.NOTICE
	case "ERROR":
		return logging.ERROR
	default:
		return logging.ERROR
	}
}

func GetLogger() *logging.Logger{
	if Logger == nil{
		cfg := readConfig()
		setupLogger(cfg)
	}
	return Logger
}
//
//
