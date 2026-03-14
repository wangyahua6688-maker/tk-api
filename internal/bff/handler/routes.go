package handler

import (
	"net/http"

	"tk-common/utils/httpresp"

	"github.com/zeromicro/go-zero/rest"
)

// RegisterHandlers 注册 tk-api(BFF) 对外 HTTP 路由。
// 设计原则：
// 1. 路由层只做协议映射，不承载业务逻辑；
// 2. 业务逻辑全部通过 gRPC 调用下游微服务；
// 3. 接口路径保持与 tk-web 现网调用兼容。
func RegisterHandlers(server *rest.Server, h *PublicHandler) {
	server.AddRoutes([]rest.Route{
		{
			Method: http.MethodGet,
			Path:   "/health",
			Handler: func(w http.ResponseWriter, _ *http.Request) {
				httpresp.OK(w, map[string]interface{}{
					"status":  "ok",
					"service": "tk-api-bff",
				})
			},
		},
		{Method: http.MethodGet, Path: "/api/v1/public/home", Handler: h.HomeOverview},
		{Method: http.MethodGet, Path: "/api/v1/public/live-scene", Handler: h.LiveScenePage},
		{Method: http.MethodGet, Path: "/api/v1/public/lottery-categories", Handler: h.LotteryCategories},
		// 论坛接口新主路径：/public/user/*
		{Method: http.MethodGet, Path: "/api/v1/public/user/topics", Handler: h.TopicList},
		{Method: http.MethodGet, Path: "/api/v1/public/user/topics/:id/detail", Handler: h.TopicDetail},
		{Method: http.MethodGet, Path: "/api/v1/public/user/users/:id/history-topics", Handler: h.AuthorHistory},
		// 兼容别名：保留 /public/forum/*，后续灰度下线。
		{Method: http.MethodGet, Path: "/api/v1/public/forum/topics", Handler: h.TopicList},
		{Method: http.MethodGet, Path: "/api/v1/public/forum/topics/:id/detail", Handler: h.TopicDetail},
		{Method: http.MethodGet, Path: "/api/v1/public/forum/users/:id/history-topics", Handler: h.AuthorHistory},
		// 高手推荐接口新主路径：/public/user/*。
		{Method: http.MethodGet, Path: "/api/v1/public/user/expert-boards", Handler: h.ExpertBoards},
		// 兼容别名：保留 /public/forum/*。
		{Method: http.MethodGet, Path: "/api/v1/public/forum/expert-boards", Handler: h.ExpertBoards},
		// 用户认证接口：验证码、注册、登录、资料。
		{Method: http.MethodPost, Path: "/api/v1/public/user/auth/sms-code", Handler: h.SendSMSCode},
		{Method: http.MethodPost, Path: "/api/v1/public/user/auth/register", Handler: h.RegisterByPhone},
		{Method: http.MethodPost, Path: "/api/v1/public/user/auth/login/password", Handler: h.LoginByPassword},
		{Method: http.MethodPost, Path: "/api/v1/public/user/auth/login/sms", Handler: h.LoginBySMS},
		{Method: http.MethodGet, Path: "/api/v1/public/user/profile", Handler: h.Profile},
		{Method: http.MethodGet, Path: "/api/v1/public/lottery-cards", Handler: h.LotteryCards},
		{Method: http.MethodGet, Path: "/api/v1/public/special-lotteries/:id/dashboard", Handler: h.LotteryDashboard},
		{Method: http.MethodGet, Path: "/api/v1/public/special-lotteries/:id/history", Handler: h.DrawHistory},
		{Method: http.MethodGet, Path: "/api/v1/public/draw-records/:id/detail", Handler: h.DrawDetail},
		{Method: http.MethodGet, Path: "/api/v1/public/lottery-info/:id/detail", Handler: h.LotteryDetail},
		{Method: http.MethodGet, Path: "/api/v1/public/lottery-info/:id/history", Handler: h.LotteryHistory},
		{Method: http.MethodGet, Path: "/api/v1/public/lottery-info/:id/results", Handler: h.LotteryResults},
		{Method: http.MethodGet, Path: "/api/v1/public/lottery-info/:id/vote-record", Handler: h.VoteRecord},
		{Method: http.MethodPost, Path: "/api/v1/public/lottery-info/:id/vote", Handler: h.Vote},
	})
}
