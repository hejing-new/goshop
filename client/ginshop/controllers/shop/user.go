package itying

import (
	"ginshop/models"
	"math"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserController struct {
	BaseController
}

func (con UserController) Index(c *gin.Context) {
	// c.String(http.StatusOK, "首页")
	var tpl = "itying/user/welcome.html"
	con.Render(c, tpl, gin.H{})
}
func (con UserController) OrderList(c *gin.Context) {
	// 当前页
	page, _ := models.Int(c.Query("page"))
	if page == 0 {
		page = 1
	}
	pageSize := 2

	// 当前用户
	user := models.User{}
	models.Cookie.Get(c, "userinfo", &user)

	// 按订单状态筛选，非法值当作"全部"
	orderStatus, statusErr := models.Int(c.Query("orderStatus"))
	if statusErr != nil || orderStatus < 0 {
		orderStatus = -1
	}

	// 按商品标题模糊搜索：先查出当前用户命中的订单 id，再用 IN 过滤
	keywords := c.Query("keywords")
	var matchOrderIds []int
	if keywords != "" {
		// ★ 带上 uid：只搜自己的订单明细，不会把别人的订单号带进来
		models.DB.Model(&models.OrderItem{}).
			Where("product_title LIKE ?", "%"+keywords+"%").
			Where("uid = ?", user.Id).
			Pluck("order_id", &matchOrderIds)

		// ★ 一条都没搜到就直接返回空列表。
		//   原来的写法会拼出 "id in ()"，MySQL 直接报语法错误，整页 500
		if len(matchOrderIds) == 0 {
			con.Render(c, "itying/user/order.html", gin.H{
				"order":       []models.Order{},
				"page":        page,
				"keywords":    keywords,
				"orderStatus": orderStatus,
				"totalPages":  float64(0),
			})
			return
		}
	}

	// 条件全部用占位符，不做字符串拼接
	buildQuery := func() *gorm.DB {
		q := models.DB.Model(&models.Order{}).Where("uid = ?", user.Id)
		if len(matchOrderIds) > 0 {
			q = q.Where("id IN ?", matchOrderIds)
		}
		if orderStatus >= 0 {
			q = q.Where("order_status = ?", orderStatus)
		}
		return q
	}

	// 获取当前用户下面的订单
	orderList := []models.Order{}
	buildQuery().Offset((page - 1) * pageSize).Limit(pageSize).
		Preload("OrderItem").Order("add_time desc").Find(&orderList)

	// 获取总数量
	var count int64
	buildQuery().Count(&count)

	con.Render(c, "itying/user/order.html", gin.H{
		"order":       orderList,
		"page":        page,
		"keywords":    keywords,
		"orderStatus": orderStatus,
		"totalPages":  math.Ceil(float64(count) / float64(pageSize)),
	})
}
func (con UserController) OrderInfo(c *gin.Context) {

	id, err := models.Int(c.Query("id"))
	if err != nil {
		c.Redirect(302, "/user/order")
		return // ★ 原来少了 return，参数错误时还会带着 id=0 继续往下查
	}
	user := models.User{}
	models.Cookie.Get(c, "userinfo", &user)
	order := []models.Order{}
	models.DB.Where("id=? And uid=?", id, user.Id).Preload("OrderItem").Find(&order)

	if len(order) == 0 {
		c.Redirect(302, "/user/order")
		return
	}
	var tpl = "itying/user/order_info.html"
	con.Render(c, tpl, gin.H{
		"order": order[0],
	})
}
