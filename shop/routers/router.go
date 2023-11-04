package routers

import (
	"github.com/beego/beego/v2/server/web"
	"shop/controllers"
)

func init() {
	web.Router("/", &controllers.HomeAction{}, "*:ShowHome") //初始化首页
	web.Router("/pay.html", &controllers.PayController{}, "*:GCashPay")
	web.Router("/payment_callback", &controllers.PayController{}, "post:PaymentCallback")
	web.Router("/pay_requst.html", &controllers.ScanShopController{})
	web.Router("/scan.html", &controllers.ScanShopController{}, "*:ScanRender")
	web.Router("/error.html", &controllers.HomeAction{}, "*:ErrorPage")
	web.SetStaticPath("/payment_result.html", "static/payment_result.html")
}
