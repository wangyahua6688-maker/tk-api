package handler

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/wangyahua6688-maker/tk-common/utils/codes"
	"github.com/wangyahua6688-maker/tk-common/utils/httpresp"
	tkv1 "github.com/wangyahua6688-maker/tk-proto/gen/go/tk/v1"
	"google.golang.org/grpc"
)

// lotteryBusinessClient 定义开奖模块依赖的业务域接口。
type lotteryBusinessClient interface {
	ListCards(ctx context.Context, in *tkv1.ListCardsRequest, opts ...grpc.CallOption) (*tkv1.JsonDataReply, error)
	LotteryDashboard(ctx context.Context, in *tkv1.IDRequest, opts ...grpc.CallOption) (*tkv1.LotteryDashboardReply, error)
	DrawHistory(ctx context.Context, in *tkv1.DrawHistoryRequest, opts ...grpc.CallOption) (*tkv1.LotteryHistoryReply, error)
	DrawDetail(ctx context.Context, in *tkv1.IDRequest, opts ...grpc.CallOption) (*tkv1.LotteryDrawDetailReply, error)
	LotteryDetail(ctx context.Context, in *tkv1.IDRequest, opts ...grpc.CallOption) (*tkv1.LotteryDetailReply, error)
	LotteryHistory(ctx context.Context, in *tkv1.IDRequest, opts ...grpc.CallOption) (*tkv1.LotteryHistoryReply, error)
	LotteryResults(ctx context.Context, in *tkv1.IDRequest, opts ...grpc.CallOption) (*tkv1.LotteryDetailReply, error)
	VoteRecord(ctx context.Context, in *tkv1.VoteRecordRequest, opts ...grpc.CallOption) (*tkv1.JsonDataReply, error)
	Vote(ctx context.Context, in *tkv1.VoteRequest, opts ...grpc.CallOption) (*tkv1.JsonDataReply, error)
}

// LotteryHandler 负责彩票详情、开奖、投票等接口。
type LotteryHandler struct {
	business lotteryBusinessClient
}

// NewLotteryHandler 创建彩票模块处理器。
func NewLotteryHandler(business lotteryBusinessClient) *LotteryHandler {
	return &LotteryHandler{business: business}
}

// LotteryCards 彩种列表接口：按 category 返回彩种封面卡片。
func (h *LotteryHandler) LotteryCards(w http.ResponseWriter, r *http.Request) {
	// 读取前端传入的分类标识（为空表示默认分类）。
	category := strings.TrimSpace(r.URL.Query().Get("category"))
	// 将 HTTP 请求映射为业务域 gRPC 请求。
	resp, err := h.business.ListCards(r.Context(), &tkv1.ListCardsRequest{
		// 处理当前语句逻辑。
		Category: category,
	})
	// 统一输出：成功/业务失败/网关失败都走同一个响应封装。
	writeRPCReply(w, resp, err)
}

// LotteryDashboard 开奖看板接口：用于首页开奖区/开奖现场页头部数据展示。
func (h *LotteryHandler) LotteryDashboard(w http.ResponseWriter, r *http.Request) {
	// 从路径中提取 special lottery id。
	id, ok := mustPathID(w, r, "special-lotteries")
	// 判断条件并进入对应分支逻辑。
	if !ok {
		// mustPathID 内部已经输出错误响应，这里直接返回。
		return
	}
	// 调用业务服务获取开奖看板。
	resp, err := h.business.LotteryDashboard(r.Context(), &tkv1.IDRequest{Id: id})
	if err != nil {
		writeTypedProtoReply(w, 0, "", nil, err)
		return
	}
	if resp == nil {
		httpresp.Fail(w, http.StatusBadGateway, codes.UpstreamEmptyReply, "empty upstream response")
		return
	}
	writeTypedProtoReply(w, resp.GetCode(), resp.GetMsg(), resp.GetData(), nil)
}

