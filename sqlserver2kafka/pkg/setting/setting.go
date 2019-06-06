package setting

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

// 配置文件管理
type AppYaml struct {
	RunMode string `yaml:"runMode"`
	Database database `yaml:"database"`
	Kafka kafka `yaml:"kafka"`
	Registry string `yaml:"registry"`
}

type database struct {
	Type string `yaml:"type"`
	User string `yaml:"user"`
	Password string `yaml:"password"`
	Host string `yaml:"host"`
	Name string `yaml:"name"`
}
type kafka struct {
	Host string `yaml:"host"`
	Topic string `yaml:"topic"`
}

var Config AppYaml
func Init()  {
	// 读取配置文件
	data,err := ioutil.ReadFile("sqlserver2kafka/conf/app.yaml")
	if err != nil {
		log.Panic("read conf/app.yaml err:",err)
	}
	Config = AppYaml{}
	// yaml解析
	err = yaml.Unmarshal(data,&Config)
	if err != nil {
		log.Panic("unmarshal conf/app.yaml err:",err)
	}
}