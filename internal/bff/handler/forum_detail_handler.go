package handler

import (
	"net/http"
	"strings"

	tkv1 "tk-proto/gen/go/tk/v1"
)

// TopicDetail 论坛帖子详情接口（新路径 + 兼容别名）。
// 路径策略：
// 1) 新路径：/public/user/topics/:id/detail；
// 2) 兼容别名：/public/forum/topics/:id/detail。
func (h *PublicHandler) TopicDetail(w http.ResponseWriter, r *http.Request) {
	// 命中兼容别名时回传迁移提示头，便于调用方灰度切换。
	if strings.Contains(r.URL.Path, "/public/forum/") {
		// 调用w.Header完成当前处理。
		w.Header().Set("Deprecation", "true")
		// 调用w.Header完成当前处理。
		w.Header().Set("Sunset", "Thu, 31 Dec 2026 23:59:59 GMT")
		// 更新当前变量或字段值。
		w.Header().Set("Link", "</api/v1/public/user/topics/:id/detail>; rel=\"successor-version\"")
	}

	// 统一解析帖子ID。
	postID, ok := mustPathID(w, r, "topics")
	// 判断条件并进入对应分支逻辑。
	if !ok {
		// 返回当前处理结果。
		return
	}

	// 转发用户域 RPC，获取详情聚合数据。
	resp, err := h.svcCtx.User.ForumTopicDetail(r.Context(), &tkv1.ForumTopicDetailRequest{
		// 处理当前语句逻辑。
		PostId: postID,
	})
	// 调用writeRPCReply完成当前处理。
	writeRPCReply(w, resp, err)
}
