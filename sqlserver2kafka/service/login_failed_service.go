package service

import (
	"encoding/json"
	"github.com/log-collection/sqlserver2kafka/models"
	"github.com/log-collection/sqlserver2kafka/pkg/gettime"
	"github.com/log-collection/sqlserver2kafka/pkg/logging"
	"go.uber.org/zap"
	"time"
)

type loginFailedToJson struct {
	DbType string `json:"dbtype"`
	Date string `json:"date"`
	ProcessInfo string `json:"processInfo"`
	Text string `json:"text"`
	ObjectName string `json:"objectname"`
	DbHostname string `json:"dbhostname"`
}

func GetLoginFailed2Kafka(aIDFlag chan int,threadPool chan bool,allFlagWrite chan Registry) {
	var loginFailedToJson loginFailedToJson
	tmpAIDFlag := <- aIDFlag
	tmpId := tmpAIDFlag + 100
	loginFailed,err := models.GetLoginFailed(tmpAIDFlag,tmpId)
	if err != nil {
		logging.Warn("查询loginFailed出错",zap.Error(err))
		return
	}
	tempNum := len(loginFailed)
	if tempNum == 0 {
		logging.Debug("查询loginFailed表结果为空",zap.Int("logdate起始位置",tmpAIDFlag,))
		logging.Debug("查询loginFailed表结果为空",zap.Int("返回数据长度为",tempNum))
		time.Sleep(10*time.Second)
		aIDFlag <- tmpAIDFlag // 如果结果为空，则不增加ID值，并把值在传回管道给其他使用
		<- threadPool
		return
	}

	//logging.Debug("查询成功", zap.Int("返回数据长度为", tempNum))
	aIDFlag <- tmpAIDFlag + tempNum // 加上传回数据的长度，为新的查询起点
	for _, v := range loginFailed {
		loginFailedToJson.DbType = "sqlserver"
		loginFailedToJson.Date = v.Date
		loginFailedToJson.ProcessInfo = v.ProcessInfo
		loginFailedToJson.Text = v.Text
		loginFailedToJson.DbHostname = v.DbHostname
		loginFailedToJson.ObjectName = "loginfailed"
		loginFailedJsonBytes,err := json.Marshal(loginFailedToJson)
		if err != nil {
			logging.Warn("解析loginFailedSQL数据失败",zap.Error(err))
			return
		}
		logging.Debug(string(loginFailedJsonBytes))
		//models.KafkaProducer(loginFailedJsonBytes)
		<- threadPool
	}
	// 将flag相关信息传入管道，由RegisterWrite函数写入文件中
	loginFailedRegistry := Registry{
		Offset: tmpAIDFlag,
		Table: "loginFailed",
		Timestamp: gettime.GetTimeStr(),
	}
	allFlagWrite <- loginFailedRegistry
}
