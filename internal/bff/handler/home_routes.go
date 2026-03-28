package handler

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest"
)

// RegisterHomeRoutes 注册首页相关路由。
func RegisterHomeRoutes(server *rest.Server, home *HomeHandler) {
	server.AddRoutes([]rest.Route{
		{Method: http.MethodGet, Path: "/api/v1/public/home", Handler: home.HomeOverview},
		{Method: http.MethodGet, Path: "/api/v1/public/live-scene", Handler: home.LiveScenePage},
		{Method: http.MethodGet, Path: "/api/v1/public/lottery-categories", Handler: home.LotteryCategories},
	})
}
