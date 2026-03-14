package handler

import (
	"net/http"
	"strings"

	tkv1 "tk-proto/tk/v1"
)

// AuthorHistory 论坛作者历史贴列表接口。
// 路径策略：
// 1) 新路径：/public/user/users/:id/history-topics；
// 2) 兼容别名：/public/forum/users/:id/history-topics。
func (h *PublicHandler) AuthorHistory(w http.ResponseWriter, r *http.Request) {
	// 命中兼容别名时回传迁移提示头。
	if strings.Contains(r.URL.Path, "/public/forum/") {
		w.Header().Set("Deprecation", "true")
		w.Header().Set("Sunset", "Thu, 31 Dec 2026 23:59:59 GMT")
		w.Header().Set("Link", "</api/v1/public/user/users/:id/history-topics>; rel=\"successor-version\"")
	}

	// 从路径中解析用户ID。
	userID, ok := mustPathID(w, r, "users")
	if !ok {
		return
	}

	// 读取查询参数：limit/issue/year。
	limit := parseIntOrDefault(r.URL.Query().Get("limit"), 30)
	issue := strings.TrimSpace(r.URL.Query().Get("issue"))
	yearRaw := strings.TrimSpace(r.URL.Query().Get("year"))
	year := parseIntOrDefault(yearRaw, 0)
	if yearRaw == "" {
		year = 0
	}

	// 转发到用户域 RPC。
	resp, err := h.svcCtx.User.ForumAuthorHistory(r.Context(), &tkv1.ForumAuthorHistoryRequest{
		UserId: userID,
		Limit:  int32(limit),
		Issue:  issue,
		Year:   int32(year),
	})
	writeRPCReply(w, resp, err)
}
