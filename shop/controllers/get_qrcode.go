package controllers

import (
	"net/http"
	"shop/models"
)

func GetQRCodeHandler(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")

	// 从数据库获取二维码信息
	qrCodeRecord, err := models.GetQRCodeRecordByToken(token)
	if err != nil {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	// 返回二维码图片路径
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"qrcodePath":"` + qrCodeRecord.QRCodePath + `"}`))
}
