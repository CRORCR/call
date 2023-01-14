package contract

import (
	"io/ioutil"
	"os"
	"strings"

	"github.com/jinzhu/configor"
)

type ConfigClass struct {
	Conf *Config
}

var (
	GlobalConfig = ConfigClass{}
)

func InitConfig(configFilePtr *string, secretFilePtr *string) {
	GlobalConfig.LoadConfig(configFilePtr, secretFilePtr)
}

type ConfigDB struct {
	User   ConfigMysql `yaml:"user"`
	Engine ConfigMysql `yaml:"engine"`
	VIP    ConfigMysql `yaml:"vip"`
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
	ConfigCurrency string `yaml:"config"`
	Balance        string `yaml:"balance"`
	Group          string `yaml:"group"`
	PathLocal      string `yaml:"pathLocal"`
	PathDebug      string `yaml:"pathDebug"`
	PathProd       string `yaml:"pathProd"`
}
type ConfigParams struct {
	PreDeliveryHours float64 `yaml:"pre_delivery_hours"`
}

type Config struct {
	DB        ConfigDB    `yaml:"mysql"`
	JwtPubKey string      `yaml:"jwtPubKey"`
	App       ConfigApp   `yaml:"app"`
	Redis     ConfigRedis `yaml:"redis"`
	Urls      UrlParams   `yaml:"urls"`
}

func init() {
	InitConfig(nil, nil)
}

func (this *ConfigClass) LoadConfig(configFilePtr *string, secretFilePtr *string) {
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
	this.Conf = new(Config)
	// 判断配置文件的来源
	if os.Getenv("CONFIG_ETCD") != "" {
		return
	}

	// 从配置文件中加载
	err := configor.Load(this.Conf, configfiles...)
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
