package logger

import (
	"encoding/json"
	"github.com/op/go-logging"
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
	//LogFile,_ = os.OpenFile("test-" + time.Now().Format("2006-01-02 15:04:05"+".log"),os.O_APPEND|os.O_WRONLY|os.O_CREATE,0666)
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
		backend.SetLevel(mapLevel(output.Level),"")
		backends = append(backends,backend)
		//logging.SetBackend(backend)
	}

	logging.SetBackend(backends...)
}

//CRITICAL Level = iota
//	ERROR
//	WARNING
//	NOTICE
//	INFO
//	DEBUG
func mapLevel(level string) logging.Level {
	switch level {
	case "DEBUG":
		return logging.DEBUG
	case "INFO":
		return logging.INFO
	case "WARNING":
		return logging.WARNING
	case "NOTICE":
		return logging.NOTICE
	case "ERROR":
		return logging.ERROR
	default:
		return logging.ERROR
	}
}

func SetLogger(logName string){
	//now := time.Now().Format("2006-01-02 15:04:05")
	//var err error
	//LogFile,err = os.OpenFile(logName+ "-" + now + "-info.log",os.O_APPEND|os.O_WRONLY|os.O_CREATE,0666)
	//debugLogfile,err = os.OpenFile(logName+ "-" + now + "-debug.log",os.O_APPEND|os.O_WRONLY|os.O_CREATE,0666)
	//if err != nil{
	//	log.Fatalf(err.Error())
	//}
	//var debugBackend *logging.LogBackend
	//var debugformat logging.Formatter
	//if debugstdout{
	//	debugBackend = logging.NewLogBackend(os.Stdout,"",0)
	//}else{
	//	debugBackend = logging.NewLogBackend(debugLogfile,"",0)
	//}
	//InfoBackend := logging.NewLogBackend(LogFile,"",0)

	//if debugstdout{
	//	debugformat = logging.MustStringFormatter(`%{color:reset}[%{level:.5s}] %{time:15:04:05} %{shortfile}  %{message}`)
	//}else{
	//	debugformat = logging.MustStringFormatter(`[%{level:.5s}] %{time:15:04:05} %{shortfile} %{callpath:1}: %{message}`)
	}
//	infoformat := logging.MustStringFormatter(`[%{level:.5s}] %{time:15:04:05}  %{shortfile} %{callpath:1}: %{message}`)
//
//	debugbandf := logging.NewBackendFormatter(debugBackend,debugformat)
//	infobandf := logging.NewBackendFormatter(InfoBackend,infoformat)
//
//	backend1level := logging.AddModuleLevel(debugbandf)
//	backend1level.SetLevel(logging.DEBUG,"")
//	backend2level := logging.AddModuleLevel(infobandf)
//	backend2level.SetLevel(logging.INFO,"")
//	logging.SetBackend(backend1level,backend2level)
//	Logger = logging.MustGetLogger("main")
//}
//
//func SetLoggerSTD(){
//	var debugBackend *logging.LogBackend
//	var debugformat logging.Formatter
//	debugBackend = logging.NewLogBackend(os.Stdout,"",0)
//	debugformat = logging.MustStringFormatter(`%{color:reset}[%{level:.5s}] %{time:15:04:05} %{shortfile}  %{message}`)
//	//infoformat := logging.MustStringFormatter(`[%{level:.5s}] %{time:15:04:05}  %{shortfile} %{callpath:1}: %{message}`)
//
//	debugbandf := logging.NewBackendFormatter(debugBackend,debugformat)
//	//infobandf := logging.NewBackendFormatter(InfoBackend,infoformat)
//
//	backend1level := logging.AddModuleLevel(debugbandf)
//	backend1level.SetLevel(logging.DEBUG,"")
//	logging.SetBackend(backend1level)
//	Logger = logging.MustGetLogger("main")
//}
//
//func GetLogger() *logging.Logger{
//	if Logger == nil{
//		SetLogger("chaincode-checker-log")
//	}
//	return Logger
//}
//
//func GetLoggerWithFileName(logName string) *logging.Logger{
//	if Logger == nil{
//		SetLogger(logName)
//	}
//	return Logger
//}
//
func GetLogger() *logging.Logger{
	if Logger == nil{
		cfg := readConfig()
		setupLogger(cfg)
	}
	return Logger
}
//
//
