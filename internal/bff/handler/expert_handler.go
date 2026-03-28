package handler

import (
	"context"
	"net/http"
	"strings"

	tkv1 "github.com/wangyahua6688-maker/tk-proto/gen/go/tk/v1"
	"google.golang.org/grpc"
)

// expertUserClient 定义高手榜模块依赖的用户域接口。
type expertUserClient interface {
	ExpertBoards(ctx context.Context, in *tkv1.ExpertBoardsRequest, opts ...grpc.CallOption) (*tkv1.JsonDataReply, error)
}

// ExpertHandler 负责高手榜相关接口。
type ExpertHandler struct {
	user expertUserClient
}

// NewExpertHandler 创建高手榜模块处理器。
func NewExpertHandler(user expertUserClient) *ExpertHandler {
	return &ExpertHandler{user: user}
}

// ExpertBoards 高手推荐榜单接口。
// 路径策略：
// 1) 主路径：/public/user/expert-boards；
// 2) 兼容别名：/public/forum/expert-boards。
func (h *ExpertHandler) ExpertBoards(w http.ResponseWriter, r *http.Request) {
	// 兼容别名命中时，回传迁移提示头，便于调用方逐步切换。
	if strings.Contains(r.URL.Path, "/public/forum/") {
		// 调用w.Header完成当前处理。
		w.Header().Set("Deprecation", "true")
		// 调用w.Header完成当前处理。
		w.Header().Set("Sunset", "Thu, 31 Dec 2026 23:59:59 GMT")
		// 更新当前变量或字段值。
		w.Header().Set("Link", "</api/v1/public/user/expert-boards>; rel=\"successor-version\"")
	}

	// 读取分页和彩种筛选参数。
	limit := parseIntOrDefault(r.URL.Query().Get("limit"), 10)
	// 定义并初始化当前变量。
	lotteryCode := strings.TrimSpace(r.URL.Query().Get("lottery_code"))

	// 转发到用户域 RPC。
	resp, err := h.user.ExpertBoards(r.Context(), &tkv1.ExpertBoardsRequest{
		// 调用int32完成当前处理。
		Limit: int32(limit),
		// 处理当前语句逻辑。
		LotteryCode: lotteryCode,
	})
	// 调用writeRPCReply完成当前处理。
	writeRPCReply(w, resp, err)
}
