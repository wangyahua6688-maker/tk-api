package handler

import (
	"net/http"

	"github.com/wangyahua6688-maker/tk-common/utils/httpresp"
	"github.com/zeromicro/go-zero/rest"
)

// RegisterHealthRoutes 注册健康检查路由。
func RegisterHealthRoutes(server *rest.Server) {
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
	})
}
