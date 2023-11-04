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
    "net/url"
    "io/ioutil"
    "fmt"
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

// CheckPaymentStatus 是一个新的方法，用于检查支付状态
func (c *PayController) CheckPaymentStatus(refTradeID string) (string, error) {
	// 构建请求URL
	apiURL := "https://api.onepayph.com/api/v1/payment_intent/" + refTradeID

	// 创建HTTP客户端
	client := &http.Client{}

	// 创建请求
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		logs.Error("创建请求失败:", err)
		return "", err
	}

	// 设置认证头部
	req.Header.Set("Authorization", "Bearer juancash")

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		logs.Error("发送请求失败:", err)
		return "", err
	}
	defer resp.Body.Close()

	// 检查响应状态码
	if resp.StatusCode != http.StatusOK {
		logs.Error("请求失败，状态码：", resp.StatusCode)
		return "", err
	}

	// 解析响应体为JSON
	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		logs.Error("解析响应体失败:", err)
		return "", err
	}

	// 获取并返回状态字段
	status, ok := result["status"].(string)
	if !ok {
		logs.Error("响应体中不存在status字段")
		return "", err
	}

	return status, nil
}

// PaymentCallback 处理上游渠道的支付回调
func (c *PayController) PaymentCallback() {
    // 记录接收到的回调请求
    logs.Info("Received payment callback")

    // 读取请求体
    body, err := ioutil.ReadAll(c.Ctx.Request.Body)
    if err != nil {
        logs.Error("Error reading request body:", err)
        c.Ctx.ResponseWriter.WriteHeader(http.StatusInternalServerError)
        return
    }
    logs.Info("Callback request body:", string(body))

    // 解析请求体中的JSON数据
    var callbackData map[string]interface{}
    err = json.Unmarshal(body, &callbackData)
    if err != nil {
        logs.Error("Error unmarshalling request body:", err)
        c.Ctx.ResponseWriter.WriteHeader(http.StatusBadRequest)
        return
    }
    logs.Info("Parsed callback data:", callbackData)

    // 验证回调的真实性（示例中省略了验证步骤）

    // 更新订单状态等业务逻辑（示例中省略了业务逻辑）

    // 记录业务处理的结果
    logs.Info("Payment callback processed successfully")

    // 重定向到显示结果的页面，并附带回调数据
    c.Data["json"] = callbackData
    // c.ServeJSON()
    // 构建查询参数
    params := url.Values{}
    for key, value := range callbackData {
        // 假设回调数据中的值都是字符串类型
        params.Add(key, fmt.Sprintf("%v", value))
    }
    queryString := params.Encode()

    // // 重定向到显示结果的页面，并附带回调数据
    // logs.Info("重定向到显示结果的页面，并附带回调数据")
    // c.Redirect("/payment_result.html?" + queryString, 302)
    // 返回 JSON 响应而不是重定向
    c.Data["json"] = map[string]string{"redirect": "/payment_result.html?" + queryString}
    c.ServeJSON()
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
        "amount":       "10",
        "currency":     "php",
        "callback_url": "http://merchant.onepayph.com:12308/payment_callback",
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
    logs.Info("请求成功！")
    // 读取响应
    if response.StatusCode == http.StatusOK {
        var result map[string]interface{}
        json.NewDecoder(response.Body).Decode(&result)
        
        // 提取二维码内容
        qrCodeContent, ok := result["qr_code_content"].(string)
        if !ok {
            logs.Error("二维码内容不存在")
            c.Ctx.ResponseWriter.WriteHeader(http.StatusInternalServerError)
            return
        }

        
        // 重定向到scan.html，并传递二维码内容
        // c.Redirect("/scan.html?qrCode="+url.QueryEscape(qrCodeContent), 302)
        c.Redirect("/scan.html?qrCode="+url.QueryEscape(qrCodeContent)+"&refTradeID="+url.QueryEscape(refTradeID), 302)

    } else {
        // 处理错误
        c.Ctx.ResponseWriter.WriteHeader(response.StatusCode)
    }
}


