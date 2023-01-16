package config

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"
)

type Config struct {
	App   *App   `yaml:"app"`
	Mysql *Mysql `yaml:"mysql"`
}

type App struct {
	Env  string `yaml:"env"`
	Port uint64 `yaml:"port"`
}

type Mysql struct {
	Name string `yaml:"name"`
}

func (a *App) ParamsStruct(b []byte) error {
	if err := json.Unmarshal(b, a); err != nil {
		return err
	}
	return nil
}

// 配置文件加载

// 每个服务的结构体定义在这里

func InitConfig() {
	InitConfigV2()
	return

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
	paramsConfig(*config)
}

func paramsConfig(configEnv string) {
	mysql := new(Mysql)
	app := new(App)
	var conf = &Config{}
	conf.App = app
	conf.Mysql = mysql

	envFolder := "./config/" + configEnv
	// 找到文件夹下面所有文件
	// 每一个文件读取，并反序列化成对应结构体
	// 检查文件夹是否存在
	if _, err := os.Stat(envFolder); os.IsNotExist(err) {
		panic(fmt.Sprintf("folder %s not exist,err:%v", envFolder, err))
	}

	// 读取每个文件
	files, err := ioutil.ReadDir(envFolder)
	if err != nil {
		panic("文件不存在")
	}

	for _, file := range files {
		if !strings.Contains(file.Name(), "app") {
			continue
		}
		fmt.Println("找到app文件", file.Name())

		loadVal, err := loadConfigFile(envFolder, file.Name())
		if err != nil {
			panic(err)
		}

		fileName := strings.Split(file.Name(), ".")
		if len(fileName) == 0 {
			continue
		}
		fmt.Println("--jhhh--", string(loadVal))
		switch fileName[0] {
		case "app":
			if err = json.Unmarshal(loadVal, app); err != nil {
				panic(err)
			}
		}
		fmt.Println("-------", *app)
	}
}

func loadConfigFile(folder string, file string) ([]byte, error) {
	//  判断文件是否以yaml或者yml作为后缀
	s := strings.Split(file, ".")
	if len(s) != 2 || s[1] != "yaml" {
		return nil, fmt.Errorf(" %s is illegal", file)
	}

	// 读取文件内容
	bf, err := ioutil.ReadFile(filepath.Join(folder, file))
	if err != nil {
		return nil, fmt.Errorf("read %s failed,err:%v", file, err)
	}
	fmt.Println("读物文件", string(bf))
	// 直接针对文本做环境变量的替换
	//bf = replace(bf, conf.envMaps)
	// 解析对应的文件
	c := make(map[string]interface{}, 0)
	if err = yaml.Unmarshal(bf, &c); err != nil {
		return nil, fmt.Errorf("read %s failed,err:%v", file, err)
	}

	conf, _ := json.Marshal(c)
	return conf, nil
}
