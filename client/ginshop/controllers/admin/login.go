package admin

import (
	"context"
	"encoding/json"
	"fmt"
	"ginshop/models"
	pbRbac "ginshop/proto/rbacLogin"
	"net/http"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type LoginController struct {
	BaseController
}

func (con LoginController) Index(c *gin.Context) {
	//验证md5是否正确
	// fmt.Println(models.Md5("123456"))   e10adc3949ba59abbe56e057f20f883e

	c.HTML(http.StatusOK, "admin/login/login.html", gin.H{})

}

func (con LoginController) DoLogin(c *gin.Context) {

	captchaId := c.PostForm("captchaId")
	username := c.PostForm("username")
	password := c.PostForm("password")
	verifyValue := c.PostForm("verifyValue")
	fmt.Println(username, password)
	//1、验证验证码是否正确
	if flag := models.VerifyCaptcha(captchaId, verifyValue); flag {
		//2、调用rbac的微服务 查询数据库 判断用户以及密码是否存在
		ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
		defer cancel()

		rbacClient := pbRbac.NewRbacLoginService("rbac", models.RbacClient)
		res, err := rbacClient.Login(ctx, &pbRbac.LoginRequest{
			Username: username,
			Password: models.Md5(password),
		})
		// ★ 原来写的是 res, _：rbac 不可用时 res 为 nil，下一行 res.IsLogin 直接空指针 panic
		if err != nil || res == nil {
			fmt.Println("调用权限服务失败:", err)
			con.Error(c, "权限服务不可用，请先启动 Consul 和 rbac 服务", "/admin/login")
			return
		}
		if res.IsLogin {
			//3、执行登录 保存用户信息 执行跳转
			session := sessions.Default(c)
			//注意：session.Set没法直接保存结构体对应的切片 把结构体转换成json字符串
			userinfoSlice, _ := json.Marshal(res.Userlist)
			session.Set("userinfo", string(userinfoSlice))
			session.Save()
			con.Success(c, "登录成功", "/admin")

		} else {
			con.Error(c, "用户名或者密码错误", "/admin/login")
		}

	} else {
		con.Error(c, "验证码验证失败", "/admin/login")
	}

}

func (con LoginController) Captcha(c *gin.Context) {
	// 参数顺序是 (高, 宽, 字符数)。原来是 (34, 100, 2)，只有两位，
	// 36 个候选字符两位只有 1296 种组合，太弱；改成 4 位，宽度同步放大到 120。
	id, b64s, err := models.MakeCaptcha(40, 120, 4)
	if err != nil {
		// ★ 原来只 fmt.Println 一行，然后照常返回空的 captchaId/Image，
		//   前端拿着空值去设 src，验证码区域一片空白，看不出是服务没起。
		fmt.Println("生成验证码失败:", err)
		c.JSON(http.StatusOK, gin.H{
			"success":      false,
			"message":      "验证码服务不可用，请先启动 Consul 和 captcha 服务",
			"captchaId":    "",
			"captchaImage": models.CaptchaErrImg,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success":      true,
		"captchaId":    id,
		"captchaImage": b64s,
	})
}

func (con LoginController) LoginOut(c *gin.Context) {
	session := sessions.Default(c)
	session.Delete("userinfo")
	session.Save()
	con.Success(c, "退出登录成功", "/admin/login")
}
