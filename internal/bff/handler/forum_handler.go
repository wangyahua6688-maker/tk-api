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
		// 调用w.Header完成当前处理。
		w.Header().Set("Deprecation", "true")
		// 调用w.Header完成当前处理。
		w.Header().Set("Sunset", "Thu, 31 Dec 2026 23:59:59 GMT")
		// 更新当前变量或字段值。
		w.Header().Set("Link", "</api/v1/public/user/topics>; rel=\"successor-version\"")
	}

	// 定义并初始化当前变量。
	limit := parseIntOrDefault(r.URL.Query().Get("limit"), 20)
	// 定义并初始化当前变量。
	feed := strings.TrimSpace(r.URL.Query().Get("feed"))
	// 定义并初始化当前变量。
	keyword := strings.TrimSpace(r.URL.Query().Get("keyword"))
	// 定义并初始化当前变量。
	issue := strings.TrimSpace(r.URL.Query().Get("issue"))
	// 定义并初始化当前变量。
	year := parseIntOrDefault(r.URL.Query().Get("year"), 0)
	// 判断条件并进入对应分支逻辑。
	if strings.TrimSpace(r.URL.Query().Get("year")) == "" {
		// 更新当前变量或字段值。
		year = 0
	}
	// 定义并初始化当前变量。
	resp, err := h.svcCtx.User.ForumTopics(r.Context(), &tkv1.ForumTopicsRequest{
		// 调用int32完成当前处理。
		Limit: int32(limit),
		// 处理当前语句逻辑。
		Feed: feed,
		// 处理当前语句逻辑。
		Keyword: keyword,
		// 处理当前语句逻辑。
		Issue: issue,
		// 调用int32完成当前处理。
		Year: int32(year),
	})
	// 调用writeRPCReply完成当前处理。
	writeRPCReply(w, resp, err)
}
