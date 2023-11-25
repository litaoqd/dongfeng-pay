package main

import (
	"flag"
	"fmt"
	"log"
	_ "shop/routers"

	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
	"gopkg.in/ini.v1"
)

var ConfPath string

func main() {
	// 定义命令行参数
	flag.StringVar(&ConfPath, "conf", "./conf/app.conf", "config file path")
	flag.Parse()

	fmt.Println("ConfPath:", ConfPath)

	// 加载 INI 配置文件
	cfg, err := ini.Load(ConfPath)
	if err != nil {
		log.Fatalf("Fail to read file: %v", err)
	}

	// 示例：读取一个配置项
	someConfigValue := cfg.Section("mysql").Key("dbbase").String()
	fmt.Println("DBBase:", someConfigValue)
	RegisterLogs()
	beego.Run()
}

/**
** 注册日志信息
 */
func RegisterLogs() {
	logs.SetLogger(logs.AdapterFile,
		`{
						"filename":"../.../logs/legend.log",
						"level":4,
						"maxlines":0,
						"maxsize":0,
						"daily":true,
						"maxdays":10,
						"color":true
					}`)

	f := &logs.PatternLogFormatter{
		Pattern:    "%F:%n|%w%t>> %m",
		WhenFormat: "2006-01-02",
	}

	logs.RegisterFormatter("pattern", f)
	_ = logs.SetGlobalFormatter("pattern")
}
