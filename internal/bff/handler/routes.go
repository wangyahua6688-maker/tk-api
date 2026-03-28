package handler

import "github.com/zeromicro/go-zero/rest"

// RegisterHandlers 注册 tk-api(BFF) 对外 HTTP 路由。
// 修复点：
//   - 短信发送 / 注册 / 登录接口统一包装 withIPContext，
//     将客户端真实 IP 注入 context，供 tk-user 的 IP 维度频控使用
func RegisterHandlers(server *rest.Server, home *HomeHandler, lottery *LotteryHandler, forum *ForumHandler, expert *ExpertHandler, auth *UserAuthHandler) {
	RegisterHealthRoutes(server)
	RegisterHomeRoutes(server, home)
	RegisterLotteryRoutes(server, lottery)
	RegisterForumRoutes(server, forum)
	RegisterExpertRoutes(server, expert)
	RegisterAuthRoutes(server, auth)
}
