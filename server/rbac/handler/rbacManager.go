package handler

import (
	"context"

	"rbac/models"
	pb "rbac/proto/rbacManager"
)

type RbacManager struct{}

//获取管理员
func (e *RbacManager) ManagerGet(ctx context.Context, req *pb.ManagerGetRequest, res *pb.ManagerGetResponse) error {
	managerList := []models.Manager{}
	// ★ 全部用占位符传参：原来 where += " AND username=" + req.Username 是直接拼接，
	//   传一个 1 or 1=1 就能查出所有管理员（SQL 注入）
	db := models.DB.Model(&models.Manager{})
	if req.Id > 0 {
		db = db.Where("id = ?", req.Id)
	}
	if len(req.Username) > 0 {
		db = db.Where("username = ?", req.Username)
	}
	db.Preload("Role").Find(&managerList)
	//处理数据
	var tempList []*pb.ManagerModel
	for _, v := range managerList {
		tempList = append(tempList, &pb.ManagerModel{
			Id:       int64(v.Id),
			Username: v.Username,
			Mobile:   v.Mobile,
			Email:    v.Email,
			Status:   int64(v.Status),
			RoleId:   int64(v.RoleId),
			AddTime:  int64(v.AddTime),
			IsSuper:  int64(v.IsSuper),
			Role: &pb.RoleModel{
				Title:       v.Role.Title,
				Description: v.Role.Description,
			},
		})
	}

	res.ManagerList = tempList
	return nil
}

//增加管理员
func (e *RbacManager) ManagerAdd(ctx context.Context, req *pb.ManagerAddRequest, res *pb.ManagerAddResponse) error {

	//执行增加管理员
	manager := models.Manager{
		Username: req.Username,
		Password: req.Password,
		Email:    req.Email,
		Mobile:   req.Mobile,
		RoleId:   int(req.RoleId),
		Status:   int(req.Status),
		AddTime:  int(req.AddTime),
	}
	err := models.DB.Create(&manager).Error
	if err != nil {
		res.Success = false
		res.Message = "增加数据失败"
	} else {
		res.Success = true
		res.Message = "增加数据成功"
	}
	return err
}

//修改管理员
func (e *RbacManager) ManagerEdit(ctx context.Context, req *pb.ManagerEditRequest, res *pb.ManagerEditResponse) error {

	//执行修改
	manager := models.Manager{Id: int(req.Id)}
	models.DB.Find(&manager)
	manager.Username = req.Username
	manager.Email = req.Email
	manager.Mobile = req.Mobile
	manager.RoleId = int(req.RoleId)

	//注意：判断密码是否为空 为空表示不修改密码 不为空表示修改密码
	if req.Password != "" {
		manager.Password = req.Password
	}
	err := models.DB.Save(&manager).Error
	if err != nil {
		res.Success = false
		res.Message = "修改数据失败"
	} else {
		res.Success = true
		res.Message = "修改数据成功"
	}
	return err
}

//删除管理员
func (e *RbacManager) ManagerDelete(ctx context.Context, req *pb.ManagerDeleteRequest, res *pb.ManagerDeleteResponse) error {

	manager := models.Manager{Id: int(req.Id)}
	err := models.DB.Delete(&manager).Error
	if err != nil {
		res.Success = false
		res.Message = "删除数据失败"
	} else {
		res.Success = true
		res.Message = "删除数据成功"
	}

	return err
}
