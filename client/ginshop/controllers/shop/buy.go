package itying

import (
	"fmt"
	"ginshop/models"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type BuyController struct {
	BaseController
}

func (con BuyController) Checkout(c *gin.Context) {
	//1、获取购物车中选择的商品

	cartList := []models.Cart{}
	models.Cookie.Get(c, "cartList", &cartList)

	orderList := []models.Cart{}
	var allPrice float64
	var allNum int

	for i := 0; i < len(cartList); i++ {
		if cartList[i].Checked {
			allPrice += cartList[i].Price * float64(cartList[i].Num)
			orderList = append(orderList, cartList[i])
			allNum += cartList[i].Num
		}
	}

	//2、获取当前用户的收货地址

	user := models.User{}
	models.Cookie.Get(c, "userinfo", &user)
	addressList := []models.Address{}
	models.DB.Where("uid = ?", user.Id).Order("id desc").Find(&addressList)

	//3、生成签名
	orderSign := models.Md5(models.GetRandomNum())
	session := sessions.Default(c)
	session.Set("orderSign", orderSign)
	session.Save()

	//4、判断orderList数据是否存在
	if len(orderList) == 0 {
		c.Redirect(302, "/")
		return
	}

	con.Render(c, "itying/buy/checkout.html", gin.H{
		"orderList":   orderList,
		"allPrice":    allPrice,
		"allNum":      allNum,
		"addressList": addressList,
		"orderSign":   orderSign,
	})

}

/*
提交订单执行结算
   1、获取用户信息 获取用户的收货地址信息
   2、获取购买商品的信息
   3、把订单信息放在订单表，把商品信息放在商品表
   4、删除购物车里面的选中数据
   5、跳转到支付页面
*/
func (con BuyController) DoCheckout(c *gin.Context) {
	//0、防止重复提交订单
	orderSignClient := c.PostForm("orderSign")
	session := sessions.Default(c)
	orderSignSession := session.Get("orderSign")
	orderSignServer, ok := orderSignSession.(string)
	if !ok || orderSignServer == "" || orderSignClient != orderSignServer {
		c.Redirect(302, "/")
		c.Abort()
		return
	}
	session.Delete("orderSign")
	session.Save()

	// 1、获取用户信息 获取用户的收货地址信息
	user := models.User{}
	models.Cookie.Get(c, "userinfo", &user)

	addressResult := []models.Address{}
	models.DB.Where("uid = ? AND default_address=1", user.Id).Find(&addressResult)
	if len(addressResult) == 0 {
		c.Redirect(302, "/buy/checkout")
		return
	}

	// 2、取购物车里勾选的商品
	//    ★ 只信任 goodsId 和 num，价格、标题一律回库里查（Cookie 里的价格可以被改）
	cartList := []models.Cart{}
	models.Cookie.Get(c, "cartList", &cartList)

	type orderLine struct {
		goods models.Goods
		num   int
		color string
	}
	lines := make([]orderLine, 0, len(cartList))
	var allPrice float64

	for i := 0; i < len(cartList); i++ {
		if !cartList[i].Checked || cartList[i].Num <= 0 {
			continue
		}
		goods := models.Goods{}
		if err := models.DB.Where("id = ? AND is_delete = 0", cartList[i].Id).First(&goods).Error; err != nil {
			// 商品已下架或不存在，回购物车让用户重新确认
			c.Redirect(302, "/cart")
			return
		}
		allPrice += goods.Price * float64(cartList[i].Num) // ★ 用数据库里的价格
		lines = append(lines, orderLine{goods: goods, num: cartList[i].Num, color: cartList[i].GoodsColor})
	}
	if len(lines) == 0 {
		c.Redirect(302, "/cart")
		return
	}

	// 3、事务：建订单主表 → 逐条建明细 → 扣库存，任何一步失败整体回滚
	var order models.Order
	err := models.DB.Transaction(func(tx *gorm.DB) error {
		order = models.Order{
			OrderId:     models.GetOrderId(),
			Uid:         user.Id,
			AllPrice:    allPrice,
			Phone:       addressResult[0].Phone,
			Name:        addressResult[0].Name,
			Address:     addressResult[0].Address,
			PayStatus:   0,
			PayType:     0,
			OrderStatus: 0,
			AddTime:     int(models.GetUnix()),
		}
		if err := tx.Create(&order).Error; err != nil {
			return err
		}

		items := make([]models.OrderItem, 0, len(lines))
		for _, line := range lines {
			// ★ 乐观锁扣库存：where 里带 goods_number >= n，靠 RowsAffected 判断是否扣成功
			res := tx.Exec("UPDATE goods SET goods_number = goods_number - ? WHERE id = ? AND goods_number >= ?",
				line.num, line.goods.Id, line.num)
			if res.Error != nil {
				return res.Error
			}
			if res.RowsAffected == 0 {
				return fmt.Errorf("商品[%s]库存不足", line.goods.Title)
			}

			items = append(items, models.OrderItem{
				OrderId:      order.Id,
				Uid:          user.Id,
				ProductTitle: line.goods.Title,
				ProductId:    line.goods.Id,
				ProductImg:   line.goods.GoodsImg,
				ProductPrice: line.goods.Price, // 成交价快照，来自数据库
				ProductNum:   line.num,
				GoodsVersion: line.goods.GoodsVersion,
				GoodsColor:   line.color,
			})
		}
		return tx.CreateInBatches(items, 100).Error
	})
	if err != nil {
		fmt.Println("下单失败:", err)
		c.Redirect(302, "/buy/checkout")
		return
	}

	// 4、事务成功后才清理购物车里已下单的项（按商品 ID 匹配）
	orderedGoods := make(map[int]bool, len(lines))
	for _, line := range lines {
		orderedGoods[line.goods.Id] = true
	}
	noSelectCartList := make([]models.Cart, 0, len(cartList))
	for i := 0; i < len(cartList); i++ {
		if cartList[i].Checked && orderedGoods[cartList[i].Id] {
			continue
		}
		noSelectCartList = append(noSelectCartList, cartList[i])
	}
	models.Cookie.Set(c, "cartList", noSelectCartList)

	c.Redirect(302, "/buy/pay?orderId="+models.String(order.Id))
}

//支付
func (con BuyController) Pay(c *gin.Context) {

	orderId, err := models.Int(c.Query("orderId"))
	if err != nil {
		c.Redirect(302, "/")
	}
	//获取用户信息
	user := models.User{}
	models.Cookie.Get(c, "userinfo", &user)
	//获取订单信息
	order := models.Order{}
	models.DB.Where("id = ?", orderId).Find(&order)
	if order.Uid != user.Id {
		c.Redirect(302, "/")
		return
	}
	//获取订单对应的商品

	orderItems := []models.OrderItem{}
	models.DB.Where("order_id = ?", orderId).Find(&orderItems)

	con.Render(c, "itying/buy/pay.html", gin.H{
		"order":      order,
		"orderItems": orderItems,
	})
}

func (con BuyController) OrderPayStatus(c *gin.Context) {

	id, err := models.Int(c.Query("id"))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "传入参数错误",
		})
		return
	}
	//获取用户信息
	user := models.User{}
	models.Cookie.Get(c, "userinfo", &user)

	//获取主订单信息
	order := models.Order{}
	models.DB.Where("id=?", id).Find(&order)

	//判断当前数据是否合法
	if user.Id != order.Uid {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "非法请求",
		})
		return
	}

	//判断是否支付
	if order.PayStatus == 1 && order.OrderStatus == 1 {
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "支付成功",
		})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "支付成功",
		})
		return
	}

}
