package models

import (
	"github.com/jinzhu/gorm"
)

type DmlAll struct {
	AID int `gorm:"column:AID"`
	ServerName string `gorm:"column:Server_Name"`
	EventTime string `gorm:"column:EventTime"`
	EventName string `gorm:"column:EventName"`
	ClientAppName string `gorm:"column:ClientAppName"`
	ClientHostName string `gorm:"column:ClientHostName"`
	BatchText string `gorm:"column:batch_text"`
	DatabaseName string `gorm:"column:DatabaseName"`
	UserName string `gorm:"column:UserName"`
}

func (DmlAll) TableName() string {
	// 返回真正的表名
	return "dml_all"
}

func GetDml(aID int,tmpID int) (dmlAll []DmlAll,err error) {
	err = db.Select("AID,Server_Name,EventTime,EventName,ClientAppName,ClientHostName,batch_text,DatabaseName,UserName",).Where("AID BETWEEN ? AND ?", aID,tmpID).Find(&dmlAll).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		// 返回值为列表，则不返回gorm.ErrRecordNotFound这个了
		// 这里是不管ErrRecordNotFound这个错误，在调用这个函数的时候判断下dmlAll是否为空
		return dmlAll,err
		//return dmlAllNil,err
	}
	return dmlAll,nil

}