package admin

import (
	"context"
	"ginshop/models"
	pbRbac "ginshop/proto/rbacRole"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type RoleController struct {
	BaseController
}

func (con RoleController) Index(c *gin.Context) {
	rbacClient := pbRbac.NewRbacRoleService("rbac", models.RbacClient)
	res, _ := rbacClient.RoleGet(context.Background(), &pbRbac.RoleGetRequest{})
	c.HTML(http.StatusOK, "admin/role/index.html", gin.H{
		"roleList": res.RoleList,
	})
}
func (con RoleController) Add(c *gin.Context) {
	c.HTML(http.StatusOK, "admin/role/add.html", gin.H{})
}
func (con RoleController) DoAdd(c *gin.Context) {

	title := strings.Trim(c.PostForm("title"), " ")
	description := strings.Trim(c.PostForm("description"), " ")
	if title == "" {
		con.Error(c, "角色的标题不能为空", "/admin/role/add")
		return
	}

	rbacClient := pbRbac.NewRbacRoleService("rbac", models.RbacClient)
	res, _ := rbacClient.RoleAdd(context.Background(), &pbRbac.RoleAddRequest{
		Title:       title,
		Description: description,
		AddTime:     models.GetUnix(),
		Status:      1,
	})
	if !res.Success {
		con.Error(c, "增加角色失败 请重试", "/admin/role/add")
	} else {
		con.Success(c, "增加角色成功", "/admin/role")
	}

}
func (con RoleController) Edit(c *gin.Context) {

	id, err := models.Int(c.Query("id"))
	if err != nil {
		con.Error(c, "传入数据错误", "/admin/role")
	} else {
		rbacClient := pbRbac.NewRbacRoleService("rbac", models.RbacClient)
		res, _ := rbacClient.RoleGet(context.Background(), &pbRbac.RoleGetRequest{
			Id: int64(id),
		})

		c.HTML(http.StatusOK, "admin/role/edit.html", gin.H{
			"role": res.RoleList[0],
		})
	}

}
func (con RoleController) DoEdit(c *gin.Context) {

	id, err1 := models.Int(c.PostForm("id"))
	if err1 != nil {
		con.Error(c, "传入数据错误", "/admin/role")
		return
	}
	title := strings.Trim(c.PostForm("title"), " ")
	description := strings.Trim(c.PostForm("description"), " ")

	if title == "" {
		con.Error(c, "角色的标题不能为空", "/admin/role/edit")
	}
	//调用微服务修改
	rbacClient := pbRbac.NewRbacRoleService("rbac", models.RbacClient)
	res, _ := rbacClient.RoleEdit(context.Background(), &pbRbac.RoleEditRequest{
		Id:          int64(id),
		Title:       title,
		Description: description,
	})
	if !res.Success {
		con.Error(c, "修改数据失败", "/admin/role/edit?id="+models.String(id))
	} else {
		con.Success(c, "修改数据成功", "/admin/role/edit?id="+models.String(id))
	}
}
func (con RoleController) Delete(c *gin.Context) {
	id, err := models.Int(c.Query("id"))
	if err != nil {
		con.Error(c, "传入数据错误", "/admin/role")
	} else {
		rbacClient := pbRbac.NewRbacRoleService("rbac", models.RbacClient)
		res, _ := rbacClient.RoleDelete(context.Background(), &pbRbac.RoleDeleteRequest{
			Id: int64(id),
		})
		if res.Success {
			con.Success(c, "删除数据成功", "/admin/role")
		} else {
			con.Error(c, "删除数据失败", "/admin/role")
		}
	}
}

func (con RoleController) Auth(c *gin.Context) {
	//1、获取角色id
	roleId, err := models.Int(c.Query("id"))
	if err != nil {
		con.Error(c, "传入数据错误", "/admin/role")
		return
	}
	//调用微服务显示授权页面
	rbacClient := pbRbac.NewRbacRoleService("rbac", models.RbacClient)
	res, _ := rbacClient.RoleAuth(context.Background(), &pbRbac.RoleAuthRequest{
		RoleId: int64(roleId),
	})

	c.HTML(http.StatusOK, "admin/role/auth.html", gin.H{
		"roleId":     roleId,
		"accessList": res.AccessList,
	})

}

func (con RoleController) DoAuth(c *gin.Context) {
	//获取角色id
	roleId, err1 := models.Int(c.PostForm("role_id"))
	if err1 != nil {
		con.Error(c, "传入数据错误", "/admin/role")
		return
	}
	//获取权限id  切片
	accessIds := c.PostFormArray("access_node[]")
	//调用微服务执行授权
	rbacClient := pbRbac.NewRbacRoleService("rbac", models.RbacClient)
	res, _ := rbacClient.RoleDoAuth(context.Background(), &pbRbac.RoleDoAuthRequest{
		RoleId:    int64(roleId),
		AccessIds: accessIds,
	})
	if res.Success {
		con.Success(c, "授权成功", "/admin/role/auth?id="+models.String(roleId))
	} else {
		con.Error(c, "授权失败", "/admin/role/auth?id="+models.String(roleId))
	}
}
