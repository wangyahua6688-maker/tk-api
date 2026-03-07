package handler

import (
	"net/http"
	"strings"

	tkv1 "tk-proto/tk/v1"
)

// ExpertBoards 高手推荐榜单接口。
// 路径策略：
// 1) 主路径：/public/user/expert-boards；
// 2) 兼容别名：/public/forum/expert-boards。
func (h *PublicHandler) ExpertBoards(w http.ResponseWriter, r *http.Request) {
	// 兼容别名命中时，回传迁移提示头，便于调用方逐步切换。
	if strings.Contains(r.URL.Path, "/public/forum/") {
		w.Header().Set("Deprecation", "true")
		w.Header().Set("Sunset", "Thu, 31 Dec 2026 23:59:59 GMT")
		w.Header().Set("Link", "</api/v1/public/user/expert-boards>; rel=\"successor-version\"")
	}

	// 读取分页和彩种筛选参数。
	limit := parseIntOrDefault(r.URL.Query().Get("limit"), 10)
	lotteryCode := strings.TrimSpace(r.URL.Query().Get("lottery_code"))

	// 转发到用户域 RPC。
	resp, err := h.svcCtx.User.ExpertBoards(r.Context(), &tkv1.ExpertBoardsRequest{
		Limit:       int32(limit),
		LotteryCode: lotteryCode,
	})
	writeRPCReply(w, resp, err)
}
