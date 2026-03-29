package handler

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest"
)

// RegisterForumRoutes 注册论坛相关路由。
func RegisterForumRoutes(server *rest.Server, forum *ForumHandler) {
	server.AddRoutes([]rest.Route{
		{Method: http.MethodGet, Path: "/api/v1/public/user/topics", Handler: forum.TopicList},
		{Method: http.MethodGet, Path: "/api/v1/public/user/topics/:id/detail", Handler: forum.TopicDetail},
		{Method: http.MethodGet, Path: "/api/v1/public/user/users/:id/history-topics", Handler: forum.AuthorHistory},
	})
}
