package utils

import (
	"crypto/rand"
	"encoding/hex"
	"log"
	"os"
	"path/filepath"

	"github.com/beego/beego/v2/core/config"
	"github.com/skip2/go-qrcode"
)

func ensureDir(dirName string) {
	_, err := os.Stat(dirName)
	if os.IsNotExist(err) {
		errDir := os.MkdirAll(dirName, 0755)
		if errDir != nil {
			log.Fatal(err)
		}
	}
}

func GenerateQRCode(content string) (string, string, error) {
	// 从配置文件获取存储路径
	storagePath, _ := config.String("qrCodeStoragePath")
	log.Printf("QR Code storage path: %s\n", storagePath)

	// 生成二维码文件名
	token, err := GenerateToken() // 修改这里
	if err != nil {
		log.Printf("Error generating token for QR code file name: %v\n", err)
		return "", "", err // 处理 GenerateToken 的错误
	}
	qrCodeFileName := "qrcode_" + token + ".png"
	qrCodePath := filepath.Join(storagePath, qrCodeFileName)
	log.Printf("Generated QR code file path: %s\n", qrCodePath)

	// 生成二维码
	// ensureDir(storagePath)
	// err = qrcode.WriteFile(content, qrcode.Medium, 256, qrCodePath)
	// if err != nil {
	// 	log.Printf("Error generating QR code: %v\n", err)
	// 	return "", "", err
	// }

	// 在新的 goroutine 中生成并保存二维码
	go func() {
		ensureDir(storagePath)
		err := qrcode.WriteFile(content, qrcode.Medium, 256, qrCodePath)
		if err != nil {
			log.Printf("Error generating QR code: %v\n", err)
			// 处理错误，例如通过通道发送错误信息或记录错误
		} else {
			log.Println("QR code generated successfully")
		}
	}()

	log.Println("QR code generated successfully")
	return qrCodeFileName, token, nil
}

// GenerateToken 生成一个随机的令牌字符串
func GenerateToken() (string, error) {
	bytes := make([]byte, 16) // 生成一个16字节（128位）的随机数
	if _, err := rand.Read(bytes); err != nil {
		log.Printf("Error generating random token: %v\n", err)
		return "", err
	}
	token := hex.EncodeToString(bytes) // 将随机数转换为十六进制字符串
	log.Printf("Generated token: %s\n", token)
	return token, nil
}
