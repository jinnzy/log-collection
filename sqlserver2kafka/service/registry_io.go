package service

import (
	"encoding/json"
	"fmt"
	"github.com/log-collection/sqlserver2kafka/pkg/logging"
	"github.com/log-collection/sqlserver2kafka/pkg/setting"
	"go.uber.org/zap"
	"io/ioutil"
	"log"
	"os"
	"time"
)

type Registry struct {
	Offset int
	Table string
	Timestamp string
}
// [{"table":"ddl","offset":668427496,"timestamp":"2019-04-11T22:01:01.645943782+08:00"},{"table":"dml","offset":19916300,"timestamp":"2019-04-11T22:00:21.031866487+08:00"}] 格式
func RegistryRead() (flag []Registry) {
	var err error
	buf, err := ioutil.ReadFile(setting.Config.Registry)
	if err != nil {
		log.Println("error:",err)
		logging.Error("读取registry文件失败",zap.Error(err))
		return
	}
	// 定义数组内容类型是map k为string v 为interface{}
	flag = []Registry{} // 赋空值
	// 反序列化，解析成map类型，以便修改
	err = json.Unmarshal(buf,&flag)
	if err != nil {
		logging.Error("解析registry失败",zap.Error(err))
		return
	}
	logging.Debug("读取registry",zap.String("flag",fmt.Sprintln(flag)))
	return flag
	// 循环数组看里面有几组map，并修改对应的值
	//for index,v := range flag {
	//	v["offset"] = 2
	//	flag[index] = v
	//}
}
func RegisterWrite(allFlagWrite chan Registry )  {
	flag := [3]Registry{}
	for {
		select {
		case v := <- allFlagWrite:
			logging.Debug("接受allFlagWrite传来的值",zap.String("值为",fmt.Sprintln(v)))
			// 将传来的结构体值进行序列化
			if v.Table == "ddl" {
				flag[0] = v
			}
			if v.Table == "dml" {
				flag[1] = v
			}
			if v.Table == "loginFailed" {
				flag[2] = v
			}
			logging.Debug("flag",zap.String("dllRegistry",fmt.Sprintln(flag)))
			jsonStrTemp,err := json.Marshal(flag)
			if err != nil {
				logging.Error("序列化allFlag值失败",zap.Error(err))
				return
			}
			logging.Debug("写入文件",zap.String("info",string(jsonStrTemp)))
			// 打开文件
			file,err := os.OpenFile(setting.Config.Registry, os.O_WRONLY|os.O_TRUNC|os.O_CREATE|os.O_SYNC,0755)
			if err != nil {
				logging.Error("打开或创建registry文件失败",zap.Error(err))
				return
			}
			_,err = file.WriteString(string(jsonStrTemp))
			if err != nil {
				logging.Error("写入文件失败",zap.Error(err))
				return
			}
			err = file.Sync()
			if err != nil {
				logging.Error("实时写入文件错误",zap.Error(err))
			}
			defer  file.Close()

		default:
			logging.Debug("等待allFlagWrite channel传来数据")
			time.Sleep(2*time.Second)
		}
	}
}
//func WriteFile(jsonStrTemp []byte)  {
//
//
//}
