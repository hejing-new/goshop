package handler

import (
	"context"
	"rbac/models"
	"strconv"

	pb "rbac/proto/rbacRole"

	"gorm.io/gorm"
)

type RbacRole struct{}

//获取角色
func (e *RbacRole) RoleGet(ctx context.Context, req *pb.RoleGetRequest, res *pb.RoleGetResponse) error {
	roleList := []models.Role{}
	where := "1=1"
	if req.Id > 0 {
		where += " AND id=" + strconv.Itoa(int(req.Id))
	}
	models.DB.Where(where).Find(&roleList)

	//处理数据
	var tempList []*pb.RoleModel
	for _, v := range roleList {
		tempList = append(tempList, &pb.RoleModel{
			Id:          int64(v.Id),
			Title:       v.Title,
			Description: v.Description,
			Status:      int64(v.Status),
			AddTime:     int64(v.AddTime),
		})
	}
	res.RoleList = tempList
	return nil
}

//增加角色
func (e *RbacRole) RoleAdd(ctx context.Context, req *pb.RoleAddRequest, res *pb.RoleAddResponse) error {

	role := models.Role{}
	role.Title = req.Title
	role.Description = req.Description
	role.Status = int(req.Status)
	role.AddTime = int(req.AddTime)

	err := models.DB.Create(&role).Error
	if err != nil {
		res.Success = false
		res.Message = "增加数据失败"
	} else {
		res.Success = true
		res.Message = "增加数据成功"
	}
	return err
}

//修改角色
func (e *RbacRole) RoleEdit(ctx context.Context, req *pb.RoleEditRequest, res *pb.RoleEditResponse) error {

	role := models.Role{Id: int(req.Id)}
	models.DB.Find(&role)
	role.Title = req.Title
	role.Description = req.Description

	err := models.DB.Save(&role).Error
	if err != nil {
		res.Success = false
		res.Message = "修改数据失败"
	} else {
		res.Success = true
		res.Message = "修改数据成功"
	}
	return nil
}

//删除角色
func (e *RbacRole) RoleDelete(ctx context.Context, req *pb.RoleDeleteRequest, res *pb.RoleDeleteResponse) error {
	role := models.Role{Id: int(req.Id)}
	err := models.DB.Delete(&role).Error
	if err != nil {
		res.Success = false
		res.Message = "删除数据失败"
	} else {
		res.Success = true
		res.Message = "删除数据成功"
	}
	return nil
}

//授权
func (e *RbacRole) RoleAuth(ctx context.Context, req *pb.RoleAuthRequest, res *pb.RoleAuthResponse) error {
	//1、获取角色id  req.RoleId
	//2、获取所有的权限
	accessList := []models.Access{}
	models.DB.Where("module_id=?", 0).Preload("AccessItem", func(db *gorm.DB) *gorm.DB {
		return db.Order("access.sort DESC")
	}).Order("sort DESC").Find(&accessList)

	//3、获取当前角色拥有的权限 ，并把权限id放在一个map对象里面
	roleAccess := []models.RoleAccess{}
	models.DB.Where("role_id=?", req.RoleId).Find(&roleAccess)
	roleAccessMap := make(map[int]int)
	for _, v := range roleAccess {
		roleAccessMap[v.AccessId] = v.AccessId
	}

	//4、循环遍历所有的权限数据，判断当前权限的id是否在角色权限的Map对象中,如果是的话给当前数据加入checked属性
	for i := 0; i < len(accessList); i++ {
		if _, ok := roleAccessMap[int(accessList[i].Id)]; ok {
			accessList[i].Checked = true
		}
		for j := 0; j < len(accessList[i].AccessItem); j++ {
			if _, ok := roleAccessMap[int(accessList[i].AccessItem[j].Id)]; ok {
				accessList[i].AccessItem[j].Checked = true
			}
		}
	}
	//处理数据
	var tempList []*pb.AccessModel
	for _, v := range accessList {
		var tempItemList []*pb.AccessModel
		for _, k := range v.AccessItem {
			tempItemList = append(tempItemList, &pb.AccessModel{
				Id:          int64(k.Id),
				ModuleName:  k.ModuleName,
				ActionName:  k.ActionName,
				Type:        int64(k.Type),
				Url:         k.Url,
				ModuleId:    int64(k.ModuleId),
				Sort:        int64(k.Sort),
				Description: k.Description,
				Status:      int64(k.Status),
				Checked:     k.Checked,
				AddTime:     int64(k.AddTime),
			})
		}
		tempList = append(tempList, &pb.AccessModel{
			Id:          int64(v.Id),
			ModuleName:  v.ModuleName,
			ActionName:  v.ActionName,
			Type:        int64(v.Type),
			Url:         v.Url,
			ModuleId:    int64(v.ModuleId),
			Sort:        int64(v.Sort),
			Description: v.Description,
			Status:      int64(v.Status),
			AddTime:     int64(v.AddTime),
			Checked:     v.Checked,
			AccessItem:  tempItemList,
		})
	}

	res.AccessList = tempList

	return nil
}

//执行授权
func (e *RbacRole) RoleDoAuth(ctx context.Context, req *pb.RoleDoAuthRequest, res *pb.RoleDoAuthResponse) error {
	//1、删除当前角色对应的权限
	roleAccess := models.RoleAccess{}
	models.DB.Where("role_id=?", req.RoleId).Delete(&roleAccess)

	//2、增加当前角色对应的权限
	for _, v := range req.AccessIds {
		roleAccess.RoleId = int(req.RoleId)
		accessId, _ := strconv.Atoi(v)
		roleAccess.AccessId = accessId
		models.DB.Create(&roleAccess)
	}
	res.Success = true
	res.Message = "授权成功"
	return nil
}

//权限判断
func (e *RbacRole) MiddlewaresAuth(ctx context.Context, req *pb.MiddlewaresAuthRequest, res *pb.MiddlewaresAuthResponse) error {
	// 1、根据角色获取当前角色的权限列表,然后把权限id放在一个map类型的对象里面
	roleAccess := []models.RoleAccess{}
	models.DB.Where("role_id=?", req.RoleId).Find(&roleAccess)
	roleAccessMap := make(map[int]int)
	for _, v := range roleAccess {
		roleAccessMap[v.AccessId] = v.AccessId
	}
	// 2、获取当前访问的url对应的权限id 判断权限id是否在角色对应的权限
	// pathname      /admin/manager
	access := models.Access{}
	models.DB.Where("url = ?", req.UrlPath).Find(&access)
	//3、判断当前访问的url对应的权限id 是否在权限列表的id中
	if _, ok := roleAccessMap[int(access.Id)]; !ok {
		res.HasPermission = false
	} else {
		res.HasPermission = true
	}
	return nil
}
