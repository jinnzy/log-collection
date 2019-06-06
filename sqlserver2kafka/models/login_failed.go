package models

import (
	"github.com/jinzhu/gorm"
)

type LoginFailed struct {
	Date string `gorm:"column:LogDate"`
	ProcessInfo string `column:"ProcessInfo"`
	Text string `gorm:"column:Text"`
	DbHostname string `json:"ServerName"`
}

func (LoginFailed) TableName() string {
	// 返回真正的表名
	return "Login_failed"
}

func GetLoginFailed(aID int,tmpID int) (loginFailed []LoginFailed,err error) {
	//err = db.Order("LogDate")
	//err = db.Select("LogDate,ProcessInfo,Text").Offset(aID).Limit(tmpID).Find(&loginFailed).Error
	err = db.Raw("select * from(        SELECT ROW_NUMBER() OVER(ORDER BY [LogDate])RN,*    FROM [audit].[dbo].[Login_failed])a where RN between ? and ?",aID,tmpID).Scan(&loginFailed).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return loginFailed,err
	}
	return loginFailed,nil
}