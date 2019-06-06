package service

import (
	"encoding/json"
	"github.com/log-collection/sqlserver2kafka/models"
	"github.com/log-collection/sqlserver2kafka/pkg/gettime"
	"github.com/log-collection/sqlserver2kafka/pkg/logging"
	"go.uber.org/zap"
	"time"
)

type dmlToJson struct {
	DbType string `json:"dbtype"`
	DbHostname string `json:"dbhostname"`
	Date string `json:"date"`
	EventName string `json:"eventname"`
	ClientAppName string `json:"clientappname"`
	ClientHostName string `json:"clienthostname"`
	SQL string `json:"sql"`
	DatabaseName string `json:"databasename"`
	DbUsername string `json:"dbusername"`
	ObjectName string `json:"objectname"`
}

func GetDml2Kafka(aIDFlag chan int,threadPool chan bool,allFlagWrite chan Registry) {
	var dmlToJson dmlToJson
	tmpAIDFlag := <- aIDFlag
	tmpId := tmpAIDFlag + 100
	dmbAll,err := models.GetDml(tmpAIDFlag,tmpId)
	if err != nil {
		logging.Warn("查询dmb出错",zap.Error(err))
		return
	}
	tempNum := len(dmbAll)
	if tempNum == 0 {
		logging.Debug("查询dml表结果为空",zap.Int("AID起始位置",tmpAIDFlag,))
		logging.Debug("查询dml表结果为空",zap.Int("返回数据长度为",tempNum))
		time.Sleep(10*time.Second)
		aIDFlag <- tmpAIDFlag // 如果结果为空，则不增加ID值，并把值在传回管道给其他使用
		<- threadPool
		return
	}
	//logging.Debug("查询成功", zap.Int("返回数据长度为", tempNum))
	// 循环查询结果，发送到kafka
	for _, v := range dmbAll {
		dmlToJson.DbType = "sqlserver"
		dmlToJson.DbHostname = v.ServerName
		dmlToJson.Date = v.EventTime
		dmlToJson.EventName = v.EventName
		dmlToJson.ClientAppName = v.ClientAppName
		dmlToJson.ClientHostName = v.ClientHostName
		dmlToJson.SQL = v.BatchText
		dmlToJson.DatabaseName = v.DatabaseName
		dmlToJson.DbUsername = v.UserName
		dmlToJson.ObjectName = "dml"
		dmlJsonBytes,err := json.Marshal(dmlToJson)
		if err != nil {
			logging.Warn("解析dmlSQL数据失败",zap.Error(err))
		}
		logging.Debug(string(dmlJsonBytes))
		models.KafkaProducer(dmlJsonBytes)
		<- threadPool
	}
	tmpAIDFlag = tmpAIDFlag+tempNum // 加上传回数据的长度，为新的查询起点
	aIDFlag <- tmpAIDFlag  // 发送到管道，以便让其他goroutine使用
	dmlRegistry := Registry{
		Offset: tmpAIDFlag,
		Table: "dml",
		Timestamp: gettime.GetTimeStr(),
	}
	allFlagWrite <- dmlRegistry

}
