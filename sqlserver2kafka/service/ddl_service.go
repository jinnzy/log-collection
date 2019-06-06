package service

import (
	"encoding/json"
	"fmt"
	"github.com/log-collection/sqlserver2kafka/models"
	"github.com/log-collection/sqlserver2kafka/pkg/gettime"
	"github.com/log-collection/sqlserver2kafka/pkg/logging"
	"go.uber.org/zap"
	"time"
)

type ddlToJson struct {
	DbType string `json:"dbtype"`
	EventName string `json:"eventname"`
	Date string `json:"date"`
	DbHostname string `json:"dbhostname"`
	DbUsername string `json:"dbusername"`
	DatabaseName string `json:"databasename"`
	ObjectName string `json:"objectname"`
	SQL string `json:"sql"`
	ClientAppName string `json:"clientappname"`
	ClientHostName string `json:"clienthostname"`
	ClientIP string `json:"clientip"`
}

func GetDdl2Kafka(aIDFlag chan int,threadPool chan bool,allFlagWrite chan Registry) {
	var ddlToJson ddlToJson
	tmpAIDFlag := <- aIDFlag
	tmpId := tmpAIDFlag + 100
	dmbAll,err := models.GetDdl(tmpAIDFlag,tmpId)
	if err != nil {
		logging.Warn("查询ddl出错",zap.Error(err))
		return
	}
	tempNum := len(dmbAll)
	if tempNum == 0 {
		logging.Debug("查询ddl表结果为空",zap.Int("AID起始位置",tmpAIDFlag,))
		logging.Debug("查询ddl表结果为空",zap.Int("返回数据长度为",tempNum))
		time.Sleep(10*time.Second)
		aIDFlag <- tmpAIDFlag // 如果结果为空，则不增加ID值，并把值在传回管道给其他使用
		<- threadPool
		return
	}

	//logging.Debug("查询成功", zap.Int("返回数据长度为", tempNum))
	tmpAIDFlag = tmpAIDFlag+tempNum // 加上传回数据的长度，为新的查询起点
	aIDFlag <- tmpAIDFlag  // 发送到管道，以便让其他goroutine使用
	for _, v := range dmbAll {
		ddlToJson.DbType = "sqlserver"
		ddlToJson.EventName = v.EventType
		ddlToJson.Date = v.EventTime
		ddlToJson.ClientHostName = v.ClientHostName
		ddlToJson.DbHostname = v.ServerName
		ddlToJson.DbUsername = v.UserName
		ddlToJson.DatabaseName = v.DatabaseName
		ddlToJson.ObjectName = v.ObjectName
		ddlToJson.SQL = v.TSQLCommand
		ddlToJson.ClientAppName = v.ClientAppName
		ddlToJson.ClientHostName = v.ClientHostName
		ddlToJson.ClientIP = v.ClientIP
		ddlJsonBytes,err := json.Marshal(ddlToJson)
		if err != nil {
			logging.Warn("解析ddlSQL数据失败",zap.Error(err))
		}
		logging.Debug(string(ddlJsonBytes))
		models.KafkaProducer(ddlJsonBytes)
		// 读取数据，解除占用的并发数
		<- threadPool
		// 将flag相关信息传入管道，由RegisterWrite函数写入文件中
		ddlRegistry := Registry{
			Offset: tmpAIDFlag,
			Table: "ddl",
			Timestamp: gettime.GetTimeStr(),
		}
		logging.Debug("传入allFlagWrite值",zap.String("dllRegistry",fmt.Sprintln(ddlRegistry)))
		allFlagWrite <- ddlRegistry
	}

}
