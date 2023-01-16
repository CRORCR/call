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

type Configuration struct {
	Conf *model.Config
}

var (
	GlobalConfig = Configuration{}
)

func InitConfig() {
	config := flag.String("config", "", "Configuration file")
	flag.Parse()
	if config == nil {
		panic("Please enter Configuration file")
	}
	switch *config {
	case EnvProduction, EnvTesting, EnvDevelopment:
	default:
		panic("Please enter Configuration file")
	}
	fmt.Println("--------------")
	fmt.Println("./config/" + *config)
	fmt.Println("./config/")
	//GlobalConfig.LoadConfig("./config/"+*config+"/", "")
	//GlobalConfig.LoadConfig("./config/", "")
	GlobalConfig.LoadConfig("./config/testing/", "")
	v, _ := json.Marshal(GlobalConfig)
	fmt.Println("读取---", string(v))
}

func (c *Configuration) LoadConfig(configFilePtr string, secretFilePtr string) {
	configPath := ``
	secretPath := ``
	if len(configFilePtr) == 0 {
		configPath = "./config/"
	} else {
		configPath = configFilePtr
	}
	if len(secretFilePtr) == 0 {
		secretPath = "./secret/"
	} else {
		secretPath = secretFilePtr
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
				fmt.Println("加入的文件----", dirname+f.Name())
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
