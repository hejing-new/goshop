package admin

import (
	"context"
	"ginshop/models"
	pbManager "ginshop/proto/rbacManager"
	pbRbac "ginshop/proto/rbacRole"
	pbRole "ginshop/proto/rbacRole"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type ManagerController struct {
	BaseController
}

func (con ManagerController) Index(c *gin.Context) {

	rbacClient := pbManager.NewRbacManagerService("rbac", models.RbacClient)
	res, _ := rbacClient.ManagerGet(context.Background(), &pbManager.ManagerGetRequest{})

	c.HTML(http.StatusOK, "admin/manager/index.html", gin.H{
		"managerList": res.ManagerList,
	})

}
func (con ManagerController) Add(c *gin.Context) {
	//获取所有的角色
	rbacClient := pbRole.NewRbacRoleService("rbac", models.RbacClient)
	res, _ := rbacClient.RoleGet(context.Background(), &pbRbac.RoleGetRequest{})
	c.HTML(http.StatusOK, "admin/manager/add.html", gin.H{
		"roleList": res.RoleList,
	})
}
func (con ManagerController) DoAdd(c *gin.Context) {

	roleId, err1 := models.Int(c.PostForm("role_id"))
	if err1 != nil {
		con.Error(c, "传入数据错误", "/admin/manager/add")
		return
	}
	username := strings.Trim(c.PostForm("username"), " ")
	password := strings.Trim(c.PostForm("password"), " ")
	email := strings.Trim(c.PostForm("email"), " ")
	mobile := strings.Trim(c.PostForm("mobile"), " ")
	//用户名和密码长度是否合法
	if len(username) < 2 || len(password) < 6 {
		con.Error(c, "用户名或者密码的长度不合法", "/admin/manager/add")
		return
	}
	//判断管理是否存在
	rbacClient := pbManager.NewRbacManagerService("rbac", models.RbacClient)
	res, _ := rbacClient.ManagerGet(context.Background(), &pbManager.ManagerGetRequest{
		Username: username,
	})

	if len(res.ManagerList) > 0 {
		con.Error(c, "此管理员已存在", "/admin/manager/add")
		return
	}
	//执行增加管理员
	addResult, _ := rbacClient.ManagerAdd(context.Background(), &pbManager.ManagerAddRequest{
		Username: username,
		Password: models.Md5(password),
		Email:    email,
		Mobile:   mobile,
		RoleId:   int64(roleId),
		Status:   1,
		AddTime:  int64(models.GetUnix()),
	})

	if !addResult.Success {
		con.Error(c, "增加管理员失败", "/admin/manager/add")
		return
	}
	con.Success(c, "增加管理员成功", "/admin/manager")
}

func (con ManagerController) Edit(c *gin.Context) {
	//获取管理员
	id, err := models.Int(c.Query("id"))
	if err != nil {
		con.Error(c, "传入数据错误", "/admin/manager")
		return
	}

	//获取管理员
	managerClient := pbManager.NewRbacManagerService("rbac", models.RbacClient)
	resManager, _ := managerClient.ManagerGet(context.Background(), &pbManager.ManagerGetRequest{
		Id: int64(id),
	})

	//获取所有的角色
	roleClient := pbRole.NewRbacRoleService("rbac", models.RbacClient)
	resRole, _ := roleClient.RoleGet(context.Background(), &pbRbac.RoleGetRequest{})

	c.HTML(http.StatusOK, "admin/manager/edit.html", gin.H{
		"manager":  resManager.ManagerList[0],
		"roleList": resRole.RoleList,
	})
}

func (con ManagerController) DoEdit(c *gin.Context) {
	id, err1 := models.Int(c.PostForm("id"))
	if err1 != nil {
		con.Error(c, "传入数据错误", "/admin/manager")
		return
	}
	roleId, err2 := models.Int(c.PostForm("role_id"))
	if err2 != nil {
		con.Error(c, "传入数据错误", "/admin/manager")
		return
	}
	username := strings.Trim(c.PostForm("username"), " ")
	password := strings.Trim(c.PostForm("password"), " ")
	email := strings.Trim(c.PostForm("email"), " ")
	mobile := strings.Trim(c.PostForm("mobile"), " ")

	if len(mobile) > 11 {
		con.Error(c, "mobile长度不合法", "/admin/manager/edit?id="+models.String(id))
		return
	}
	//注意：判断密码是否为空 为空表示不修改密码 不为空表示修改密码
	if password != "" {
		//判断密码长度是否合法
		if len(password) < 6 {
			con.Error(c, "密码的长度不合法 密码长度不能小于6位", "/admin/manager/edit?id="+models.String(id))
			return
		}
		password = models.Md5(password)
	}

	managerClient := pbManager.NewRbacManagerService("rbac", models.RbacClient)
	editResult, _ := managerClient.ManagerEdit(context.Background(), &pbManager.ManagerEditRequest{
		Id:       int64(id),
		Username: username,
		Password: password,
		Email:    email,
		Mobile:   mobile,
		RoleId:   int64(roleId),
	})
	if !editResult.Success {
		con.Error(c, "修改数据失败", "/admin/manager/edit?id="+models.String(id))
		return
	}
	con.Success(c, "修改数据成功", "/admin/manager")
}

func (con ManagerController) Delete(c *gin.Context) {
	id, err := models.Int(c.Query("id"))
	if err != nil {
		con.Error(c, "传入数据错误", "/admin/manager")
	} else {
		managerClient := pbManager.NewRbacManagerService("rbac", models.RbacClient)
		managerRes, _ := managerClient.ManagerDelete(context.Background(), &pbManager.ManagerDeleteRequest{
			Id: int64(id),
		})
		if managerRes.Success {
			con.Success(c, "删除数据成功", "/admin/manager")
			return
		}
		con.Success(c, "删除数据失败", "/admin/manager")
	}
}
