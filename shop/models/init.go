/***************************************************
 ** @Desc : This file for ...
 ** @Time : 2019/8/9 13:48
 ** @Author : yuebin
 ** @File : init
 ** @Last Modified by : yuebin
 ** @Last Modified time: 2019/8/9 13:48
 ** @Software: GoLand
****************************************************/
package models

import (
	"fmt"

	"github.com/beego/beego/v2/adapter/orm"
	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/server/web"
	_ "github.com/go-sql-driver/mysql"
)

func init() {
	dbHost, _ := web.AppConfig.String("mysql::dbhost")
	dbUser, _ := web.AppConfig.String("mysql::dbuser")
	dbPassword, _ := web.AppConfig.String("mysql::dbpasswd")
	dbBase, _ := web.AppConfig.String("mysql::dbbase")
	dbPort, _ := web.AppConfig.String("mysql::dbport")

	link := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8", dbUser, dbPassword, dbHost, dbPort, dbBase)

	logs.Info("mysql init.....", link)

	_ = orm.RegisterDriver("mysql", orm.DRMySQL)
	_ = orm.RegisterDataBase("default", "mysql", link, 30, 30)
	//注册qrcodes表
	orm.RegisterModel(new(QrCodes))
}
