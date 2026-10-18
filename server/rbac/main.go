package main

import (
	"rbac/handler"
	"rbac/models"
	pbAccess "rbac/proto/rbacAccess"
	pbLogin "rbac/proto/rbacLogin"
	pbManager "rbac/proto/rbacManager"
	pbRole "rbac/proto/rbacRole"

	"github.com/go-micro/plugins/v4/registry/consul"
	"go-micro.dev/v4"
	log "go-micro.dev/v4/logger"
	"go-micro.dev/v4/registry"
)

var (
	service = "rbac"
	version = "latest"
)

func main() {
	addr := models.Config.Section("consul").Key("addr").String()
	consulReg := consul.NewRegistry(
		registry.Addrs(addr),
	)
	// Create service
	srv := micro.NewService(
		micro.Name(service),
		micro.Version(version),
		micro.Registry(consulReg),
	)
	srv.Init()

	// 注册登录的微服务
	pbLogin.RegisterRbacLoginHandler(srv.Server(), new(handler.RbacLogin))
	// 注册role的微服务
	pbRole.RegisterRbacRoleHandler(srv.Server(), new(handler.RbacRole))
	//注册manager微服务
	pbManager.RegisterRbacManagerHandler(srv.Server(), new(handler.RbacManager))
	//注册access微服务
	pbAccess.RegisterRbacAccessHandler(srv.Server(), new(handler.RbacAccess))
	// Run service
	if err := srv.Run(); err != nil {
		log.Fatal(err)
	}
}
