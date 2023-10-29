package controllers

import (
    "fmt"
    "github.com/beego/beego/v2/server/web"
    "net/http"
    "io/ioutil"
    "merchant/models" // 确保导入正确的路径
    "merchant/sys/enum"
)

type QueryOrderController struct {
    web.Controller
}

func (c *QueryOrderController) ShowUI() {
    c.TplName = "query_order.html"
}

func (c *QueryOrderController) SearchOrder() {
    merchantOrderId := c.GetString("merchantOrderId")
    if merchantOrderId == "" {
        c.Data["json"] = map[string]string{"error": "No merchant order ID provided"}
        c.ServeJSON()
        return
    }

    // 获取当前会话中的用户信息
    us := c.GetSession(enum.UserSession)
    if us == nil {
        c.Data["json"] = map[string]string{"error": "User not logged in"}
        c.ServeJSON()
        return
    }
    u := us.(models.MerchantInfo)
    
    fmt.Println("Merchant Key:", u.MerchantKey)

    apiUrl := "https://api.onepayph.com/api/v1/payment_intent/ref_trade_id/" + merchantOrderId

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
        c.Data["json"] = map[string]string{"error": "Failed to fetch order"}
        c.ServeJSON()
        return
    }
    defer resp.Body.Close()

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        c.Data["json"] = map[string]string{"error": "Failed to read response"}
        c.ServeJSON()
        return
    }

    c.Ctx.ResponseWriter.Header().Set("Content-Type", "application/json")
    c.Ctx.ResponseWriter.Write(body)
}
