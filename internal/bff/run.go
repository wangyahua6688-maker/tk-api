package bff

import (
	"flag"

	"tk-api/internal/bff/config"
	"tk-api/internal/bff/handler"
	"tk-api/internal/bff/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest"
)

// Run 启动 tk-api BFF 服务。
func Run() {
	// 1) 读取启动参数（配置文件路径）。
	configFile := flag.String("f", "etc/tk-api.yaml", "config file")
	flag.Parse()

	// 2) 加载网关配置：HTTP 监听 + 下游 RPC 地址。
	var c config.Config
	conf.MustLoad(*configFile, &c)

	// 3) 初始化服务上下文（gRPC 客户端）。
	svcCtx := svc.NewServiceContext(c)
	// 4) 创建公开路由处理器。
	publicHandler := handler.NewPublicHandler(svcCtx)

	// 5) 启动 go-zero REST 服务（启用 CORS，便于本地联调）。
	server := rest.MustNewServer(c.RestConf, rest.WithCors("*"))
	defer server.Stop()

	// 6) 注册 HTTP 路由。
	handler.RegisterHandlers(server, publicHandler)
	// 7) 输出启动日志并进入监听。
	logx.Infof("starting tk-api bff at %s:%d", c.Host, c.Port)
	server.Start()
}
