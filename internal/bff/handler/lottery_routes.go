package handler

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest"
)

// RegisterLotteryRoutes 注册开奖与图纸相关路由。
func RegisterLotteryRoutes(server *rest.Server, lottery *LotteryHandler) {
	server.AddRoutes([]rest.Route{
		{Method: http.MethodGet, Path: "/api/v1/public/lottery-cards", Handler: lottery.LotteryCards},
		{Method: http.MethodGet, Path: "/api/v1/public/special-lotteries/:id/dashboard", Handler: lottery.LotteryDashboard},
		{Method: http.MethodGet, Path: "/api/v1/public/special-lotteries/:id/history", Handler: lottery.DrawHistory},
		{Method: http.MethodGet, Path: "/api/v1/public/draw-records/:id/detail", Handler: lottery.DrawDetail},
		{Method: http.MethodGet, Path: "/api/v1/public/lottery-info/:id/detail", Handler: lottery.LotteryDetail},
		{Method: http.MethodGet, Path: "/api/v1/public/lottery-info/:id/history", Handler: lottery.LotteryHistory},
		{Method: http.MethodGet, Path: "/api/v1/public/lottery-info/:id/results", Handler: lottery.LotteryResults},
		{Method: http.MethodGet, Path: "/api/v1/public/lottery-info/:id/vote-record", Handler: lottery.VoteRecord},
		{Method: http.MethodPost, Path: "/api/v1/public/lottery-info/:id/vote", Handler: lottery.Vote},
	})
}
