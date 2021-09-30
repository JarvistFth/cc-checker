package config

import (
	"cc-checker/logger"
	"cc-checker/utils"
	"golang.org/x/tools/go/ssa"
	"io/ioutil"
	"path/filepath"
	"runtime"
	"sigs.k8s.io/yaml"
)

var log = logger.GetLogger()

var SSconfig *Config

type Sink struct {
	Package  string `json:"package"`
	Method   string `json:"method"`
	Receiver string `json:"receiver"`
}

type Exclude struct {
	Package  string `json:"package"`
	Method   string `json:"method"`
	Receiver string `json:"receiver"`
}

type Source struct {
	Package  string `json:"package"`
	Method   string `json:"method"`
	Receiver string `json:"receiver"`
	Tag      string `json:"tag"`
}

type Config struct {
	Sources  []Source  `json:"sources"`
	Sinks    []Sink    `json:"sinks"`
	Excludes []Exclude `json:"excludes"`
}

func ReadConfig() (*Config, error) {
	_, f, _, _ := runtime.Caller(0)
	f = filepath.Join(filepath.Dir(f), "conf.yaml")
	bytes, err := ioutil.ReadFile(f)
	if err != nil {
		log.Fatalf("get sources and sink cfg files error: %s", err.Error())
		return nil, err
	}

	var config Config
	err = yaml.Unmarshal(bytes, &config)
	if err != nil {
		log.Fatalf("unmarshal json file error: %s", err.Error())
		return nil, err
	}

	SSconfig = &config
	return SSconfig, nil
}

func (c *Config) String() string {
	var ret string
	for _, s := range c.Sources {
		ret += "sources:" + s.Package + "." + s.Method + " receiver:" + s.Receiver + s.Tag + "\n"
	}

	for _, s := range c.Sinks {
		ret += "sinks:" + s.Package + "." + s.Method + " receiver:" + s.Receiver + "\n"
	}

	return ret
}

func (c *Config) IsExcluded(path, recv, name string) bool {
	for _, e := range c.Excludes {
		if e.Package == path && e.Receiver == recv && e.Method == name {
			return true
		}
	}
	return false
}

// IsSink determines whether a function is a sink.
func (c *Config) IsSink(path, recv, name string) bool {
	for _, s := range c.Sinks {
		if s.Package == path && s.Receiver == recv && s.Method == name {
			return true
		}
	}
	return false
}

func (c *Config) IsSource(path, recv, name string) (bool, string) {

	for _, s := range c.Sources {
		if s.Package == path && s.Receiver == recv && s.Method == name {
			return true, s.Tag
		}
	}
	return false, ""
}

func IsExcluded(path, recv, name string) bool {
	for _, e := range SSconfig.Excludes {
		if e.Package == path && e.Receiver == recv && e.Method == name {
			return true
		}
	}
	return false
}

// IsSink determines whether a function is a sink.
func IsSink(call ssa.CallInstruction) bool {
	callcom := call.Common()
	if fn := callcom.StaticCallee(); fn != nil {
		path, recv, name := utils.DecomposeFunction(fn)
		//log.Debugf("sink static callee: %s, %s, %s", path, recv, name)
		for _, s := range SSconfig.Sinks {
			if s.Package == path && s.Receiver == recv && s.Method == name {
				return true
			}
		}
	} else {
		//invoke mode
		if callcom.IsInvoke() {
			path, recv, name := utils.DecomposeAbstractMethod(callcom)
			//log.Debugf("sink invoke callee: %s, %s, %s", path, recv, name)
			for _, s := range SSconfig.Sinks {
				if s.Package == path && s.Receiver == recv && s.Method == name {
					return true
				}
			}
		} else {
			//dynamic call mode
		}
	}
	return false
}

func IsSource(call ssa.CallInstruction) (string, bool) {
	callcom := call.Common()
	if fn := callcom.StaticCallee(); fn != nil {
		//log.Debugf("source static fn:%s, name:%s", fn.String(), fn.Name())
		path, recv, name := utils.DecomposeFunction(fn)
		//log.Debugf("source static callee: %s, receiver:%s, name:%s", path, recv, name)
		for _, s := range SSconfig.Sources {
			if s.Package == path && s.Receiver == recv && s.Method == name {
				return s.Tag, true
			}
		}
	} else {
		//invoke mode
		if callcom.IsInvoke() {
			path, recv, name := utils.DecomposeAbstractMethod(callcom)
			//log.Debugf("source invoke callee: %s, receiver:%s, %s", path, recv, name)
			for _, s := range SSconfig.Sources {
				if s.Package == path && s.Receiver == recv && s.Method == name {
					return s.Tag, true
				}
			}
		} else {
			//dynamic call mode
		}
	}

	return "", false
}
