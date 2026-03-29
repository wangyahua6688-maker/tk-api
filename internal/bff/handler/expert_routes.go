package handler

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest"
)

// RegisterExpertRoutes 注册高手榜相关路由。
func RegisterExpertRoutes(server *rest.Server, expert *ExpertHandler) {
	server.AddRoutes([]rest.Route{
		{Method: http.MethodGet, Path: "/api/v1/public/user/expert-boards", Handler: expert.ExpertBoards},
	})
}
