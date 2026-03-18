package config

import (
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

// Config BFF 网关配置：对外提供 HTTP 接口，对内通过 gRPC 调用领域服务。
type Config struct {
	rest.RestConf                    // HTTP 服务配置（监听地址、端口、中间件）。
	BusinessRpc   zrpc.RpcClientConf // 业务域 RPC 客户端（首页、开奖、投票、现场页等）。
	UserRpc       zrpc.RpcClientConf // 用户域 RPC 客户端（论坛、评论、用户态相关能力）。
}
