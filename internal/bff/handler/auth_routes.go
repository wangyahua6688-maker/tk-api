package handler

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest"
)

// RegisterAuthRoutes 注册用户认证相关路由。
func RegisterAuthRoutes(server *rest.Server, auth *UserAuthHandler) {
	server.AddRoutes([]rest.Route{
		{
			Method:  http.MethodPost,
			Path:    "/api/v1/public/user/auth/sms-code",
			Handler: withIPContext(auth.SendSMSCode),
		},
		{
			Method:  http.MethodPost,
			Path:    "/api/v1/public/user/auth/register",
			Handler: withIPContext(auth.RegisterByPhone),
		},
		{
			Method:  http.MethodPost,
			Path:    "/api/v1/public/user/auth/login/password",
			Handler: withIPContext(auth.LoginByPassword),
		},
		{
			Method:  http.MethodPost,
			Path:    "/api/v1/public/user/auth/login/sms",
			Handler: withIPContext(auth.LoginBySMS),
		},
		{Method: http.MethodGet, Path: "/api/v1/public/user/profile", Handler: auth.Profile},
	})
}
