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

// ExpertBoards 高手推荐榜单接口，统一使用 /public/user/expert-boards 路径。
func (h *ExpertHandler) ExpertBoards(w http.ResponseWriter, r *http.Request) {
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
