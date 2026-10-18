package itying

import (
	"fmt"
	"math"
	"strconv"

	"ginshop/models"

	"github.com/gin-gonic/gin"
	"github.com/smartwalle/alipay/v3"
)

type AlipayController struct{}

func (con AlipayController) Alipay(c *gin.Context) {
	//1、获取订单号 判断此订单号是否值当前用户的
	//2、获取订单里面的支付信息

	var privateKey = "MIIEvgIBADANBgkqhkiG9w0BAQEFAASCBKgwggSkAgEAAoIBAQCPUi2yj4RJcqAKfjV/AxjqIjjUq007pGhLClbrwEoCVhgnQU9bAVFnPlDaVdO2xBfu8D/gHCwV9czZKJh51yg1kwE6Yv/hSd966PnGS9NszfxxWfTbWeC6DIHZj66nTECK+vYWX36tKxIG+juzXyfoAuyL1h58oFgF/Zwa1FzpKDTcMUIE+npf3DpMS0Uatf6TsREzVQQs3i0WVxYY3lv1Dmualr6Q3GAFE1j/xt/STE993uh8MVeLS+RcTrrjilPSVH5Z0DLAjkDSH7XUK1lIHASpgOrddEJ8MeT8L7bUJ8nDs5qQ8zfbPVJaXfsF2NYS162HkIW2bl0r5c5Mudt7AgMBAAECggEAWUHVweHNgiyH7WECkhJsvswHVrNEi0NtzGYpEfOUY/YYXsI22Lduaf0OP5u6GZXwTdeEAF+rORX2uLumkiLkINFnr2QedcEbFCHqBIwOpTF36WQbsUw9P8EwUT1BiWFcxPFctzxL2S78sCnBaol1gfHoPYJhRD5b84cpZDAjmPSJk1XAtAtKChUIskLBAsCvwlGHbx/6UQwM9eKgwo0Y67MCPW9wBjE9bRFWBfaeLszEVu3nKyOKLwcUGDXrbmBS8bj9YtyqTG7RjKZmIGuJtPKehEjlNn2ALYYMHUA5VSdVP4LrBiLVLQE5tTDedAzr3uHuWOgLaMeDnfwQg8Hj6QKBgQDGaCoZ3JRg4PbGetduHpGfWOstwZKNwqkDeYvpA03GcRYhzhqG98qdTGEYaLGkszQrh0ZcYk7fs8jsR5Lu2WTzImtiAYDCDEjSDHj5N5PfWR0r2pxrUvX8shun4QKpz4QQ0RUjeujZ2hJHkeCviF2+k5IBDvz6YtBo+H3IgR9QjwKBgQC47IVtuvlYZ0/2TAcAt40YhEgLZOr3NuR4eSnx76zf/8vRmHDEfIqvUdIzzJy3RVYT3uXiB7DYwaHr76ouP3lMOhgLgUlK4Dt8L3UbMP6Asr+6D/uggVmlIHKK76HZIdBL5nGQvTOvwE+fmEivlR5QV+cczqQdCYyNZXETmgAkVQKBgF32JMIcqZR71cLHmFDJX1Okq7P+sWY7YwmHPZA7hVDOa5nU3tE+dpEqA+2oX0DNsY5PwS2tTQc6QJRNjTNadymCCnLenVjIso/vYjc8b+ZdcKg9HsjhACgNPXWy5S0AXt4L9sPXyICreu60EkFvBl5jysh/jaUSuPqNfBxBsk/XAoGBAKYxkzzp4/wKZXfSHh0L2VemUuVCnlTtVWncYtEXeQObbX8CBJ7h2vXzj/mTs2iWfOTA11NLXCmB5FcZfpWv4ACc2U1FtSwA2BUkxZdZcfESNHMwuBEpDvrzbV3mPUvaMsxz366YC+Kw8B5biz+ZwbOtPHzMTfv2wAW3nGdkaSo9AoGBALW8kQ6QwXu3gB7q8Es1zCvSg07tnJks57nJDTtOOg6oc52RvKjndWHT6aFwp/YrESeHNvoudRhPm5nKUTbKDiLixDXcrvEU3b70dJGwWHkBdYU/s5J2o5oaclwiyiUptnIp2eA+Cu+wH1eREGmrOAdeBSFOfdjcfJTFmOCoalJR" // 必须，上一步中使用 RSA签名验签工具 生成的私钥
	var client, err = alipay.New("2021001186696588", privateKey, true)
	client.LoadAppPublicCertFromFile("crt/appCertPublicKey_2021001186696588.crt") // 加载应用公钥证书
	client.LoadAliPayRootCertFromFile("crt/alipayRootCert.crt")                   // 加载支付宝根证书
	client.LoadAliPayPublicCertFromFile("crt/alipayCertPublicKey_RSA2.crt")       // 加载支付宝公钥证书

	// 将 key 的验证调整到初始化阶段
	if err != nil {
		fmt.Println(err)
		return
	}

	// 1、按订单 id 查订单，并校验归属（防止拿别人的订单去付款）
	orderId, err := models.Int(c.Query("id"))
	if err != nil {
		c.String(200, "订单参数错误")
		return
	}
	user := models.User{}
	models.Cookie.Get(c, "userinfo", &user)

	order := models.Order{}
	if err := models.DB.Where("id = ? AND uid = ?", orderId, user.Id).First(&order).Error; err != nil {
		c.String(200, "订单不存在")
		return
	}
	if order.PayStatus == 1 {
		c.Redirect(302, "/user/order")
		return
	}

	// 2、金额和订单号都取自数据库，不能再写死
	var p = alipay.TradePagePay{}
	p.NotifyURL = "http://118.123.14.36:8005/alipayNotify" // TODO: 换成配置项，多环境要能切换
	p.ReturnURL = "http://118.123.14.36:8005/alipayReturn" // TODO: 同上
	p.Subject = "小米商城订单-" + order.OrderId
	p.OutTradeNo = order.OrderId // ★ 必须用我们自己的订单号，回调才找得到订单
	p.TotalAmount = strconv.FormatFloat(order.AllPrice, 'f', 2, 64)
	p.ProductCode = "FAST_INSTANT_TRADE_PAY"

	var url, err4 = client.TradePagePay(p)
	if err4 != nil {
		fmt.Println(err4)
	}

	var payURL = url.String()
	fmt.Println(payURL)

	c.Redirect(302, payURL)

}
func (con AlipayController) AlipayNotify(c *gin.Context) {
	fmt.Println("AlipayNotify")

	var privateKey = "MIIEvgIBADANBgkqhkiG9w0BAQEFAASCBKgwggSkAgEAAoIBAQCPUi2yj4RJcqAKfjV/AxjqIjjUq007pGhLClbrwEoCVhgnQU9bAVFnPlDaVdO2xBfu8D/gHCwV9czZKJh51yg1kwE6Yv/hSd966PnGS9NszfxxWfTbWeC6DIHZj66nTECK+vYWX36tKxIG+juzXyfoAuyL1h58oFgF/Zwa1FzpKDTcMUIE+npf3DpMS0Uatf6TsREzVQQs3i0WVxYY3lv1Dmualr6Q3GAFE1j/xt/STE993uh8MVeLS+RcTrrjilPSVH5Z0DLAjkDSH7XUK1lIHASpgOrddEJ8MeT8L7bUJ8nDs5qQ8zfbPVJaXfsF2NYS162HkIW2bl0r5c5Mudt7AgMBAAECggEAWUHVweHNgiyH7WECkhJsvswHVrNEi0NtzGYpEfOUY/YYXsI22Lduaf0OP5u6GZXwTdeEAF+rORX2uLumkiLkINFnr2QedcEbFCHqBIwOpTF36WQbsUw9P8EwUT1BiWFcxPFctzxL2S78sCnBaol1gfHoPYJhRD5b84cpZDAjmPSJk1XAtAtKChUIskLBAsCvwlGHbx/6UQwM9eKgwo0Y67MCPW9wBjE9bRFWBfaeLszEVu3nKyOKLwcUGDXrbmBS8bj9YtyqTG7RjKZmIGuJtPKehEjlNn2ALYYMHUA5VSdVP4LrBiLVLQE5tTDedAzr3uHuWOgLaMeDnfwQg8Hj6QKBgQDGaCoZ3JRg4PbGetduHpGfWOstwZKNwqkDeYvpA03GcRYhzhqG98qdTGEYaLGkszQrh0ZcYk7fs8jsR5Lu2WTzImtiAYDCDEjSDHj5N5PfWR0r2pxrUvX8shun4QKpz4QQ0RUjeujZ2hJHkeCviF2+k5IBDvz6YtBo+H3IgR9QjwKBgQC47IVtuvlYZ0/2TAcAt40YhEgLZOr3NuR4eSnx76zf/8vRmHDEfIqvUdIzzJy3RVYT3uXiB7DYwaHr76ouP3lMOhgLgUlK4Dt8L3UbMP6Asr+6D/uggVmlIHKK76HZIdBL5nGQvTOvwE+fmEivlR5QV+cczqQdCYyNZXETmgAkVQKBgF32JMIcqZR71cLHmFDJX1Okq7P+sWY7YwmHPZA7hVDOa5nU3tE+dpEqA+2oX0DNsY5PwS2tTQc6QJRNjTNadymCCnLenVjIso/vYjc8b+ZdcKg9HsjhACgNPXWy5S0AXt4L9sPXyICreu60EkFvBl5jysh/jaUSuPqNfBxBsk/XAoGBAKYxkzzp4/wKZXfSHh0L2VemUuVCnlTtVWncYtEXeQObbX8CBJ7h2vXzj/mTs2iWfOTA11NLXCmB5FcZfpWv4ACc2U1FtSwA2BUkxZdZcfESNHMwuBEpDvrzbV3mPUvaMsxz366YC+Kw8B5biz+ZwbOtPHzMTfv2wAW3nGdkaSo9AoGBALW8kQ6QwXu3gB7q8Es1zCvSg07tnJks57nJDTtOOg6oc52RvKjndWHT6aFwp/YrESeHNvoudRhPm5nKUTbKDiLixDXcrvEU3b70dJGwWHkBdYU/s5J2o5oaclwiyiUptnIp2eA+Cu+wH1eREGmrOAdeBSFOfdjcfJTFmOCoalJR" // 必须，上一步中使用 RSA签名验签工具 生成的私钥
	var client, err = alipay.New("2021001186696588", privateKey, true)
	client.LoadAppPublicCertFromFile("crt/appCertPublicKey_2021001186696588.crt") // 加载应用公钥证书
	client.LoadAliPayRootCertFromFile("crt/alipayRootCert.crt")                   // 加载支付宝根证书
	client.LoadAliPayPublicCertFromFile("crt/alipayCertPublicKey_RSA2.crt")       // 加载支付宝公钥证书

	if err != nil {
		fmt.Println(err)
		return
	}

	req := c.Request
	req.ParseForm()

	// 1、验签：防止有人伪造一条"支付成功"的通知
	ok, err := client.VerifySign(req.Form)
	if err != nil || !ok {
		fmt.Println("支付宝回调验签失败:", err)
		c.String(200, "fail")
		return
	}

	// 2、只处理交易成功的通知，其它状态也要 ACK，否则支付宝会一直重试
	if req.Form.Get("trade_status") != "TRADE_SUCCESS" {
		c.String(200, "success")
		return
	}

	outTradeNo := req.Form.Get("out_trade_no")
	payAmount, err := strconv.ParseFloat(req.Form.Get("total_amount"), 64)
	if err != nil {
		fmt.Println("支付宝回调金额解析失败:", err)
		c.String(200, "fail")
		return
	}

	// 3、按我们自己的订单号找回订单（之前 OutTradeNo 用的是当前时间，根本对不上）
	order := models.Order{}
	if err := models.DB.Where("order_id = ?", outTradeNo).First(&order).Error; err != nil {
		fmt.Println("支付宝回调订单不存在:", outTradeNo)
		c.String(200, "fail")
		return
	}

	// 4、金额必须一致，防止"改价后支付"
	if math.Abs(order.AllPrice-payAmount) > 0.01 {
		fmt.Printf("支付宝回调金额不一致 order=%s local=%.2f alipay=%.2f\n", outTradeNo, order.AllPrice, payAmount)
		c.String(200, "fail")
		return
	}

	// 5、幂等更新：pay_status=0 写进 where，重复回调时影响行数为 0，不会重复处理
	res := models.DB.Exec("UPDATE `order` SET pay_status = 1, order_status = 1, pay_time = ? WHERE id = ? AND pay_status = 0",
		models.GetUnix(), order.Id)
	if res.Error != nil {
		fmt.Println("支付宝回调更新订单失败:", res.Error)
		c.String(200, "fail")
		return
	}

	// 6、必须返回 success，否则支付宝会按策略反复重试
	c.String(200, "success")
}
func (con AlipayController) AlipayReturn(c *gin.Context) {
	c.String(200, "支付成功")
}
