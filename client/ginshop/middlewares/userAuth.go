package middlewares

import (
	"ginshop/models"

	"github.com/gin-gonic/gin"
)

func InitUserAuthMiddleware(c *gin.Context) {
	//判断用户有没有登录
	user := models.User{}
	isLogin := models.Cookie.Get(c, "userinfo", &user)
	if !isLogin || len(user.Phone) != 11 {
		c.Redirect(302, "/pass/login")
		// ★ 必须 Abort：Gin 的 handler 链由 c.Next() 驱动，
		//   中间件里只 return 只结束当前函数，后面的业务 handler 仍然会执行。
		//   不加这一行，未登录 POST /buy/doCheckout 照样会走下单逻辑（user.Id = 0）。
		c.Abort()
		return
	}
}
