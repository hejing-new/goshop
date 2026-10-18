package middlewares

import (
	"context"
	"encoding/json"
	"fmt"
	"ginshop/models"
	pbRbac "ginshop/proto/rbacRole"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"gopkg.in/ini.v1"
)

func InitAdminAuthMiddleware(c *gin.Context) {
	// 1、获取访问路径，去掉 query（前端会带 ?t=时间戳，不去掉永远匹配不上 access.url）
	pathname := strings.Split(c.Request.URL.String(), "?")[0]

	// 2、从 session 取管理员信息（session.Set 存不了结构体，存的是 JSON 字符串）
	session := sessions.Default(c)
	userinfoStr, ok := session.Get("userinfo").(string)

	// 3、没有登录：只放行 login / doLogin / captcha
	if !ok || userinfoStr == "" {
		if isAdminLoginPath(pathname) {
			return
		}
		toAdminLogin(c)
		return
	}

	var userinfoStruct []models.Manager
	if err := json.Unmarshal([]byte(userinfoStr), &userinfoStruct); err != nil ||
		len(userinfoStruct) == 0 || userinfoStruct[0].Username == "" {
		if isAdminLoginPath(pathname) {
			return
		}
		toAdminLogin(c)
		return
	}
	manager := userinfoStruct[0]

	// 4、超级管理员直接放行，不查库
	if manager.IsSuper == 1 {
		return
	}

	// 5、/admin/goods → goods，和 access.url 字段保持同一套格式
	urlPath := strings.Replace(pathname, "/admin/", "", 1)
	if excludeAuthPath("/" + urlPath) {
		return
	}

	// 6、调用 rbac 微服务做权限判断：必须带超时，并且必须判断 error
	ctx, cancel := context.WithTimeout(c.Request.Context(), 300*time.Millisecond)
	defer cancel()

	rbacClient := pbRbac.NewRbacRoleService("rbac", models.RbacClient)
	res, err := rbacClient.MiddlewaresAuth(ctx, &pbRbac.MiddlewaresAuthRequest{
		RoleId:  int64(manager.RoleId),
		UrlPath: urlPath,
	})
	if err != nil || res == nil {
		// 安全组件失败必须"拒绝"，不能放行（fail-close），否则 rbac 一挂就等于全网越权
		fmt.Printf("调用权限服务失败: %v\n", err)
		c.String(http.StatusServiceUnavailable, "权限服务不可用，请稍后重试")
		c.Abort()
		return
	}
	if !res.HasPermission {
		c.String(http.StatusForbidden, "没有权限")
		c.Abort()
		return
	}
	// 没被拦截 = 放行，链会继续走到业务 handler
}

func isAdminLoginPath(pathname string) bool {
	return pathname == "/admin/login" || pathname == "/admin/doLogin" || pathname == "/admin/captcha"
}

// 跳登录页：Redirect 之后必须 Abort，否则后面的业务 handler 仍会执行
func toAdminLogin(c *gin.Context) {
	c.Redirect(302, "/admin/login")
	c.Abort()
}

var (
	excludeAuthOnce sync.Once
	excludeAuthList []string
)

// 排除权限判断的方法
// 改造点：配置文件只在第一次调用时读一次；读不到也只是打印日志，
// 不再 os.Exit(1)——原来一个请求就能把整个进程杀掉。
func excludeAuthPath(urlPath string) bool {
	excludeAuthOnce.Do(func() {
		config, err := ini.Load("./conf/app.ini")
		if err != nil {
			fmt.Printf("读取配置文件失败: %v\n", err)
			return
		}
		raw := config.Section("").Key("excludeAuthPath").String()
		for _, v := range strings.Split(raw, ",") {
			if v = strings.TrimSpace(v); v != "" {
				excludeAuthList = append(excludeAuthList, v)
			}
		}
	})
	for _, v := range excludeAuthList {
		if v == urlPath {
			return true
		}
	}
	return false
}
