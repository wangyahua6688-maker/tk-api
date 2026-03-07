package config

import (
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

// Config BFF 网关配置：对外提供 HTTP 接口，对内通过 gRPC 调用领域服务。
type Config struct {
	rest.RestConf
	// BusinessRpc 业务域服务（首页、开奖、图纸、投票、现场页等）。
	BusinessRpc zrpc.RpcClientConf
	// UserRpc 用户域服务（论坛、评论分组等）。
	UserRpc zrpc.RpcClientConf
}
