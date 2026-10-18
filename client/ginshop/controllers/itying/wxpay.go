package itying

import (
	"fmt"
	"math"
	"net/http"
	"strconv"

	"ginshop/models"

	"github.com/gin-gonic/gin"
	"github.com/objcoding/wxpay"
	qrcode "github.com/skip2/go-qrcode"
)

// 微信回调必须是标准 XML 应答，返回不对微信会一直重试
const (
	wxpaySuccessXML = `<xml><return_code><![CDATA[SUCCESS]]></return_code><return_msg><![CDATA[OK]]></return_msg></xml>`
	wxpayFailXML    = `<xml><return_code><![CDATA[FAIL]]></return_code><return_msg><![CDATA[ERROR]]></return_msg></xml>`
)

type WxpayController struct{}

func (con WxpayController) Wxpay(c *gin.Context) {

	//1、获取订单号 判断此订单号是否值当前用户的
	//2、获取订单里面的支付信息

	//1、配置基本信息
	account := wxpay.NewAccount(
		"wx7bf3787c783116e4",
		"1502539541",
		"zhongyuantengitying6666666666666",
		false,
	)
	client := wxpay.NewClient(account)

	//2、按订单 id 查订单并校验归属，金额和订单号都取自数据库
	orderId, err := models.Int(c.Query("id"))
	if err != nil {
		c.String(http.StatusOK, "订单参数错误")
		return
	}
	user := models.User{}
	models.Cookie.Get(c, "userinfo", &user)

	order := models.Order{}
	if err := models.DB.Where("id = ? AND uid = ?", orderId, user.Id).First(&order).Error; err != nil {
		c.String(http.StatusOK, "订单不存在")
		return
	}

	//3、调用统一下单：金额单位是"分"，不能写死 1 分
	ip := c.ClientIP()
	params := make(wxpay.Params)
	params.SetString("body", "小米商城订单-"+order.OrderId).
		SetString("out_trade_no", order.OrderId). // ★ 用我们自己的订单号
		SetInt64("total_fee", int64(math.Round(order.AllPrice*100))). // 元 → 分
		SetString("spbill_create_ip", ip).
		SetString("notify_url", "http://pay.apiying.com/wxpay/notify"). //必须在商户平台的Native支付回调链接里面配置
		// SetString("trade_type", "APP")
		SetString("trade_type", "NATIVE") //网站支付需要改为NATIVE

	p, err1 := client.UnifiedOrder(params)
	if err1 != nil {
		c.String(http.StatusOK, "生成二维码失败")
		return
	}
	//4、获取code_url 生成支付二维码

	png, err := qrcode.Encode(p["code_url"], qrcode.Medium, 256)
	if err != nil {
		c.String(http.StatusOK, "生成二维码失败")
		return
	}
	c.String(http.StatusOK, string(png))
}

/*
1、发布到服务器测试
2、必须在商户平台的Native支付回调链接里面配置
3、如何接收XML的数据  c.GetRawData()
4、如何获取数据
5、如何验证数据
6、更新数据

*/
func (con WxpayController) WxpayNotify(c *gin.Context) {
	//1、获取表单传过来的xml数据
	xmlByte, _ := c.GetRawData()
	xmlStr := string(xmlByte)

	postParams := wxpay.XmlToMap(xmlStr)

	//2、校验签名
	account := wxpay.NewAccount(
		"wx7bf3787c783116e4",
		"1502539541",
		"zhongyuantengitying6666666666666",
		false,
	)
	client := wxpay.NewClient(account)
	isValidate := client.ValidSign(postParams)

	//3、验签失败直接返回 FAIL，防止伪造回调
	if !isValidate {
		fmt.Println("微信回调验签失败")
		c.String(200, wxpayFailXML)
		return
	}
	fmt.Println("微信回调参数:", postParams)

	//4、只处理支付成功的通知
	if postParams["return_code"] != "SUCCESS" || postParams["result_code"] != "SUCCESS" {
		c.String(200, wxpaySuccessXML)
		return
	}

	outTradeNo := postParams["out_trade_no"]
	totalFee, err := strconv.ParseInt(postParams["total_fee"], 10, 64) // 单位：分
	if err != nil {
		fmt.Println("微信回调金额解析失败:", err)
		c.String(200, wxpayFailXML)
		return
	}

	order := models.Order{}
	if err := models.DB.Where("order_id = ?", outTradeNo).First(&order).Error; err != nil {
		fmt.Println("微信回调订单不存在:", outTradeNo)
		c.String(200, wxpayFailXML)
		return
	}

	//5、金额校验：本地存的是元，回调给的是分
	if int64(math.Round(order.AllPrice*100)) != totalFee {
		fmt.Printf("微信回调金额不一致 order=%s local=%d wx=%d\n",
			outTradeNo, int64(math.Round(order.AllPrice*100)), totalFee)
		c.String(200, wxpayFailXML)
		return
	}

	//6、幂等更新：只有未支付的订单才会被改
	res := models.DB.Exec("UPDATE `order` SET pay_status = 1, order_status = 1, pay_time = ? WHERE id = ? AND pay_status = 0",
		models.GetUnix(), order.Id)
	if res.Error != nil {
		fmt.Println("微信回调更新订单失败:", res.Error)
		c.String(200, wxpayFailXML)
		return
	}

	//7、必须返回 SUCCESS
	c.String(200, wxpaySuccessXML)
}
