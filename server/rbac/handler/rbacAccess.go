package handler

import (
	"context"
	"strconv"

	"rbac/models"
	pb "rbac/proto/rbacAccess"
)

type RbacAccess struct{}

//获取权限
func (e *RbacAccess) AccessGet(ctx context.Context, req *pb.AccessGetRequest, res *pb.AccessGetResponse) error {
	accessList := []models.Access{}
	where := "1=1"
	if req.Id > 0 {
		where += " AND id=" + strconv.Itoa(int(req.Id))
	} else {
		where += " AND module_id = 0"
	}
	models.DB.Where(where).Preload("AccessItem").Find(&accessList)

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
			AccessItem:  tempItemList,
		})
	}

	res.AccessList = tempList
	return nil
}

//增加权限
func (e *RbacAccess) AccessAdd(ctx context.Context, req *pb.AccessAddRequest, res *pb.AccessAddResponse) error {
	access := models.Access{
		ModuleName:  req.ModuleName,
		Type:        req.Type,
		ActionName:  req.ActionName,
		Url:         req.Url,
		ModuleId:    req.ModuleId,
		Sort:        req.Sort,
		Description: req.Description,
		Status:      req.Status,
	}
	err := models.DB.Create(&access).Error
	if err != nil {
		res.Success = false
		res.Message = "增加数据失败"
	} else {
		res.Success = true
		res.Message = "增加数据成功"
	}
	return err
}

//修改权限
func (e *RbacAccess) AccessEdit(ctx context.Context, req *pb.AccessEditRequest, res *pb.AccessEditResponse) error {
	access := models.Access{Id: req.Id}
	models.DB.Find(&access)

	access.ModuleName = req.ModuleName
	access.Type = req.Type
	access.ActionName = req.ActionName
	access.Url = req.Url
	access.ModuleId = req.ModuleId
	access.Sort = req.Sort
	access.Description = req.Description
	access.Status = req.Status

	err := models.DB.Save(&access).Error
	if err != nil {
		res.Success = false
		res.Message = "修改数据失败"
	} else {
		res.Success = true
		res.Message = "修改数据成功"
	}
	return err
}

//删除权限
func (e *RbacAccess) AccessDelete(ctx context.Context, req *pb.AccessDeleteRequest, res *pb.AccessDeleteResponse) error {

	access := models.Access{Id: req.Id}
	models.DB.Find(&access)
	if access.ModuleId == 0 { //顶级模块
		accessList := []models.Access{}
		models.DB.Where("module_id = ?", access.Id).Find(&accessList)

		if len(accessList) > 0 {
			res.Success = false
			res.Message = "当前模块下面有菜单或者操作，请删除菜单或者操作以后再来删除这个数据"
		} else {
			models.DB.Delete(&access)
			res.Success = true
			res.Message = "删除数据成功"
		}
	} else { //操作 或者菜单
		models.DB.Delete(&access)
		res.Success = true
		res.Message = "删除数据成功"
	}
	return nil
}
