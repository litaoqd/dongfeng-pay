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
	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/ini.v1"
)

// 假设 ConfPath 是在 main.go 中定义的全局变量
func init() {
	// 加载 INI 配置文件
	fmt.Println("ConfPath:", ConfPath)
	cfg, err := ini.Load(ConfPath)
	if err != nil {
		logs.Error("Fail to read file: %v", err)
		return
	}

	// 从配置文件读取数据库配置
	dbHost := cfg.Section("mysql").Key("dbhost").String()
	dbUser := cfg.Section("mysql").Key("dbuser").String()
	dbPassword := cfg.Section("mysql").Key("dbpasswd").String()
	dbBase := cfg.Section("mysql").Key("dbbase").String()
	dbPort := cfg.Section("mysql").Key("dbport").String()

	link := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8", dbUser, dbPassword, dbHost, dbPort, dbBase)

	logs.Info("mysql init.....", link)

	_ = orm.RegisterDriver("mysql", orm.DRMySQL)
	_ = orm.RegisterDataBase("default", "mysql", link, 30, 30)
	// 注册qrcodes表
	orm.RegisterModel(new(QrCodes))
}
