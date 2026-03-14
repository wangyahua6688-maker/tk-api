package handler

import (
	"net/http"
	"strings"

	tkv1 "tk-proto/tk/v1"
)

// TopicDetail 论坛帖子详情接口（新路径 + 兼容别名）。
// 路径策略：
// 1) 新路径：/public/user/topics/:id/detail；
// 2) 兼容别名：/public/forum/topics/:id/detail。
func (h *PublicHandler) TopicDetail(w http.ResponseWriter, r *http.Request) {
	// 命中兼容别名时回传迁移提示头，便于调用方灰度切换。
	if strings.Contains(r.URL.Path, "/public/forum/") {
		w.Header().Set("Deprecation", "true")
		w.Header().Set("Sunset", "Thu, 31 Dec 2026 23:59:59 GMT")
		w.Header().Set("Link", "</api/v1/public/user/topics/:id/detail>; rel=\"successor-version\"")
	}

	// 统一解析帖子ID。
	postID, ok := mustPathID(w, r, "topics")
	if !ok {
		return
	}

	// 转发用户域 RPC，获取详情聚合数据。
	resp, err := h.svcCtx.User.ForumTopicDetail(r.Context(), &tkv1.ForumTopicDetailRequest{
		PostId: postID,
	})
	writeRPCReply(w, resp, err)
}
