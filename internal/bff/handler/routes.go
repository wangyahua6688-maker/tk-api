package handler

import (
	"net/http"

	"tk-common/utils/httpresp"

	"github.com/zeromicro/go-zero/rest"
)

// RegisterHandlers 注册 tk-api(BFF) 对外 HTTP 路由。
// 修复点：
//   - 短信发送 / 注册 / 登录接口统一包装 withIPContext，
//     将客户端真实 IP 注入 context，供 tk-user 的 IP 维度频控使用
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

		// ── 首页 & 开奖 ────────────────────────────────
		{Method: http.MethodGet, Path: "/api/v1/public/home", Handler: h.HomeOverview},
		{Method: http.MethodGet, Path: "/api/v1/public/live-scene", Handler: h.LiveScenePage},
		{Method: http.MethodGet, Path: "/api/v1/public/lottery-categories", Handler: h.LotteryCategories},
		{Method: http.MethodGet, Path: "/api/v1/public/lottery-cards", Handler: h.LotteryCards},
		{Method: http.MethodGet, Path: "/api/v1/public/special-lotteries/:id/dashboard", Handler: h.LotteryDashboard},
		{Method: http.MethodGet, Path: "/api/v1/public/special-lotteries/:id/history", Handler: h.DrawHistory},
		{Method: http.MethodGet, Path: "/api/v1/public/draw-records/:id/detail", Handler: h.DrawDetail},
		{Method: http.MethodGet, Path: "/api/v1/public/lottery-info/:id/detail", Handler: h.LotteryDetail},
		{Method: http.MethodGet, Path: "/api/v1/public/lottery-info/:id/history", Handler: h.LotteryHistory},
		{Method: http.MethodGet, Path: "/api/v1/public/lottery-info/:id/results", Handler: h.LotteryResults},
		{Method: http.MethodGet, Path: "/api/v1/public/lottery-info/:id/vote-record", Handler: h.VoteRecord},
		{Method: http.MethodPost, Path: "/api/v1/public/lottery-info/:id/vote", Handler: h.Vote},

		// ── 论坛（主路径 + 兼容别名）──────────────────────
		{Method: http.MethodGet, Path: "/api/v1/public/user/topics", Handler: h.TopicList},
		{Method: http.MethodGet, Path: "/api/v1/public/user/topics/:id/detail", Handler: h.TopicDetail},
		{Method: http.MethodGet, Path: "/api/v1/public/user/users/:id/history-topics", Handler: h.AuthorHistory},
		{Method: http.MethodGet, Path: "/api/v1/public/forum/topics", Handler: h.TopicList},
		{Method: http.MethodGet, Path: "/api/v1/public/forum/topics/:id/detail", Handler: h.TopicDetail},
		{Method: http.MethodGet, Path: "/api/v1/public/forum/users/:id/history-topics", Handler: h.AuthorHistory},

		// ── 高手榜 ───────────────────────────────────────
		{Method: http.MethodGet, Path: "/api/v1/public/user/expert-boards", Handler: h.ExpertBoards},
		{Method: http.MethodGet, Path: "/api/v1/public/forum/expert-boards", Handler: h.ExpertBoards},

		// ── 用户认证（注入 IP context，支持 IP 维度限流）──────
		// withIPContext 将请求的真实客户端 IP 写入 context，
		// 供 tk-user 服务的短信发送 IP 频控使用
		{
			Method:  http.MethodPost,
			Path:    "/api/v1/public/user/auth/sms-code",
			Handler: withIPContext(h.SendSMSCode), // ← IP context 注入
		},
		{
			Method:  http.MethodPost,
			Path:    "/api/v1/public/user/auth/register",
			Handler: withIPContext(h.RegisterByPhone), // ← IP context 注入（分布式锁防并发注册）
		},
		{
			Method:  http.MethodPost,
			Path:    "/api/v1/public/user/auth/login/password",
			Handler: withIPContext(h.LoginByPassword), // ← IP context 注入（登录失败计数）
		},
		{
			Method:  http.MethodPost,
			Path:    "/api/v1/public/user/auth/login/sms",
			Handler: h.LoginBySMS, // 短信登录已有验证码保护，无需单独 IP 限流
		},
		{Method: http.MethodGet, Path: "/api/v1/public/user/profile", Handler: h.Profile},
	})
}
