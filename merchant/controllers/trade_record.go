/***************************************************
 ** @Desc : This file for 交易记录
 ** @Time : 19.12.2 16:34
 ** @Author : Joker
 ** @File : trade_record
 ** @Last Modified by : Joker
 ** @Last Modified time: 19.12.2 16:34
 ** @Software: GoLand
****************************************************/
package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/beego/beego/v2/server/web"
	"merchant/models" // 确保导入正确的路径
	"merchant/sys/enum"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type TradeRecord struct {
	KeepSession
}

func (c *TradeRecord) ShowUI() {
	us := c.GetSession(enum.UserSession)
	u := us.(models.MerchantInfo)

	ranMd5 := encrypt.EncodeMd5([]byte(pubMethod.RandomString(46)))
	c.Ctx.SetCookie(enum.UserCookie, ranMd5, enum.CookieExpireTime)
	c.Ctx.SetSecureCookie(ranMd5, enum.UserCookie, ranMd5, enum.CookieExpireTime)
	c.SetSession(enum.UserCookie, ranMd5)

	c.Data["payType"] = enum.GetPayType()
	c.Data["status"] = enum.GetOrderStatus()
	c.Data["userName"] = u.MerchantName
	c.TplName = "trade_record.html"
}

// TradeQueryAndListPage 方法的修改
func (c *TradeRecord) TradeQueryAndListPage() {
	us := c.GetSession(enum.UserSession)
	u := us.(models.MerchantInfo)

	// 分页参数
	page, _ := strconv.Atoi(c.GetString("page"))
	limit, _ := strconv.Atoi(c.GetString("limit"))
	if limit == 0 {
		limit = 15
	}

	// 查询参数
	in := make(map[string]string)
	merchantNo := strings.TrimSpace(c.GetString("MerchantNo"))
	start := strings.TrimSpace(c.GetString("start"))
	end := strings.TrimSpace(c.GetString("end"))
	payType := strings.TrimSpace(c.GetString("pay_type"))
	status := strings.TrimSpace(c.GetString("status"))

	fmt.Println("！！！！！！！！start:", start)
	fmt.Println("！！！！！！！！end:", end)
	fmt.Println("！！！！！！！！status:", status)

	in["merchant_order_id"] = merchantNo
	in["pay_type_code"] = payType
	in["status"] = status
	in["merchant_uid"] = u.MerchantUid

	// 将日期字符串转换为时间戳
	startTimestamp, err := convertToTimestamp(start)
	if err != nil {
		// 处理错误
	}
	endTimestamp, err := convertToTimestamp(end)
	if err != nil {
		// 处理错误
	}

	// 使用时间戳作为查询参数
	if start != "" {
		in["update_time__gte"] = strconv.FormatInt(startTimestamp, 10)
	}
	if end != "" {
		in["update_time__lte"] = strconv.FormatInt(endTimestamp, 10)
	}

	// 计算分页数
	count := models.GetOrderProfitLenByMap(in)
	totalPage := count / limit // 计算总页数
	if count%limit != 0 {      // 不满一页的数据按一页计算
		totalPage++
	}

	// 数据获取
	var list []models.OrderProfitInfo
	if page <= totalPage {
		list = models.GetOrderProfitByMap(in, limit, (page-1)*limit)
	}

	// 数据回显
	out := make(map[string]interface{})
	out["limit"] = limit // 分页数据
	out["page"] = page
	out["totalPage"] = totalPage
	out["root"] = list // 显示数据

	c.Data["json"] = out
	c.ServeJSON()
	c.StopRun()
}

func (c *TradeRecord) ShowComplaintUI() {
	us := c.GetSession(enum.UserSession)
	u := us.(models.MerchantInfo)

	ranMd5 := encrypt.EncodeMd5([]byte(pubMethod.RandomString(46)))
	c.Ctx.SetCookie(enum.UserCookie, ranMd5, enum.CookieExpireTime)
	c.Ctx.SetSecureCookie(ranMd5, enum.UserCookie, ranMd5, enum.CookieExpireTime)
	c.SetSession(enum.UserCookie, ranMd5)

	c.Data["payType"] = enum.GetPayType()
	c.Data["status"] = enum.GetComOrderStatus()
	c.Data["userName"] = u.MerchantName
	c.TplName = "complaint_record.html"
}

