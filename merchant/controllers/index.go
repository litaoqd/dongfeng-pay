/***************************************************
 ** @Desc : This file for 首页
 ** @Time : 19.11.30 11:49
 ** @Author : Joker
 ** @File : index
 ** @Last Modified by : Joker
 ** @Last Modified time: 19.11.30 11:49
 ** @Software: GoLand
****************************************************/
package controllers

import (
	"fmt"
	"merchant/models"
	"merchant/sys/enum"
	"net/http"
    "io/ioutil"
)

type Index struct {
	KeepSession
}

// 首页
func (c *Index) ShowUI() {
	us := c.GetSession(enum.UserSession)
	u := us.(models.MerchantInfo)

	c.Data["userName"] = u.MerchantName
	c.TplName = "index.html"
}

// 加载用户账户金额信息
func (c *Index) LoadUserAccountInfo() {
	us := c.GetSession(enum.UserSession)
	u := us.(models.MerchantInfo)

	ac := models.GetAccountByUid(u.MerchantUid)

	info := make(map[string]interface{})
	// 账户余额
	info["balanceAmt"] = pubMethod.FormatFloat64ToString(ac.Balance)

	// 可用余额
	info["settAmount"] = pubMethod.FormatFloat64ToString(ac.WaitAmount)

	// 冻结金额
	info["freezeAmt"] = pubMethod.FormatFloat64ToString(ac.FreezeAmount)

	// 押款金额
	info["amountFrozen"] = pubMethod.FormatFloat64ToString(ac.LoanAmount)

	c.Data["json"] = info
	c.ServeJSON()
	c.StopRun()
}

// 加载总订单信息
func (c *Index) LoadCountOrder() {
	us := c.GetSession(enum.UserSession)
	u := us.(models.MerchantInfo)

	md := models.GetMerchantDeployByUid(u.MerchantUid)

	type orderInPayWay struct {
		PayWayName    string // 支付方式名
		OrderCount    int    // 订单数
		SucOrderCount int    // 成功订单数
		SucRate       string // 成功率
	}

	ways := make([]orderInPayWay, len(md))

	for k, v := range md {
		in := make(map[string]string)
		in["merchant_uid"] = u.MerchantUid

		ways[k].PayWayName = models.GetRoadInfoByRoadUid(v.SingleRoadUid).ProductName

		in["road_uid"] = v.SingleRoadUid
		ways[k].OrderCount = models.GetOrderLenByMap(in)

		in["status"] = enum.SUCCESS
		ways[k].SucOrderCount = models.GetOrderLenByMap(in)

		if ways[k].OrderCount == 0 {
			ways[k].SucRate = "0"
			continue
		}
		ways[k].SucRate = fmt.Sprintf("%0.4f", float64(ways[k].SucOrderCount)/float64(ways[k].OrderCount))
	}

	c.Data["json"] = ways
	c.ServeJSON()
	c.StopRun()
}

// 加载总订单数
func (c *Index) LoadOrderCount() {
	us := c.GetSession(enum.UserSession)
	u := us.(models.MerchantInfo)

	out := make(map[string]interface{})

	in := make(map[string]string)
	in["merchant_uid"] = u.MerchantUid
	out["orders"] = models.GetOrderLenByMap(in)

	in["status"] = enum.SUCCESS
	out["suc_orders"] = models.GetOrderLenByMap(in)

	if out["orders"].(int) == 0 {
		out["suc_rate"] = 0
	} else {
		out["suc_rate"] = fmt.Sprintf("%0.4f", float64(out["suc_orders"].(int))/float64(out["orders"].(int)))
	}

	c.Data["json"] = out
	c.ServeJSON()
	c.StopRun()
}

// 加载用户支付配置
func (c *Index) LoadUserPayWayUI() {
	us := c.GetSession(enum.UserSession)
	u := us.(models.MerchantInfo)

	c.Data["userName"] = u.MerchantName
	c.TplName = "pay_way.html"
}

func (c *Index) LoadUserPayWay() {
	us := c.GetSession(enum.UserSession)
	u := us.(models.MerchantInfo)

	md := models.GetMerchantDeployByUid(u.MerchantUid)

	type payConfig struct {
		No   string  // 通道编号
		Name string  // 产品名
		Rate float64 // 通道费率
	}

	ways := make([]payConfig, len(md))

	for k, v := range md {
		road := models.GetRoadInfoByRoadUid(v.SingleRoadUid)
		ways[k].No = road.RoadUid

		ways[k].Name = road.ProductName

		ways[k].Rate = road.BasicFee + v.SingleRoadPlatformRate + v.SingleRoadAgentRate
	}

	c.Data["json"] = ways
	c.ServeJSON()
	c.StopRun()
}

// GetBalance 处理获取余额的请求
func (c *Index) GetBalance() {
    // 获取当前会话中的用户信息
    us := c.GetSession(enum.UserSession)
    if us == nil {
        c.Data["json"] = map[string]string{"error": "User not logged in"}
        c.ServeJSON()
        return
    }
    u := us.(models.MerchantInfo)

    // 构建 API URL
    apiUrl := "https://api.onepayph.com/api/v1/balance/" + u.MerchantUid

    // 创建一个新的 HTTP 请求
    req, err := http.NewRequest("GET", apiUrl, nil)
    if err != nil {
        c.Data["json"] = map[string]string{"error": "Failed to create request"}
        c.ServeJSON()
        return
    }

    // 添加 Bearer Token 到请求头
    req.Header.Add("Authorization", "Bearer " + u.MerchantKey)

    // 发送请求
    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        c.Data["json"] = map[string]string{"error": "Failed to fetch balance"}
        c.ServeJSON()
        return
    }
    defer resp.Body.Close()

    // 读取响应体
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        c.Data["json"] = map[string]string{"error": "Failed to read response"}
        c.ServeJSON()
        return
    }

    // 将响应体直接写入到响应中
    c.Ctx.ResponseWriter.Header().Set("Content-Type", "application/json")
    c.Ctx.ResponseWriter.Write(body)
}