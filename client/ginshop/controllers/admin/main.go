package admin

import (
	"context"
	"encoding/json"
	"fmt"
	"ginshop/models"
	pbRbac "ginshop/proto/rbacRole"
	"net/http"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type MainController struct {
	BaseController
}

func (con MainController) Index(c *gin.Context) {
	//获取userinfo 对应的session
	session := sessions.Default(c)
	userinfo := session.Get("userinfo")
	//类型断言 来判断 userinfo是不是一个string
	userinfoStr, ok := userinfo.(string)

	if ok {
		//1、获取用户信息
		var userinfoStruct []models.Manager
		json.Unmarshal([]byte(userinfoStr), &userinfoStruct)
		//调用微服务获取当前角色拥有权限（菜单）
		ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
		defer cancel()

		rbacClient := pbRbac.NewRbacRoleService("rbac", models.RbacClient)
		res, err := rbacClient.RoleAuth(ctx, &pbRbac.RoleAuthRequest{
			RoleId: int64(userinfoStruct[0].RoleId),
		})
		// ★ 原来 res, _ 吞掉 error：rbac 不可用时 res 为 nil，下一行取 res.AccessList 直接 panic
		if err != nil || res == nil {
			fmt.Println("获取权限菜单失败:", err)
			c.String(http.StatusServiceUnavailable, "权限服务不可用，请先启动 Consul 和 rbac 服务")
			return
		}

		// fmt.Printf("%#v", accessList)
		c.HTML(http.StatusOK, "admin/main/index.html", gin.H{
			"username":   userinfoStruct[0].Username,
			"accessList": res.AccessList,
			"isSuper":    userinfoStruct[0].IsSuper,
		})
	} else {
		c.Redirect(302, "/admin/login")
	}

}

func (con MainController) Welcome(c *gin.Context) {
	c.HTML(http.StatusOK, "admin/main/welcome.html", gin.H{})
}

// 公共接口允许修改的表和字段白名单。
// 表名、字段名来自前端的 data-table / data-field，如果直接拼进 SQL 就是注入点。
var allowChangeFields = map[string]map[string]bool{
	"goods":                {"status": true, "is_hot": true, "is_best": true, "is_new": true, "sort": true, "goods_number": true},
	"goods_cate":           {"status": true, "sort": true},
	"goods_type":           {"status": true},
	"goods_type_attribute": {"status": true, "sort": true},
	"access":               {"sort": true},
	"role":                 {"status": true},
	"nav":                  {"status": true, "sort": true},
	"focus":                {"status": true, "sort": true},
}

func checkChangeParam(table, field string) bool {
	fields, ok := allowChangeFields[table]
	if !ok {
		return false
	}
	return fields[field]
}

//公共修改状态的方法：把 0/1 字段取反
func (con MainController) ChangeStatus(c *gin.Context) {
	id, err := models.Int(c.Query("id"))
	table := c.Query("table")
	field := c.Query("field")
	if err != nil || !checkChangeParam(table, field) {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "传入的参数错误",
		})
		return
	}

	// 表名和字段名已经过白名单校验，这里加反引号拼接是安全的；id 仍然走占位符
	err1 := models.DB.Exec("update `"+table+"` set `"+field+"`=ABS(`"+field+"`-1) where id=?", id).Error
	if err1 != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "修改失败 请重试",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "修改成功",
	})
}

//公共修改数值的方法
func (con MainController) ChangeNum(c *gin.Context) {
	id, err := models.Int(c.Query("id"))
	table := c.Query("table")
	field := c.Query("field")
	num, numErr := models.Int(c.Query("num"))
	if err != nil || numErr != nil || num < 0 || !checkChangeParam(table, field) {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "传入的参数错误",
		})
		return
	}

	// 数值也走占位符，不再拼进 SQL 文本
	err1 := models.DB.Exec("update `"+table+"` set `"+field+"`=? where id=?", num, id).Error
	if err1 != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "修改数据失败",
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "修改成功",
		})
	}
}

//清除缓存
func (con MainController) FlushAll(c *gin.Context) {
	models.CacheDb.FlushAll()
	con.Success(c, "清除Redis缓存数据成功", "/admin")
}
