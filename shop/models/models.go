package models

import (
	"github.com/beego/beego/v2/adapter/orm"
)

// QRCode represents the qrcodes table in the database
// orm:"qrcodes"
type QrCodes struct {
	Token         string   `orm:"column(token);pk"`
	QRCodePath    string   `orm:"column(qr_code_path)"`
	QRCodeContent string   `orm:"column(qr_code_content)"`
	CreateTime    string   `orm:"column(create_time)"`
	OrderID       string   `orm:"column(order_id)"`
	MerchantID    string   `orm:"column(merchant_id)"`
	ReturnURL     string   `orm:"column(return_url)"`
	_             struct{} `orm:"table(qrcodes)"`
}

func SaveQRCodeRecord(record *QrCodes) error {
	o := orm.NewOrm()
	_, err := o.Insert(record)
	return err
}

func GetQRCodeRecordByToken(token string) (*QrCodes, error) {
	o := orm.NewOrm()
	record := QrCodes{Token: token}
	err := o.Read(&record, "Token")
	if err == orm.ErrNoRows {
		return nil, err
	}
	return &record, nil
}
