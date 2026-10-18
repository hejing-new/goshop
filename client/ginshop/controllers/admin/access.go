package admin

import (
	"context"
	"ginshop/models"
	pbRbac "ginshop/proto/rbacAccess"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type AccessController struct {
	BaseController
}

func (con AccessController) Index(c *gin.Context) {

	rbacClient := pbRbac.NewRbacAccessService("rbac", models.RbacClient)
	res, _ := rbacClient.AccessGet(context.Background(), &pbRbac.AccessGetRequest{})

	c.HTML(http.StatusOK, "admin/access/index.html", gin.H{
		"accessList": res.AccessList,
	})

}
func (con AccessController) Add(c *gin.Context) {
	//获取顶级模块
	rbacClient := pbRbac.NewRbacAccessService("rbac", models.RbacClient)
	res, _ := rbacClient.AccessGet(context.Background(), &pbRbac.AccessGetRequest{})

	c.HTML(http.StatusOK, "admin/access/add.html", gin.H{
		"accessList": res.AccessList,
	})
}
func (con AccessController) DoAdd(c *gin.Context) {
	moduleName := strings.Trim(c.PostForm("module_name"), " ")
	actionName := c.PostForm("action_name")
	accessType, err1 := models.Int(c.PostForm("type"))
	url := c.PostForm("url")
	moduleId, err2 := models.Int(c.PostForm("module_id"))
	sort, err3 := models.Int(c.PostForm("sort"))
	status, err4 := models.Int(c.PostForm("status"))
	description := c.PostForm("description")
	if err1 != nil || err2 != nil || err3 != nil || err4 != nil {
		con.Error(c, "传入参数错误", "/admin/access/add")
		return
	}
	if moduleName == "" {
		con.Error(c, "模块名称不能为空", "/admin/access/add")
		return
	}

	rbacClient := pbRbac.NewRbacAccessService("rbac", models.RbacClient)
	res, _ := rbacClient.AccessAdd(context.Background(), &pbRbac.AccessAddRequest{
		ModuleName:  moduleName,
		Type:        int64(accessType),
		ActionName:  actionName,
		Url:         url,
		ModuleId:    int64(moduleId),
		Sort:        int64(sort),
		Description: description,
		Status:      int64(status),
	})

	if !res.Success {
		con.Error(c, "增加数据失败", "/admin/access/add")
		return
	}
	con.Success(c, "增加数据成功", "/admin/access")

}

func (con AccessController) Edit(c *gin.Context) {
	//获取要修改的数据
	id, err1 := models.Int(c.Query("id"))
	if err1 != nil {
		con.Error(c, "参数错误", "/admin/access")
	}
	//获取当前id对应的access
	rbacClient := pbRbac.NewRbacAccessService("rbac", models.RbacClient)
	access, _ := rbacClient.AccessGet(context.Background(), &pbRbac.AccessGetRequest{
		Id: int64(id),
	})

	//获取顶级模块
	resAccess, _ := rbacClient.AccessGet(context.Background(), &pbRbac.AccessGetRequest{})

	c.HTML(http.StatusOK, "admin/access/edit.html", gin.H{
		"access":     access.AccessList[0],
		"accessList": resAccess.AccessList,
	})
}

func (con AccessController) DoEdit(c *gin.Context) {
	id, err1 := models.Int(c.PostForm("id"))
	moduleName := strings.Trim(c.PostForm("module_name"), " ")
	actionName := c.PostForm("action_name")
	accessType, err2 := models.Int(c.PostForm("type"))
	url := c.PostForm("url")
	moduleId, err3 := models.Int(c.PostForm("module_id"))
	sort, err4 := models.Int(c.PostForm("sort"))
	status, err5 := models.Int(c.PostForm("status"))
	description := c.PostForm("description")
	if err1 != nil || err2 != nil || err3 != nil || err4 != nil || err5 != nil {
		con.Error(c, "传入参数错误", "/admin/access")
		return
	}
	if moduleName == "" {
		con.Error(c, "模块名称不能为空", "/admin/access/edit?id="+models.String(id))
		return
	}
	rbacClient := pbRbac.NewRbacAccessService("rbac", models.RbacClient)
	accessRes, _ := rbacClient.AccessEdit(context.Background(), &pbRbac.AccessEditRequest{
		Id:          int64(id),
		ModuleName:  moduleName,
		Type:        int64(accessType),
		ActionName:  actionName,
		Url:         url,
		ModuleId:    int64(moduleId),
		Sort:        int64(sort),
		Description: description,
		Status:      int64(status),
	})

	if !accessRes.Success {
		con.Error(c, "修改数据失败", "/admin/access/edit?id="+models.String(id))
	} else {
		con.Success(c, "修改数据成功", "/admin/access/edit?id="+models.String(id))
	}

}

func (con AccessController) Delete(c *gin.Context) {
	id, err := models.Int(c.Query("id"))
	if err != nil {
		con.Error(c, "传入数据错误", "/admin/access")
	} else {
		//获取我们要删除的数据
		rbacClient := pbRbac.NewRbacAccessService("rbac", models.RbacClient)
		accessRes, _ := rbacClient.AccessDelete(context.Background(), &pbRbac.AccessDeleteRequest{
			Id: int64(id),
		})
		if !accessRes.Success { //顶级模块
			con.Error(c, accessRes.Message, "/admin/access")
		} else { //操作 或者菜单
			con.Success(c, "删除数据成功", "/admin/access")
		}

	}
}
