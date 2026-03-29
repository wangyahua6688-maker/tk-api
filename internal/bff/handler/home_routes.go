package handler

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest"
)

// RegisterHomeRoutes 注册首页相关路由。
// 路径规则：
// business 服务能力统一使用 /public/business/* 作为规范主路径。
func RegisterHomeRoutes(server *rest.Server, home *HomeHandler) {
	server.AddRoutes([]rest.Route{
		{Method: http.MethodGet, Path: "/api/v1/public/business/home", Handler: home.HomeOverview},
		{Method: http.MethodGet, Path: "/api/v1/public/business/live-scene", Handler: home.LiveScenePage},
		{Method: http.MethodGet, Path: "/api/v1/public/business/lottery-categories", Handler: home.LotteryCategories},
	})
}
