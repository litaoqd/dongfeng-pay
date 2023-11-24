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
	c.Ctx.WriteString(responseURL)
}