// DrawHistory 开奖区历史开奖接口（按彩种维度）。
func (h *LotteryHandler) DrawHistory(w http.ResponseWriter, r *http.Request) {
	// 提取彩种 ID。
	id, ok := mustPathID(w, r, "special-lotteries")
	// 判断条件并进入对应分支逻辑。
	if !ok {
		// 返回当前处理结果。
		return
	}
	// 读取排序与展示参数。
	orderMode := strings.TrimSpace(r.URL.Query().Get("order_mode"))
	// 定义并初始化当前变量。
	showFive := parseIntOrDefault(strings.TrimSpace(r.URL.Query().Get("show_five")), 1) == 1
	// 定义并初始化当前变量。
	limit := parseIntOrDefault(strings.TrimSpace(r.URL.Query().Get("limit")), 80)
	// 定义并初始化当前变量。
	resp, err := h.business.DrawHistory(r.Context(), &tkv1.DrawHistoryRequest{
		// 处理当前语句逻辑。
		SpecialLotteryId: id,
		// 处理当前语句逻辑。
		OrderMode: orderMode,
		// 处理当前语句逻辑。
		ShowFive: showFive,
		// 调用int32完成当前处理。
		Limit: int32(limit),
	})
	if err != nil {
		writeTypedProtoReply(w, 0, "", nil, err)
		return
	}
	if resp == nil {
		httpresp.Fail(w, http.StatusBadGateway, codes.UpstreamEmptyReply, "empty upstream response")
		return
	}
	writeTypedProtoReply(w, resp.GetCode(), resp.GetMsg(), resp.GetData(), nil)
}

// DrawDetail 开奖区开奖详情接口（按开奖记录ID查询）。
func (h *LotteryHandler) DrawDetail(w http.ResponseWriter, r *http.Request) {
	// 定义并初始化当前变量。
	id, ok := mustPathID(w, r, "draw-records")
	// 判断条件并进入对应分支逻辑。
	if !ok {
		// 返回当前处理结果。
		return
	}
	// 定义并初始化当前变量。
	resp, err := h.business.DrawDetail(r.Context(), &tkv1.IDRequest{Id: id})
	if err != nil {
		writeTypedProtoReply(w, 0, "", nil, err)
		return
	}
	if resp == nil {
		httpresp.Fail(w, http.StatusBadGateway, codes.UpstreamEmptyReply, "empty upstream response")
		return
	}
	writeTypedProtoReply(w, resp.GetCode(), resp.GetMsg(), resp.GetData(), nil)
}

// LotteryDetail 彩票详情接口：返回大图、投票、评论分组、推荐图纸等详情页核心数据。
func (h *LotteryHandler) LotteryDetail(w http.ResponseWriter, r *http.Request) {
	// 提取图纸 ID。
	id, ok := mustPathID(w, r, "lottery-info")
	// 判断条件并进入对应分支逻辑。
	if !ok {
		// 返回当前处理结果。
		return
	}
	// 调用业务域详情聚合接口。
	resp, err := h.business.LotteryDetail(r.Context(), &tkv1.IDRequest{Id: id})
	if err != nil {
		writeTypedProtoReply(w, 0, "", nil, err)
		return
	}
	if resp == nil {
		httpresp.Fail(w, http.StatusBadGateway, codes.UpstreamEmptyReply, "empty upstream response")
		return
	}
	writeTypedProtoReply(w, resp.GetCode(), resp.GetMsg(), resp.GetData(), nil)
}

// LotteryHistory 历史开奖接口：返回当前彩种的多期开奖记录。
func (h *LotteryHandler) LotteryHistory(w http.ResponseWriter, r *http.Request) {
	// 提取图纸 ID。
	id, ok := mustPathID(w, r, "lottery-info")
	// 判断条件并进入对应分支逻辑。
	if !ok {
		// 返回当前处理结果。
		return
	}
	// 转发到业务域历史开奖接口。
	resp, err := h.business.LotteryHistory(r.Context(), &tkv1.IDRequest{Id: id})
	if err != nil {
		writeTypedProtoReply(w, 0, "", nil, err)
		return
	}
	if resp == nil {
		httpresp.Fail(w, http.StatusBadGateway, codes.UpstreamEmptyReply, "empty upstream response")
		return
	}
	writeTypedProtoReply(w, resp.GetCode(), resp.GetMsg(), resp.GetData(), nil)
}

