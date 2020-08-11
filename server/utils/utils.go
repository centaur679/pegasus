package utils

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// Conf client配置
type Conf struct {
	Port int `yaml:"port"`
}

// GetConf 读取配置
func (c *Conf) GetConf(filePath string) *Conf {
	yamlFile, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Println(err.Error())
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		fmt.Println(err.Error())
	}
	return c
}

// Validate  校验配置
func (c *Conf) Validate() bool {
	switch {
	case c.Port <= 0:
		return false
	default:
		return true
	}
}
