package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"shop/models"
	"shop/utils"

	"github.com/beego/beego/v2/core/config"
	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/server/web"
)

type CheckoutController struct {
	web.Controller
}

type CheckoutRequest struct {
	QRCodeContent string `json:"qrcode"`
	CreateTime    string `json:"create_time"`
	OrderID       string `json:"order_id"`
	MerchantID    string `json:"merchant_id"`
	ReturnURL     string `json:"return_url"`
}

func (c *CheckoutController) CheckoutRender() {
	// 添加逻辑来准备数据
	c.TplName = "pay/checkout.html" // 确保这里的路径是正确的
}

func (c *CheckoutController) GetQRCode() {
	token := c.GetString("token")
	logs.Info("Received request for QR code with token:", token)

	// 从数据库获取二维码信息
	qrCodeRecord, err := models.GetQRCodeRecordByToken(token)
	if err != nil {
		logs.Error("QR code not found for token:", token, "Error:", err)
		c.Ctx.ResponseWriter.WriteHeader(http.StatusNotFound)
		c.Ctx.WriteString("Not Found")
		return
	}

	// 从配置文件获取 qrCodeGetPath
	qrCodeGetPath, err := web.AppConfig.String("qrCodeGetPath")
	if err != nil {
		logs.Error("Error retrieving qrCodeGetPath from config:", err)
		c.Ctx.ResponseWriter.WriteHeader(http.StatusInternalServerError)
		c.Ctx.WriteString("Internal Server Error")
		return
	}

	// 拼接完整的二维码路径
	fullQRCodePath := qrCodeGetPath + "/" + qrCodeRecord.QRCodePath
	logs.Info("Full QR code path:", fullQRCodePath)

	c.Data["json"] = map[string]string{"qrcodePath": fullQRCodePath}
	c.ServeJSON()
}

func (c *CheckoutController) CheckoutHandler() {
	// 读取请求体
	body, err := ioutil.ReadAll(c.Ctx.Request.Body)
	if err != nil {
		logs.Error("Error reading request body:", err)
		c.Ctx.ResponseWriter.WriteHeader(http.StatusInternalServerError)
		return
	}

	var req CheckoutRequest
	if err := json.Unmarshal(body, &req); err != nil {
		log.Printf("Error parsing JSON request: %v\n", err)
		c.Ctx.ResponseWriter.WriteHeader(http.StatusBadRequest)
		return
	}

	// log.Printf("request: %v\n", req)

	// 使用 req 中的字段代替 c.GetString()
	qrCodeContent := req.QRCodeContent
	createTime := req.CreateTime
	orderID := req.OrderID
	merchantID := req.MerchantID
	returnUrl := req.ReturnURL

	log.Printf("Received checkout request: qrCodeContent=%s, createTime=%s, orderID=%s, merchantID=%s, returnUrl=%s\n", qrCodeContent, createTime, orderID, merchantID, returnUrl)

	// 生成二维码图片并保存
	qrCodePath, err := utils.GenerateQRCode(qrCodeContent)
	if err != nil {
		log.Printf("Error generating QR code: %v\n", err)
		c.Ctx.ResponseWriter.WriteHeader(http.StatusInternalServerError)
		return
	}

	// 生成Token
	token, err := utils.GenerateToken()
	if err != nil {
		log.Printf("Error generating token: %v\n", err)
		c.Ctx.ResponseWriter.WriteHeader(http.StatusInternalServerError)
		return
	}

	// 存储信息到数据库
	qrCodeRecord := models.QrCodes{
		Token:         token,
		QRCodePath:    qrCodePath,
		QRCodeContent: qrCodeContent,
		CreateTime:    createTime,
		OrderID:       orderID,
		MerchantID:    merchantID,
		ReturnURL:     returnUrl,
	}
	err = models.SaveQRCodeRecord(&qrCodeRecord)
	if err != nil {
		log.Printf("Error saving QR code record to database: %v\n", err)
		c.Ctx.ResponseWriter.WriteHeader(http.StatusInternalServerError)
		return
	}

	// 构建返回的URL
	checkoutPageURL, _ := config.String("checkoutPageURL")
	responseURL := fmt.Sprintf("%s?token=%s", checkoutPageURL, token)

	log.Printf("Checkout URL generated: %s\n", responseURL)

	// 返回URL给OnePay
	// c.Ctx.WriteString(responseURL)

	// 构建返回给OnePay的JSON对象
	response := map[string]string{
		"url":   responseURL,
		"token": token,
	}

	// 设置响应类型为JSON
	c.Ctx.ResponseWriter.Header().Set("Content-Type", "application/json")
	c.Ctx.ResponseWriter.WriteHeader(http.StatusOK)

	// 返回JSON给OnePay
	json.NewEncoder(c.Ctx.ResponseWriter).Encode(response)

}
