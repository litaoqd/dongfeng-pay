package routers

import (
	"shop/controllers"

	"github.com/beego/beego/v2/server/web"
)

func init() {
	web.Router("/", &controllers.HomeAction{}, "*:ShowHome") //初始化首页
	web.Router("/pay.html", &controllers.PayController{}, "*:GCashPay")
	web.Router("/payment_callback", &controllers.PayController{}, "post:PaymentCallback")
	web.Router("/pay_requst.html", &controllers.ScanShopController{})
	web.Router("/scan.html", &controllers.ScanShopController{}, "*:ScanRender")
	web.Router("/error.html", &controllers.HomeAction{}, "*:ErrorPage")
	web.SetStaticPath("/payment_result.html", "static/payment_result.html")
	web.Router("/check_payment_status", &controllers.PayController{}, "get:CheckPaymentStatus")
	// 使用Beego的方式添加checkout路由
	web.Router("/checkout", &controllers.CheckoutController{}, "post:CheckoutHandler")
	web.Router("/checkout.html", &controllers.CheckoutController{}, "*:CheckoutRender")
	web.Router("/get_qrcode", &controllers.CheckoutController{}, "get:GetQRCode")

}
