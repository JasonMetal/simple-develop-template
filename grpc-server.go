package main

import (
	"github.com/JasonMetal/simple-develop-template/internal/constant"
	"github.com/JasonMetal/simple-develop-template/pkg/support-go/bootstrap"
	grpcHealth "google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
)

func main() {
	bootstrap.SetProjectName(constant.ProjectName)
	bootstrap.Init()
	s := bootstrap.NewGrpcServer()
	// 健康检测
	grpc_health_v1.RegisterHealthServer(s, grpcHealth.NewServer())

	// 业务服务
	bootstrap.RunServer(s, constant.GrpcServiceHostPort)
}
