package handler

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest"
)

// RegisterHomeRoutes 注册首页相关路由。
// 路径规则：
// 1) business 服务能力统一使用 /public/business/* 作为规范主路径；
// 2) 旧的 /public/* 路径保留兼容，并在响应头中返回迁移提示。
func RegisterHomeRoutes(server *rest.Server, home *HomeHandler) {
	server.AddRoutes([]rest.Route{
		{Method: http.MethodGet, Path: "/api/v1/public/business/home", Handler: home.HomeOverview},
		{Method: http.MethodGet, Path: "/api/v1/public/business/live-scene", Handler: home.LiveScenePage},
		{Method: http.MethodGet, Path: "/api/v1/public/business/lottery-categories", Handler: home.LotteryCategories},
		{Method: http.MethodGet, Path: "/api/v1/public/home", Handler: home.HomeOverview},
		{Method: http.MethodGet, Path: "/api/v1/public/live-scene", Handler: home.LiveScenePage},
		{Method: http.MethodGet, Path: "/api/v1/public/lottery-categories", Handler: home.LotteryCategories},
	})
}
