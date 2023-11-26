/***************************************************
 ** @Desc : c file for ...
 ** @Time : 2019/8/19 18:13
 ** @Author : yuebin
 ** @File : add
 ** @Last Modified by : yuebin
 ** @Last Modified time: 2019/8/19 18:13
 ** @Software: GoLand
****************************************************/
package controllers

import (
	"boss/models/merchant" // 确保导入正确的模型路径
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	beego "github.com/beego/beego/v2/server/web"
)

type MainController struct {
	beego.Controller
	// ... 其他必要的结构体成员 ...
}

// 商户信息结构
type MerchantInfo struct {
	MerchantUid string
	MerchantKey string
	// ... 其他必要的字段 ...
}

// 商户余额信息结构
type MerchantBalanceInfo struct {
	MerchantName string          `json:"merchantName"`
	BalanceData  json.RawMessage `json:"balanceData"`
}

// 获取所有商户信息和余额
func (c *MainController) GetAllMerchantsData() {
	log.Println("开始获取所有商户信息")

	// 从数据库获取所有商户信息
	merchants := merchant.GetAllMerchant()
	// if err != nil {
	// 	log.Printf("获取商户信息失败: %v\n", err)
	// 	// 处理错误
	// 	return
	// }

	log.Printf("共找到 %d 个商户\n", len(merchants))

	allData := make(map[string]MerchantBalanceInfo)

	// 遍历商户并获取余额信息
	for _, m := range merchants {
		// log.Printf("获取商户 %s 的余额信息\n", m.MerchantUid)
		balanceData, err := getMerchantBalance(m.MerchantUid, m.MerchantKey)
		if err != nil {
			log.Printf("获取商户 %s 的余额信息失败: %v\n", m.MerchantUid, err)
			// 处理错误
			continue
		}

		allData[m.MerchantUid] = MerchantBalanceInfo{
			MerchantName: m.MerchantName, // 假设商户名字字段为 MerchantName
			BalanceData:  balanceData,
		}
	}

	// 将所有数据合并为一个大的 JSON 对象
	data, err := json.Marshal(allData)
	if err != nil {
		log.Printf("合并数据失败: %v\n", err)
		// 处理错误
		return
	}

	log.Printf("合并后的数据: %s\n", string(data))
	// 返回数据给前端
	c.Ctx.WriteString(string(data))
}

// 从 API 获取单个商户的余额信息
func getMerchantBalance(merchantUid, merchantKey string) ([]byte, error) {
	apiUrl := fmt.Sprintf("https://api.onepayph.com/api/v1/balance/%s", merchantUid)
	req, err := http.NewRequest("GET", apiUrl, nil)
	if err != nil {
		log.Printf("创建请求失败: %v\n", err)
		return nil, err
	}

	req.Header.Add("Authorization", "Bearer "+merchantKey)
	log.Printf(apiUrl + ":" + merchantKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("执行请求失败: %v\n", err)
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("读取响应失败: %v\n", err)
		return nil, err
	}

	// log.Printf("body: %s\n", body)

	return json.RawMessage(body), nil
}
