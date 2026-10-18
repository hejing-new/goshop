package handler

import (
	"context"
	"fmt"
	"rbac/models"
	pb "rbac/proto"
)

type Rbac struct{}

func (e *Rbac) Login(ctx context.Context, req *pb.LoginRequest, res *pb.LoginResponse) error {
	managerList := []models.Manager{}
	err := models.DB.Where("username=? AND password=?", req.Username, req.Password).Find(&managerList)
	if err != nil {
		fmt.Println(err)
	}

	var templist []*pb.ManagerModel
	for i := 0; i < len(managerList); i++ {
		templist = append(templist, &pb.ManagerModel{
			Id:       int64(managerList[i].Id),
			Username: managerList[i].Username,
			Password: managerList[i].Password,
			Mobile:   managerList[i].Mobile,
			Email:    managerList[i].Email,
			Status:   int64(managerList[i].Status),
			RoleId:   int64(managerList[i].RoleId),
			AddTime:  int64(managerList[i].AddTime),
			IsSuper:  int64(managerList[i].IsSuper),
		})
	}

	if len(managerList) > 0 {
		res.IsLogin = true

	} else {
		res.IsLogin = false
	}
	res.Userlist = templist
	return nil
}
