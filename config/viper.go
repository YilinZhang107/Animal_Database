/*
* @Author: Oatmeal107
* @Date:   2023/6/12 13:18
 */

package config

import (
	"fmt"
	"github.com/spf13/viper"
)

type Config struct {
	MysqlConfig MysqlConfig `yaml:"mysql" mapstructure:"mysql"`
}

type MysqlConfig struct {
	DbHost   string `yaml:"DbHost"`
	Port     string `yaml:"port"`
	Config   string `yaml:"config"`
	DbName   string `yaml:"DbName"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

var (
	Conf                    *Config //全局配置文件存储在这里
	ServerPort              string
	UploadTemplateExcelPath string

	// 地图相关
	Provinces   []string
	Prov2city   map[string][]string
	City2county map[string][]string
)

// InitViper  读取配置文件
func InitViper(p []string, p2c map[string][]string, c2c map[string][]string) {
	v := viper.New()
	v.SetConfigFile("./config/config.yaml")
	// 设置配置文件的名字
	//v.SetConfigName("config")
	//// 设置配置文件的类型
	//v.SetConfigType("yaml")
	//// 添加配置文件的路径，指定 config 目录下寻找
	//v.AddConfigPath("./config")

	// 寻找配置文件并读取
	err := v.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
	if err = v.Unmarshal(&Conf); err != nil {
		panic(err)
	}

	ServerPort = v.GetString("serverPort")
	UploadTemplateExcelPath = v.GetString("uploadTemplateExcelPath")

	// 地图相关
	Provinces, Prov2city, City2county = p, p2c, c2c
}
