package models

import (
	"fmt"
	"github.com/EDDYCJY/go-gin-example/pkg/gormType"
	"github.com/EDDYCJY/go-gin-example/pkg/setting"
	"gorm.io/gorm"
	"log"
	"time"
)

var UDB *gorm.DB

type Paging struct {
	PageNum  int `gorm:"column:-" SDB:"-" json:"page_num" form:"page_num"`
	PageSize int `gorm:"column:-" SDB:"-" json:"page_size" form:"page_size"`
}

type Page struct {
	PageIndex int `json:"page_index" form:"page_index" default:"1"`
	PageSize  int `json:"page_size" form:"page_size" default:"10"`
}

type QueryPaging struct {
	Page
	TenantId  string `json:"tenant_id" valid:"Required;MaxSize(32)" validate:"required" maxLength:"32" example:"10000487779b444bb964985fcba10000"`
	ProjectId string `json:"project_id" valid:"MaxSize(32)" maxLength:"32"`
	OptUserId string `authField:"user_id" json:"opt_user_id" valid:"MaxSize(32)" maxLength:"32" swaggerignore:"true"`
}

type PathItem struct {
	Id        int64  `paramField:"id" json:"id" valid:"Required;Min(1)" validate:"required" minimum:"1"`
	TenantId  string `authField:"tenant_id" json:"tenant_id" valid:"MaxSize(32)" maxLength:"32" localField:"10000487779b444bb964985fcba10000"`
	OptUserId string `authField:"user_id" json:"opt_user_id" valid:"MaxSize(32)" maxLength:"32" localField:"cb18f103df4548a8ac770eda491fe95b"`
}

type PathItemInt struct {
	Id        int    `paramField:"id" json:"id" valid:"Required;Min(1)" validate:"required" minimum:"1"`
	TenantId  string `authField:"tenant_id" json:"tenant_id" valid:"MaxSize(32)" maxLength:"32" localField:"10000487779b444bb964985fcba10000"`
	OptUserId string `authField:"user_id" json:"opt_user_id" valid:"MaxSize(32)" maxLength:"32" localField:"cb18f103df4548a8ac770eda491fe95b"`
}

type PathStringItem struct {
	Id        string `paramField:"id" json:"id" valid:"Required;MaxSize(32)" validate:"required" maxLength:"32"`
	TenantId  string `authField:"tenant_id" json:"tenant_id" valid:"MaxSize(32)" maxLength:"32"`
	OptUserId string `authField:"user_id" json:"opt_user_id" valid:"MaxSize(32)" maxLength:"32"`
}

type DBEntity interface {
}

type ApiResp[T any] struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data T      `json:"data"`
}

// Setup initializes the database instance
func Setup() {
	var err error
	beginTime := time.Now()
	cfg := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local",
		setting.DatabaseSetting.User,
		setting.DatabaseSetting.Password,
		setting.DatabaseSetting.Host,
		setting.DatabaseSetting.Name)
	UDB, err = gormType.NewMySqlDb(cfg)
	if err != nil {
		log.Printf("[mysql]failed to connect mysql,error: %s", err.Error())
	}
	log.Printf("[mysql]build a ticket database,cost %v", time.Since(beginTime))
}
