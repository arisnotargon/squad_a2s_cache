package config

import (
	"fmt"
	"os"
	"sync"

	"gopkg.in/yaml.v3"
)

type AppConfig struct {
	DevName string `yaml:"dev_name"`
}

var (
	Conf *AppConfig
	once sync.Once
)

func init() {
	once.Do(initConf)
}

func initConf() {
	if Conf == nil {
		dataBytes, err := os.ReadFile("config.yml")
		if err != nil {
			fmt.Printf("读取config.yml文件失败:[%+v]\n", err)
			return
		}
		Conf = &AppConfig{}
		err = yaml.Unmarshal(dataBytes, Conf)
		if err != nil {
			fmt.Printf("Unmarshal config.yml文件失败:[%+v], %s\n", err, string(dataBytes))
			return
		}
	}
}
