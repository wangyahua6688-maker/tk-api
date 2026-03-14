package handler

import (
	"net/http"
	"strings"

	tkv1 "tk-proto/tk/v1"
)

// TopicList 用户域帖子列表接口（旧名兼容，实际走新 ForumTopics RPC）。
// 路径迁移策略：
// 1) 新路径：/public/user/topics；
// 2) 兼容别名：/public/forum/topics。
func (h *PublicHandler) TopicList(w http.ResponseWriter, r *http.Request) {
	// 兼容别名命中时回传迁移提示，便于前端和网关逐步切换到 /public/user/*。
	if strings.Contains(r.URL.Path, "/public/forum/") {
		w.Header().Set("Deprecation", "true")
		w.Header().Set("Sunset", "Thu, 31 Dec 2026 23:59:59 GMT")
		w.Header().Set("Link", "</api/v1/public/user/topics>; rel=\"successor-version\"")
	}

	limit := parseIntOrDefault(r.URL.Query().Get("limit"), 20)
	feed := strings.TrimSpace(r.URL.Query().Get("feed"))
	keyword := strings.TrimSpace(r.URL.Query().Get("keyword"))
	issue := strings.TrimSpace(r.URL.Query().Get("issue"))
	year := parseIntOrDefault(r.URL.Query().Get("year"), 0)
	if strings.TrimSpace(r.URL.Query().Get("year")) == "" {
		year = 0
	}
	resp, err := h.svcCtx.User.ForumTopics(r.Context(), &tkv1.ForumTopicsRequest{
		Limit:   int32(limit),
		Feed:    feed,
		Keyword: keyword,
		Issue:   issue,
		Year:    int32(year),
	})
	writeRPCReply(w, resp, err)
}