// LotteryResults 结果接口：当前与详情接口数据结构保持一致，便于前端复用。
func (h *LotteryHandler) LotteryResults(w http.ResponseWriter, r *http.Request) {
	// 提取图纸 ID。
	id, ok := mustPathID(w, r, "lottery-info")
	// 判断条件并进入对应分支逻辑。
	if !ok {
		// 返回当前处理结果。
		return
	}
	// 业务域内部会复用详情逻辑。
	resp, err := h.business.LotteryResults(r.Context(), &tkv1.IDRequest{Id: id})
	if err != nil {
		writeTypedProtoReply(w, 0, "", nil, err)
		return
	}
	if resp == nil {
		httpresp.Fail(w, http.StatusBadGateway, codes.UpstreamEmptyReply, "empty upstream response")
		return
	}
	writeTypedProtoReply(w, resp.GetCode(), resp.GetMsg(), resp.GetData(), nil)
}

// VoteRecord 查询当前设备投票记录。
func (h *LotteryHandler) VoteRecord(w http.ResponseWriter, r *http.Request) {
	// 提取图纸 ID。
	id, ok := mustPathID(w, r, "lottery-info")
	// 判断条件并进入对应分支逻辑。
	if !ok {
		// 返回当前处理结果。
		return
	}
	// 采集客户端指纹信息，用于查询当前设备是否已投票。
	resp, err := h.business.VoteRecord(r.Context(), &tkv1.VoteRecordRequest{
		// 处理当前语句逻辑。
		LotteryInfoId: id,
		// DeviceID 优先从请求头读取，兼容 query 兜底。
		DeviceId: getDeviceID(r),
		// IP 用于风控与审计。
		ClientIp: getClientIP(r),
		// UA 参与指纹计算（DeviceID 缺失时）。
		UserAgent: strings.TrimSpace(r.UserAgent()),
	})
	// 调用writeRPCReply完成当前处理。
	writeRPCReply(w, resp, err)
}

// Vote 提交投票：入参只接受 option_id。
func (h *LotteryHandler) Vote(w http.ResponseWriter, r *http.Request) {
	// 提取图纸 ID。
	id, ok := mustPathID(w, r, "lottery-info")
	// 判断条件并进入对应分支逻辑。
	if !ok {
		// 返回当前处理结果。
		return
	}

	// 声明当前变量。
	var reqBody struct {
		// 处理当前语句逻辑。
		OptionID uint64 `json:"option_id"`
	}
	// 读取请求体并限制最大大小，避免异常大包攻击。
	body, err := io.ReadAll(io.LimitReader(r.Body, 1<<20))
	// 判断条件并进入对应分支逻辑。
	if err != nil {
		// 调用httpresp.Fail完成当前处理。
		httpresp.Fail(w, http.StatusBadRequest, codes.InvalidRequestBody, "invalid request body")
		// 返回当前处理结果。
		return
	}
	// option_id 必填且必须是正整数。
	if len(body) == 0 || json.Unmarshal(body, &reqBody) != nil || reqBody.OptionID == 0 {
		// 调用httpresp.Fail完成当前处理。
		httpresp.Fail(w, http.StatusBadRequest, codes.OptionIDRequired, "option_id is required")
		// 返回当前处理结果。
		return
	}

	// 将 HTTP 请求转换为 gRPC 投票请求。
	resp, err := h.business.Vote(r.Context(), &tkv1.VoteRequest{
		// 处理当前语句逻辑。
		LotteryInfoId: id,
		// 处理当前语句逻辑。
		OptionId: reqBody.OptionID,
		// 指纹参数交给业务域做限流与去重。
		DeviceId: getDeviceID(r),
		// 调用getClientIP完成当前处理。
		ClientIp: getClientIP(r),
		// 调用strings.TrimSpace完成当前处理。
		UserAgent: strings.TrimSpace(r.UserAgent()),
	})
	// 调用writeRPCReply完成当前处理。
	writeRPCReply(w, resp, err)
}
