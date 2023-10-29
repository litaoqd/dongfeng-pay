package controllers

import (
    "github.com/beego/beego/v2/server/web"
)

type QueryOrderController struct {
    web.Controller
}

func (c *QueryOrderController) ShowUI() {
    c.TplName = "query_order.html"
}
