/***************************************************
 ** @Desc : This file for ...
 ** @Time : 2019/12/18 17:16
 ** @Author : yuebin
 ** @File : pay
 ** @Last Modified by : yuebin
 ** @Last Modified time: 2019/12/18 17:16
 ** @Software: GoLand
****************************************************/
package controllers

import (
	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/server/web"
	"strconv"
	"strings"
	"net/http"
	"encoding/json" // This is for json
    "bytes"        // This is for bytes
)

type PayController struct {
	web.Controller
}

func (c *PayController) Pay() {
	orderNo := strings.TrimSpace(c.GetString("orderid"))
	flash := web.NewFlash()
	if orderNo == "" {
		flash.Error("订单号为空")
		flash.Store(&c.Controller)
		c.Redirect("/error.html", 302)
		return
	}
	amount := strings.TrimSpace(c.GetString("amount"))
	if !c.judgeAmount(amount) {
		flash.Error("金额有误")
		flash.Store(&c.Controller)
		c.Redirect("/error.html", 302)
		return
	}
	isScan := strings.TrimSpace(c.GetString("SCAN"))
	isH5 := strings.TrimSpace(c.GetString("H5"))
	isKj := strings.TrimSpace(c.GetString("KJ"))
	if strings.Contains(isScan, "SCAN") {
		//扫码
		scanShop := new(ScanShopController)
		scanShop.Prepare()
		scanShop.Params["orderPrice"] = amount
		scanShop.Params["payWayCode"] = isScan
		scanShop.Params["orderNo"] = orderNo
		response := scanShop.Shop(c.Ctx.Request.Host)
		if response.Code == 200 {
			str := "/scan.html?" + "orderNo=" + orderNo + "&orderPrice=" + amount + "&qrCode=" + response.Qrcode + "&payWayCode=" + isScan
			c.Redirect(str, 302)
		} else {
			flash.Error(response.Msg)
			flash.Store(&c.Controller)
			c.Redirect("/error.html", 302)
		}
	} else if strings.Contains(isH5, "H5") {

	} else if strings.Contains(isKj, "FAST") {

	} else {
		flash.Error("不存在这样的支付类型")
		flash.Store(&c.Controller)
		c.Redirect("/error.html", 302)
		return
	}
}

func (c *PayController) judgeAmount(amount string) bool {
	_, err := strconv.ParseFloat(amount, 64)
	if err != nil {
		logs.Error("输入金额有误")
		return false
	}

	return true
}

// GCashPay 处理GCash支付请求
func (c *PayController) GCashPay() {
    // 从表单中获取订单编号
    orderNo := c.GetString("orderid")
    // 在订单编号前添加 'test'
    refTradeID := "test" + orderNo

    // 构建请求体
    requestBody := map[string]interface{}{
        "ref_trade_id": refTradeID,
        "amount":       "100",
        "currency":     "php",
        "callback_url": "This is a test callback_url - pay in",
        "customer": map[string]string{
            "first_name": "wowo",
            "last_name":  "totot",
            "email":      "wwwd@sina.com",
        },
        "payment_method": map[string]interface{}{
            "type": "direct_debit",
            "direct_debit": map[string]string{
                "channel": "jcpay_qrcode",
            },
        },
        "description": "AbC",
    }

    // 将请求体编码为JSON
    jsonValue, _ := json.Marshal(requestBody)

    // 创建POST请求
    request, err := http.NewRequest("POST", "https://api.onepayph.com/api/v1/payment_intent", bytes.NewBuffer(jsonValue))
    if err != nil {
        logs.Error("创建请求失败:", err)
        c.Ctx.ResponseWriter.WriteHeader(http.StatusInternalServerError)
        return
    }

    // 添加认证头部
    request.Header.Set("Authorization", "Bearer juancash")
    request.Header.Set("Content-Type", "application/json")

    // 发送请求
    client := &http.Client{}
    response, err := client.Do(request)
    if err != nil {
        logs.Error("发送请求失败:", err)
        c.Ctx.ResponseWriter.WriteHeader(http.StatusInternalServerError)
        return
    }
    defer response.Body.Close()

    // 读取响应
    if response.StatusCode == http.StatusOK {
        var result map[string]interface{}
        json.NewDecoder(response.Body).Decode(&result)
        // 处理结果，例如发送二维码到前端
        c.Data["json"] = result["qr_code_content"]
        c.ServeJSON()
    } else {
        // 处理错误
        c.Ctx.ResponseWriter.WriteHeader(response.StatusCode)
    }
}
