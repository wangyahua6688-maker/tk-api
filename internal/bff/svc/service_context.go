package svc

import (
	"tk-api/internal/bff/config"
	tkv1 "tk-proto/gen/go/tk/v1"

	"github.com/zeromicro/go-zero/zrpc"
)

// ServiceContext BFF 的依赖集合。
type ServiceContext struct {
	// Config 保存网关配置。
	Config config.Config
	// Business 业务域 gRPC 客户端（首页+开奖+图纸+投票+现场页）。
	Business tkv1.BusinessServiceClient
	// User 用户域 gRPC 客户端（论坛帖子与评论能力）。
	User tkv1.UserServiceClient
}

// NewServiceContext 创建ServiceContext实例。
func NewServiceContext(c config.Config) *ServiceContext {
	// 初始化业务域客户端。
	businessClient := zrpc.MustNewClient(c.BusinessRpc)
	// 初始化用户域客户端。
	userClient := zrpc.MustNewClient(c.UserRpc)

	// 将两个领域客户端注入 BFF 上下文，供 handler 统一复用。
	// 注意：这里不做业务逻辑，只做依赖装配。
	return &ServiceContext{
		// 处理当前语句逻辑。
		Config: c,
		// 调用tkv1.NewBusinessServiceClient完成当前处理。
		Business: tkv1.NewBusinessServiceClient(businessClient.Conn()),
		// 调用tkv1.NewUserServiceClient完成当前处理。
		User: tkv1.NewUserServiceClient(userClient.Conn()),
	}
}
