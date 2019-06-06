package models

import (
	"github.com/jinzhu/gorm"
)

type DdlAll struct {
	AID int `gorm:"column:AID"`
	EventType string `column:"EventType"`
	EventTime string `gorm:"column:PostTime"`
	ServerName string `gorm:"column:ServerName"`
	UserName string `gorm:"column:LoginName"`
	DatabaseName string `gorm:"column:DatabaseName"`
	ObjectName string `gorm:"column:ObjectName"`
	TSQLCommand string `gorm:"column:TSQLCommand"`
	ClientAppName string `gorm:"column:client_program"`
	ClientHostName string `gorm:"column:client_name"`
	ClientIP string `gorm:"column:client_IP"`
}

func (DdlAll) TableName() string {
	// 返回真正的表名
	return "ddl_all"
}

func GetDdl(aID int,tmpID int) (ddlAll []DdlAll,err error) {
	err = db.Select("AID,EventType,PostTime,ServerName,LoginName,DatabaseName,ObjectName,TSQLCommand,client_program,client_name,client_IP").Where("AID BETWEEN ? AND ?", aID,tmpID).Find(&ddlAll).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		// 返回值为列表，则不返回gorm.ErrRecordNotFound这个了
		// 这里是不管ErrRecordNotFound这个错误，在调用这个函数的时候判断下ddlAll是否为空
		return ddlAll,err
		//return ddlAllNil,err
	}
	return ddlAll,nil

}