/***************************************************
 ** @Desc : This file for 首页js
 ** @Time : 19.12.2 14:46
 ** @Author : Joker
 ** @File : index
 ** @Last Modified by : Joker
 ** @Last Modified time: 19.12.2 14:46
 ** @Software: GoLand
 ****************************************************/

let index = {
    getAccountInfo: function() {
        $.ajax({
            type: "GET",
            url: "/index/loadInfo/",
            success: function(result) {
                $("#balanceAmt").text(result["balanceAmt"]);
                $("#settAmount").text(result["settAmount"]);
                $("#freezeAmt").text(result["freezeAmt"]);
                $("#amountFrozen").text(result["amountFrozen"]);
            },
            error: function(XMLHttpRequest) {
                toastr.info('something is wrong, code: ' + XMLHttpRequest.status);
            }
        });
    },
    getOrdersInfo: function() {
        $.ajax({
            type: "GET",
            url: "/index/loadOrders",
            success: function(result) {
                $("#orders").text(result["orders"]);
                $("#suc_orders").text(result["suc_orders"]);
                $("#suc_rate").text((result["suc_rate"] * 100).toFixed(2) + "%");
            },
            error: function(XMLHttpRequest) {
                toastr.info('something is wrong, code: ' + XMLHttpRequest.status);
            }
        });
    },
    get_account_balance: function() {
        $.ajax({
            type: "GET",
            url: "/withdraw/balance",
            success: function(resp) {
                $("#balance").val(resp.balance);
                $("#sett_fee").html(resp.fee);
            },
            error: function(XMLHttpRequest) {
                toastr.info('something is wrong, code: ' + XMLHttpRequest.status);
            }
        });
    },
    loadTradeRecord: function() {
        $.ajax({
            type: "GET",
            url: "/index/load_count_order",
            success: function(res) {
                let con = "";
                $.each(res, function(index, item) {
                    if (item.PayWayName === "") { return true; }
                    con += `<div class="project"><div class="row bg-white has-shadow"><div class="left-col col-lg-1 d-flex align-items-center justify-content-between"><small>` + (index + 1) + `</small></div><div class="left-col col-lg-3 d-flex align-items-center justify-content-between"><small>` + item.PayWayName + `</small></div><div class="left-col col-lg-2 d-flex align-items-center justify-content-between"><small>` + item.OrderCount + `</small></div><div class="left-col col-lg-2 d-flex align-items-center justify-content-between"><small>` + item.SucOrderCount + `</small></div><div class="left-col col-lg-2 d-flex align-items-center justify-content-between"><small>` + (item.SucRate * 100).toFixed(2) + `%</small></div></div></div>`;
                });
                $("#your_showtime").html(con);
            },
            error: function(XMLHttpRequest) {
                toastr.info('something is wrong, code: ' + XMLHttpRequest.status);
            }
        });
    },
    loadPayWay: function() {
        $.ajax({
            type: "GET",
            url: "/index/pay_way",
            success: function(res) {
                let con = "";
                $.each(res, function(index, item) {
                    if (item.Name === "") { return true; }
                    con += `<div class="project"><div class="row bg-white has-shadow"><div class="left-col col-lg-1 d-flex align-items-center justify-content-between"><small>` + (index + 1) + `</small></div><div class="left-col col-lg-3 d-flex align-items-center justify-content-between"><small>` + item.Name + `</small></div><div class="left-col col-lg-2 d-flex align-items-center justify-content-between"><small>` + item.Rate + `%</small></div></div></div>`;
                });
                $("#your_showtime").html(con);
            },
            error: function(XMLHttpRequest) {
                toastr.info('something is wrong, code: ' + XMLHttpRequest.status);
            }
        });
    },
    
    updateBalance: function() {
        console.log("!!!!!!!!!!!!!!!!updateBalance function called");
        $.ajax({
            url: "/index/get_balance",
            type: "GET",
            success: function(response) {
                // 更新账户余额和押款金额
                $("#balanceAmt").text(response.cash_balance);
                $("#amountFrozen").text(response.receivable_balance);
                console.log("Balance Amount: ", response.cash_balance);
                console.log("Amount Frozen: ", response.receivable_balance);
            },
            error: function(error) {
                console.log("Error fetching balance: ", error);
            }
        });
    },
    
    init: function() {
        // 初始加载一次余额
        this.updateBalance();

        // 每分钟更新余额
        setInterval(this.updateBalance.bind(this), 60000);
    }
};

// 页面加载时调用 init 方法
$(function() {
    index.init();
});
