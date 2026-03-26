package handler

import (
	"net/http"
	"strings"

	tkv1 "tk-proto/gen/go/tk/v1"
)

// AuthorHistory 论坛作者历史贴列表接口。
// 路径策略：
// 1) 新路径：/public/user/users/:id/history-topics；
// 2) 兼容别名：/public/forum/users/:id/history-topics。
func (h *PublicHandler) AuthorHistory(w http.ResponseWriter, r *http.Request) {
	// 命中兼容别名时回传迁移提示头。
	if strings.Contains(r.URL.Path, "/public/forum/") {
		// 调用w.Header完成当前处理。
		w.Header().Set("Deprecation", "true")
		// 调用w.Header完成当前处理。
		w.Header().Set("Sunset", "Thu, 31 Dec 2026 23:59:59 GMT")
		// 更新当前变量或字段值。
		w.Header().Set("Link", "</api/v1/public/user/users/:id/history-topics>; rel=\"successor-version\"")
	}

	// 从路径中解析用户ID。
	userID, ok := mustPathID(w, r, "users")
	// 判断条件并进入对应分支逻辑。
	if !ok {
		// 返回当前处理结果。
		return
	}

	// 读取查询参数：limit/issue/year。
	limit := parseIntOrDefault(r.URL.Query().Get("limit"), 30)
	// 定义并初始化当前变量。
	issue := strings.TrimSpace(r.URL.Query().Get("issue"))
	// 定义并初始化当前变量。
	yearRaw := strings.TrimSpace(r.URL.Query().Get("year"))
	// 定义并初始化当前变量。
	year := parseIntOrDefault(yearRaw, 0)
	// 判断条件并进入对应分支逻辑。
	if yearRaw == "" {
		// 更新当前变量或字段值。
		year = 0
	}

	// 转发到用户域 RPC。
	resp, err := h.svcCtx.User.ForumAuthorHistory(r.Context(), &tkv1.ForumAuthorHistoryRequest{
		// 处理当前语句逻辑。
		UserId: userID,
		// 调用int32完成当前处理。
		Limit: int32(limit),
		// 处理当前语句逻辑。
		Issue: issue,
		// 调用int32完成当前处理。
		Year: int32(year),
	})
	// 调用writeRPCReply完成当前处理。
	writeRPCReply(w, resp, err)
}
