/***************************************************
 ** @Desc : This file for 银行编码
 ** @Time : 19.12.4 10:42
 ** @Author : Joker
 ** @File : bank_info
 ** @Last Modified by : Joker
 ** @Last Modified time: 19.12.4 10:42
 ** @Software: GoLand
****************************************************/
package enum

const (
	ICBC   = "ICBC"
	ABC    = "ABC"
	BOC    = "BOC"
	CCB    = "CCB"
)

var bankInfo = map[string]string{
	ICBC:   "BDO",
	ABC:    "UNION BANK",
	BOC:    "BPI",
	CCB:    "GCASH",
}

func GetBankInfo() map[string]string {
	return bankInfo
}
