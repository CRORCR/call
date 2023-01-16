package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/jinzhu/configor"
)

type ConfigClass struct {
	Conf *ConfigV2
}

var (
	GlobalConfig = ConfigClass{}
)

func InitConfigV2() {
	GlobalConfig.LoadConfig(nil, nil)
	v, _ := json.Marshal(GlobalConfig)
	fmt.Println("读取---", string(v))
}

type ConfigDB struct {
	User ConfigMysql `yaml:"user"`
}

type ConfigMysql struct {
	Host          string `yaml:"host"`
	Username      string `yaml:"username"`
	Password      string `yaml:"password"`
	Database      string `yaml:"database"`
	Port          uint32 `yaml:"port"`
	IsAutoMigrate bool   `yaml:"is_auto_migrate"`
	LogMode       bool   `yaml:"log_mode"`
	MaxIdleConns  int    `yaml:"max_idle_conns"`
	MaxOpenConns  int    `yaml:"max_open_conns"`
}

type ConfigRedis struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
	DB   int    `yaml:"db"`
}

type ConfigApp struct {
	ENV  string `yaml:"env"`
	Port string `yaml:"port"`
}
type UrlParams struct {
	Idcade     string `yaml:"idcade"`
	Verify     string `yaml:"verify"`
	Type       string `yaml:"type"`
	ApiKey     string `yaml:"apiKey"`
	SecretKey  string `yaml:"secretKey"`
	ImageUrl   string `yaml:"imageUrl"`
	GatewayUrl string `yaml:"gatewayUrl"`
}

type ConfigParams struct {
	PreDeliveryHours float64 `yaml:"pre_delivery_hours"`
}

type ConfigV2 struct {
	DB        ConfigDB    `yaml:"mysql"`
	JwtPubKey string      `yaml:"jwtPubKey"`
	App       ConfigApp   `yaml:"app"`
	Redis     ConfigRedis `yaml:"redis"`
	Urls      UrlParams   `yaml:"urls"`
}

func (c *ConfigClass) LoadConfig(configFilePtr *string, secretFilePtr *string) {
	configPath := ``
	secretPath := ``
	if configFilePtr == nil {
		configPath = "./config/"
	} else {
		configPath = *configFilePtr
	}
	if secretFilePtr == nil {
		secretPath = "./secret/"
	} else {
		secretPath = *secretFilePtr
	}
	configfiles := GetConfigFiles(configPath, secretPath)
	c.Conf = new(ConfigV2)

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