// 投诉列表查询分页
func (c *TradeRecord) ComplaintQueryAndListPage() {
	us := c.GetSession(enum.UserSession)
	u := us.(models.MerchantInfo)

	// 分页参数
	page, _ := strconv.Atoi(c.GetString("page"))
	limit, _ := strconv.Atoi(c.GetString("limit"))
	if limit == 0 {
		limit = 15
	}

	// 查询参数
	in := make(map[string]string)
	merchantNo := strings.TrimSpace(c.GetString("MerchantNo"))
	start := strings.TrimSpace(c.GetString("start"))
	end := strings.TrimSpace(c.GetString("end"))
	payType := strings.TrimSpace(c.GetString("pay_type"))
	status := strings.TrimSpace(c.GetString("status"))

	in["merchant_order_id"] = merchantNo
	in["pay_type_code"] = payType
	if strings.Compare("YES", status) == 0 {
		in["freeze"] = enum.YES
	} else {
		in["refund"] = enum.YES
	}
	in["merchant_uid"] = u.MerchantUid

	if start != "" {
		in["update_time__gte"] = start
	}
	if end != "" {
		in["update_time__lte"] = end
	}

	// 计算分页数
	count := models.GetOrderLenByMap(in)
	totalPage := count / limit // 计算总页数
	if count%limit != 0 {      // 不满一页的数据按一页计算
		totalPage++
	}

	// 数据获取
	var list []models.OrderInfo
	if page <= totalPage {
		list = models.GetOrderByMap(in, limit, (page-1)*limit)
	}

	// 数据回显
	out := make(map[string]interface{})
	out["limit"] = limit // 分页数据
	out["page"] = page
	out["totalPage"] = totalPage
	out["root"] = list // 显示数据

	c.Data["json"] = out
	c.ServeJSON()
	c.StopRun()
}

type TradeRecordController struct {
	web.Controller
}

func (c *TradeRecordController) Fetch() {
	// 获取请求参数
	start := c.GetString("start")
	end := c.GetString("end")
	status := c.GetString("status")
	fmt.Println("****************start:", start)
	fmt.Println("****************end:", end)
	fmt.Println("****************status:", status)

	startTimestamp, err := convertToTimestamp(start)
	if err != nil {
		// 处理错误
	}
	endTimestamp, err := convertToTimestamp(end)
	if err != nil {
		// 处理错误
	}

	fmt.Println("****************startTimestamp:", startTimestamp)
	fmt.Println("****************endTimestamp:", endTimestamp)
	// 获取 token
	token, err := c.GetToken()

	if err != nil {
		c.Data["json"] = map[string]interface{}{"error": "获取 token 失败"}
		c.ServeJSON()
		return
	}

	// 构建 API 请求
	apiUrl := fmt.Sprintf("https://api.onepayph.com/api/v1/disburse_order?start=%d&end=%d&status=%s&page_size=3", startTimestamp, endTimestamp, status)
	fmt.Println("****************apiUrl:", apiUrl)

	client := &http.Client{}
	req, err := http.NewRequest("GET", apiUrl, nil)
	if err != nil {
		c.Data["json"] = map[string]interface{}{"error": "创建请求失败"}
		c.ServeJSON()
		return
	}
	req.Header.Add("Authorization", "Bearer "+token)
	fmt.Println("****************Header:", req.Header)

	// 发起请求
	resp, err := client.Do(req)
	if err != nil {
		c.Data["json"] = map[string]interface{}{"error": "请求 API 失败"}
		c.ServeJSON()
		return
	}
	defer resp.Body.Close()

	// 解析响应
	var result interface{}
	json.NewDecoder(resp.Body).Decode(&result)

	// 返回结果
	c.Data["json"] = result
	fmt.Println("****************result:", result)
	c.ServeJSON()
}

func convertToTimestamp(dateTimeStr string) (int64, error) {
	if dateTimeStr == "" {
		return 0, nil
	}

	layout := "2006-01-02 15:04:00" // Go 的时间格式化布局
	t, err := time.Parse(layout, dateTimeStr)
	if err != nil {
		return 0, err
	}

	return t.UnixMilli(), nil
}

// GetToken 从会话中获取 token
func (c *TradeRecordController) GetToken() (string, error) {
	us := c.GetSession(enum.UserSession)
	if us == nil {
		c.Data["json"] = map[string]string{"error": "User not logged in"}
		c.ServeJSON()
		return "", fmt.Errorf("用户未登录")
	}
	u := us.(models.MerchantInfo)

	fmt.Println("Merchant Key:", u.MerchantKey)
	return u.MerchantKey, nil
}
