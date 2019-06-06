package models

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mssql"
	"github.com/log-collection/sqlserver2kafka/pkg/logging"
	"github.com/log-collection/sqlserver2kafka/pkg/setting"
	"go.uber.org/zap"
)
var db *gorm.DB
func Init()  {
	var (
		err error
		dbType,dbName,user,password,host string
	)
	dbType = setting.Config.Database.Type
	dbName = setting.Config.Database.Name
	user = setting.Config.Database.User
	password = setting.Config.Database.Password
	host = setting.Config.Database.Host
	dbInfo := "sqlserver://"+user+":"+password+"@"+host+"?database="+dbName
	logging.Info("数据库信息",zap.String("dbInfo",dbInfo))
	db,err = gorm.Open(dbType,dbInfo)
	if err != nil {
		logging.Panic("open db false",zap.Error(err))
	}
	db.SingularTable(true) // 禁用复数表名
	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(100)
	// 开启sql debug日志
	db.LogMode(true)

}
