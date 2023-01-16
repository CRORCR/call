package config

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/CRORCR/call/internal/model"
	"github.com/jinzhu/configor"
)

// 配置文件根据模块不同，拆分多个文件存储，对不同性质的配置，根据环境不同读取不同配置，通过运行参数决定。
// 对于私密的配置，放到secret目录，普通配置放到config目录。

type Configuration struct {
	Conf *model.Config
}

func InitConfig() *Configuration {
	config := &Configuration{}
	envPackage := flag.String("config", "", "Configuration file")
	flag.Parse()
	if envPackage == nil {
		panic("Please enter Configuration file")
	}

	switch *envPackage {
	case EnvProduction, EnvTesting, EnvDevelopment:
	default:
		panic("Please enter Configuration file")
	}

	config.LoadConfig(envPackage, nil)

	v, _ := json.Marshal(config)
	fmt.Println("读取配置:", string(v))
	return config
}

// LoadConfig 配置必须要放到config和secret目录
func (c *Configuration) LoadConfig(configFilePtr *string, secretFilePtr *string) {
	configPath := "./config/"
	secretPath := "./secret/"
	if configFilePtr != nil { // 好像有一个拼接路径的包，替换掉 todo
		configPath = fmt.Sprintf("%s/%s/", configPath, *configFilePtr)
	}
	if secretFilePtr != nil {
		secretPath = *secretFilePtr
	}
	configfiles := GetConfigFiles(configPath, secretPath)
	c.Conf = new(model.Config)

	// 从配置文件中加载
	err := configor.Load(c.Conf, configfiles...)
	if err != nil {
		msg := "Failed to load config file !!! " + err.Error()
		panic(msg)
	}
}

func GetConfigFiles(dirs ...string) []string {
	configfiles := make([]string, 10)
	for i := 0; i < len(dirs); i++ {
		dir := dirs[i]
		configfiles = walkDir(configfiles, dir)
	}

	return deleteEmpty(configfiles)
}

func walkDir(configfiles []string, dirname string) []string {
	files, err := ioutil.ReadDir(dirname)
	if err == nil {
		for _, f := range files {
			if strings.Contains(f.Name(), ".yaml") {
				configfiles = append(configfiles, dirname+f.Name())
			}
		}
	}
	return configfiles
}

func deleteEmpty(configfiles []string) []string {
	var retConfigfiles []string
	for _, configfile := range configfiles {
		if configfile != "" {
			retConfigfiles = append(retConfigfiles, configfile)
		}
	}
	return retConfigfiles
}
