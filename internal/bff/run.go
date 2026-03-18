package bff

import (
	"flag"
	"strings"

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
	// 调用flag.Parse完成当前处理。
	flag.Parse()

	// 2) 加载网关配置：HTTP 监听 + 下游 RPC 地址。
	var c config.Config
	// 调用conf.MustLoad完成当前处理。
	conf.MustLoad(*configFile, &c)

	// 3) 初始化服务上下文（gRPC 客户端）。
	svcCtx := svc.NewServiceContext(c)
	// 4) 创建公开路由处理器。
	publicHandler := handler.NewPublicHandler(svcCtx)

	// 5) 构建 REST 启动选项：仅对白名单来源启用 CORS。
	serverOptions := make([]rest.RunOption, 0, 1)
	// 定义并初始化当前变量。
	allowedOrigins := collectAllowedOrigins(c.CORS.AllowedOrigins)
	// 判断条件并进入对应分支逻辑。
	if len(allowedOrigins) > 0 {
		// 更新当前变量或字段值。
		serverOptions = append(serverOptions, rest.WithCors(allowedOrigins...))
	}
	// 6) 启动 go-zero REST 服务。
	server := rest.MustNewServer(c.RestConf, serverOptions...)
	// 注册延迟执行逻辑。
	defer server.Stop()

	// 7) 注册 HTTP 路由。
	handler.RegisterHandlers(server, publicHandler)
	// 8) 输出启动日志并进入监听。
	logx.Infof("starting tk-api bff at %s:%d", c.Host, c.Port)
	// 调用server.Start完成当前处理。
	server.Start()
}

// collectAllowedOrigins 归一化并过滤空白 CORS 来源配置。
func collectAllowedOrigins(origins []string) []string {
	// 定义并初始化当前变量。
	result := make([]string, 0, len(origins))
	// 遍历当前集合并处理元素。
	for _, origin := range origins {
		// 更新当前变量或字段值。
		origin = strings.TrimSpace(origin)
		// 判断条件并进入对应分支逻辑。
		if origin == "" {
			// 继续下一次循环。
			continue
		}
		// 更新当前变量或字段值。
		result = append(result, origin)
	}
	// 返回当前处理结果。
	return result
}
