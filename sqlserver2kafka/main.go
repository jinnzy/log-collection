package main

import (
	"fmt"
	//"fmt"
	"github.com/log-collection/sqlserver2kafka/models"
	"github.com/log-collection/sqlserver2kafka/pkg/logging"
	"github.com/log-collection/sqlserver2kafka/pkg/setting"
	"github.com/log-collection/sqlserver2kafka/service"
	"go.uber.org/zap"
	"os"
)

func main() {
	logging.Init()
	setting.Init()
	models.Init()
	models.KafkaInit()

	allFlagWrite := make(chan service.Registry, 1000)

	dmlAIDFlag := make(chan int, 1)
	ddlAIDFlag := make(chan int, 1)
	loginFailedFlag := make(chan int, 1)
	// 读取registry文件，
	_, err := os.Stat("./registyr")
	if err != nil {
		logging.Warn("读取registry文件不存在，将offset值设置为1")
		//  dml flag
		dmlAIDFlag <- 1
		// ddl flag
		ddlAIDFlag <- 1
		// login failed flag
		loginFailedFlag <- 1
	}
	if err == nil {
		// 返回数组 0 位置是dml 1位置是ddl 2位置是loginfailed
		tempMap := map[string]string{
			"ddl":         "",
			"dml":         "",
			"loginFailed": "",
		}
		flagData := service.RegistryRead()
		// 判断flagData的值，因为registry里面的值可能只有1个或2个，
		for _, v := range flagData {
			switch v.Table {
			case "ddl":
				if v.Offset != 0 {
					ddlAIDFlag <- v.Offset
					delete(tempMap, "ddl")
					logging.Debug("删除tempmap ddl", zap.String("tempMap", fmt.Sprintln(tempMap)))
				}
			case "dml":
				if v.Offset != 0 {
					dmlAIDFlag <- v.Offset
					delete(tempMap, "dml")
				}
			case "loginFailed":
				if v.Offset != 0 {
					loginFailedFlag <- v.Offset
					delete(tempMap, "loginFailed")
				}
			default:
				logging.Debug("switch flagData默认分支略过")
			}
		}
		if len(tempMap) > 0 {
			for kMap := range tempMap {
				switch kMap {
				case "ddl":
					logging.Debug("进入len(tempmap)",zap.String("kmap",kMap))
					ddlAIDFlag <- 1
				case "dml":
					logging.Debug("进入len(tempmap)",zap.String("kmap",kMap))
					dmlAIDFlag <- 1
				case "loginFailed":
					logging.Debug("进入len(tempmap)",zap.String("kmap",kMap))
					loginFailedFlag <- 1
				}
			}
		}
	}

	//for i := 0; i < 20;i ++ {
	threadPool := make(chan bool, 30)
	go service.RegisterWrite(allFlagWrite)
	for {
		threadPool <- true
		go service.GetDml2Kafka(dmlAIDFlag, threadPool, allFlagWrite)
		threadPool <- true
		go service.GetDdl2Kafka(ddlAIDFlag, threadPool, allFlagWrite)
		threadPool <- true
		go service.GetLoginFailed2Kafka(loginFailedFlag, threadPool, allFlagWrite)
	}

}
