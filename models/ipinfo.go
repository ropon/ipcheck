package models

import (
	"github.com/jinzhu/gorm"
	. "ipinfo/databases"
	"time"
)

type IpInfo struct {
	ID        uint   `json:"-" gorm:"primary_key,AUTO_INCREMENT"`
	Start     uint32 `json:"-" gorm:"column:start" binding:"required"`
	End       uint32 `json:"-" gorm:"column:end" binding:"required"`
	Ip        string `json:"ip" gorm:"-"`
	Zone      string `json:"zone" gorm:"column:zone"`
	Isp       string `json:"isp" gorm:"column:isp"`
	CacheTime string `json:"cachetime" gorm:"column:cachetime"`
}

func (i IpInfo) TableName() string {
	return "ipinfo"
}

//创建时更新缓存时间
func (i IpInfo) BeforeCreate(scope *gorm.Scope) error {
	_ = scope.SetColumn("CacheTime", time.Now().Format("2006-01-02 15:04:05"))
	return nil
}

//更新时更新缓存时间
func (i IpInfo) BeforeUpdate(scope *gorm.Scope) error {
	_ = scope.SetColumn("CacheTime", time.Now().Format("2006-01-02 15:04:05"))
	return nil
}

func Migrate() {
	Db.AutoMigrate(&IpInfo{})
}

//创建
func (i IpInfo) Insert() bool {
	created := Db.Create(&i)
	if created.Error != nil {
		return false
	}
	return true
}

//更新 根据给定的条件更新单个属性

func IpInfoGet(ipNum uint32) (i IpInfo) {
	Db.Where("end >= ? AND start <= ? ", ipNum, ipNum).First(&i)
	return
}
