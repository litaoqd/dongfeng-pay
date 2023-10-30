package controllers

import (
    "fmt"
    "github.com/beego/beego/v2/server/web"
    "github.com/dchest/captcha"
    "merchant/models"
    "merchant/sys"
    "merchant/sys/enum"
    "merchant/utils"
    "strings"
    "github.com/dgryski/dgoogauth"
)

var pubMethod = sys.PublicMethod{}
var encrypt = utils.Encrypt{}

type Login struct {
    web.Controller
}

func (c *Login) UserLogin() {
    // captchaCode := c.GetString("captchaCode")
    // captchaId := c.GetString("captchaId")
    userName := strings.TrimSpace(c.GetString("userName"))
    password := c.GetString("Password")
    dynamicCode := c.GetString("dynamicCode")

    var (
        flag = enum.FailedFlag
        msg  = ""
        url  = "/"

        pwdMd5 string
        ran    string
        ranMd5 string

        // verify bool
        u models.MerchantInfo
    )

    us := c.GetSession(enum.UserSession)
    if us != nil {
        url = enum.DoMainUrl
        flag = enum.SuccessFlag
        c.respondAndExit(flag, msg, url)
        return
    }

    if userName == "" || password == "" {
        msg = "登录账号或密码不能为空!"
        c.respondAndExit(enum.FailedFlag, msg, "/")
        return
    }


    u = models.GetMerchantByPhone(userName)
    if u.LoginPassword == "" {
        msg = "账户信息错误，请联系管理人员!"
        c.respondAndExit(enum.FailedFlag, msg, "/")
        return
    }

    if strings.Compare(enum.ACTIVE, u.Status) != 0 {
        msg = "登录账号或密码错误!"
        c.respondAndExit(enum.FailedFlag, msg, "/")
        return
    }

    // 验证密码
    pwdMd5 = encrypt.EncodeMd5([]byte(password))
    if strings.Compare(strings.ToUpper(pwdMd5), u.LoginPassword) != 0 {
        msg = "登录账号或密码错误!"
        c.respondAndExit(enum.FailedFlag, msg, "/")
        return
    }

    // 验证 Google Authenticator OTP
    otpConfig := &dgoogauth.OTPConfig{
        Secret:      u.MerchantSecret,
        WindowSize:  3,
        HotpCounter: 0,
    }
    
    fmt.Println("Merchant Key:", otpConfig.Secret)
    fmt.Println("dynamicCode:", dynamicCode)    
    valid, err := otpConfig.Authenticate(dynamicCode)
    if err != nil || !valid {
        msg = "动态验证码错误或已过期!"
        c.respondAndExit(enum.FailedFlag, msg, "/")
        return
    }

    c.SetSession(enum.UserSession, u)

    // 设置客户端用户信息有效保存时间
    ran = pubMethod.RandomString(46)
    ranMd5 = encrypt.EncodeMd5([]byte(ran))
    c.Ctx.SetSecureCookie(ranMd5, enum.UserCookie, ranMd5, enum.CookieExpireTime)
    c.SetSession(enum.UserCookie, ranMd5)

    url = enum.DoMainUrl
    flag = enum.SuccessFlag
    utils.LogNotice(fmt.Sprintf("&#8203;``【oaicite:0】``&#8203;用户登录成功，请求IP：%s", u.MerchantName, c.Ctx.Input.IP()))

    c.respondAndExit(flag, msg, url)
}

func (c *Login) respondAndExit(flag int, msg string, url string) {
    c.Data["json"] = pubMethod.JsonFormat(flag, "", msg, url)
    _ = c.ServeJSON()
    c.StopRun()
}

// 这里添加您项目中其他已有的方法
// 例如：Logout, ChangePassword, ResetPassword 等
// 确保这些方法与您项目的其他部分保持一致


func (c *Login) Index() {
	capt := struct {
		CaptchaId string
	}{
		captcha.NewLen(4),
	}
	c.Data["CaptchaId"] = capt.CaptchaId
	c.TplName = "login.html"
}

// 验证输入的验证码
func (c *Login) VerifyCaptcha() {
	captchaValue := c.Ctx.Input.Param(":value")
	captchaId := c.Ctx.Input.Param(":chaId")

	verify := captcha.VerifyString(captchaId, captchaValue)
	if verify {
		c.Data["json"] = pubMethod.JsonFormat(enum.SuccessFlag, "", "", "")
	} else {
		c.Data["json"] = pubMethod.JsonFormat(enum.FailedFlag, "", "验证码不匹配!", "")
	}
	c.ServeJSON()
	c.StopRun()
}

// 重绘验证码
func (c *Login) FlushCaptcha() {
	capt := struct {
		CaptchaId string
	}{
		captcha.NewLen(4),
	}
	c.Data["json"] = pubMethod.JsonFormat(enum.SuccessFlag, capt.CaptchaId, "验证码不匹配!", "")
	c.ServeJSON()
	c.StopRun()
}

// 退出登录
func (c *Login) LoginOut() {
	c.DelSession(enum.UserSession)

	c.Data["json"] = pubMethod.JsonFormat(enum.FailedFlag, "", "", "/")
	c.ServeJSON()
	c.StopRun()
}

// 对接文档
func (c *Login) PayDoc() {
	c.TplName = "api_doc/pay_doc.html"
}
