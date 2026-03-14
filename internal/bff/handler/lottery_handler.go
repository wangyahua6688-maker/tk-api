package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"tk-common/utils/codes"
	"tk-common/utils/httpresp"
	tkv1 "tk-proto/tk/v1"
)

// LotteryCards 彩种列表接口：按 category 返回彩种封面卡片。
func (h *PublicHandler) LotteryCards(w http.ResponseWriter, r *http.Request) {
	// 读取前端传入的分类标识（为空表示默认分类）。
	category := strings.TrimSpace(r.URL.Query().Get("category"))
	// 将 HTTP 请求映射为业务域 gRPC 请求。
	resp, err := h.svcCtx.Business.ListCards(r.Context(), &tkv1.ListCardsRequest{
		Category: category,
	})
	// 统一输出：成功/业务失败/网关失败都走同一个响应封装。
	writeRPCReply(w, resp, err)
}

// LotteryDashboard 开奖看板接口：用于首页开奖区/开奖现场页头部数据展示。
func (h *PublicHandler) LotteryDashboard(w http.ResponseWriter, r *http.Request) {
	// 从路径中提取 special lottery id。
	id, ok := mustPathID(w, r, "special-lotteries")
	if !ok {
		// mustPathID 内部已经输出错误响应，这里直接返回。
		return
	}
	// 调用业务服务获取开奖看板。
	resp, err := h.svcCtx.Business.LotteryDashboard(r.Context(), &tkv1.IDRequest{Id: id})
	writeRPCReply(w, resp, err)
}

// DrawHistory 开奖区历史开奖接口（按彩种维度）。
func (h *PublicHandler) DrawHistory(w http.ResponseWriter, r *http.Request) {
	// 提取彩种 ID。
	id, ok := mustPathID(w, r, "special-lotteries")
	if !ok {
		return
	}
	// 读取排序与展示参数。
	orderMode := strings.TrimSpace(r.URL.Query().Get("order_mode"))
	showFive := parseIntOrDefault(strings.TrimSpace(r.URL.Query().Get("show_five")), 1) == 1
	limit := parseIntOrDefault(strings.TrimSpace(r.URL.Query().Get("limit")), 80)
	resp, err := h.svcCtx.Business.DrawHistory(r.Context(), &tkv1.DrawHistoryRequest{
		SpecialLotteryId: id,
		OrderMode:        orderMode,
		ShowFive:         showFive,
		Limit:            int32(limit),
	})
	writeRPCReply(w, resp, err)
}

// DrawDetail 开奖区开奖详情接口（按开奖记录ID查询）。
func (h *PublicHandler) DrawDetail(w http.ResponseWriter, r *http.Request) {
	id, ok := mustPathID(w, r, "draw-records")
	if !ok {
		return
	}
	resp, err := h.svcCtx.Business.DrawDetail(r.Context(), &tkv1.IDRequest{Id: id})
	writeRPCReply(w, resp, err)
}

// LotteryDetail 彩票详情接口：返回大图、投票、评论分组、推荐图纸等详情页核心数据。
func (h *PublicHandler) LotteryDetail(w http.ResponseWriter, r *http.Request) {
	// 提取图纸 ID。
	id, ok := mustPathID(w, r, "lottery-info")
	if !ok {
		return
	}
	// 调用业务域详情聚合接口。
	resp, err := h.svcCtx.Business.LotteryDetail(r.Context(), &tkv1.IDRequest{Id: id})
	writeRPCReply(w, resp, err)
}

// LotteryHistory 历史开奖接口：返回当前彩种的多期开奖记录。
func (h *PublicHandler) LotteryHistory(w http.ResponseWriter, r *http.Request) {
	// 提取图纸 ID。
	id, ok := mustPathID(w, r, "lottery-info")
	if !ok {
		return
	}
	// 转发到业务域历史开奖接口。
	resp, err := h.svcCtx.Business.LotteryHistory(r.Context(), &tkv1.IDRequest{Id: id})
	writeRPCReply(w, resp, err)
}

// LotteryResults 结果接口：当前与详情接口数据结构保持一致，便于前端复用。
func (h *PublicHandler) LotteryResults(w http.ResponseWriter, r *http.Request) {
	// 提取图纸 ID。
	id, ok := mustPathID(w, r, "lottery-info")
	if !ok {
		return
	}
	// 业务域内部会复用详情逻辑。
	resp, err := h.svcCtx.Business.LotteryResults(r.Context(), &tkv1.IDRequest{Id: id})
	writeRPCReply(w, resp, err)
}

// VoteRecord 查询当前设备投票记录。
func (h *PublicHandler) VoteRecord(w http.ResponseWriter, r *http.Request) {
	// 提取图纸 ID。
	id, ok := mustPathID(w, r, "lottery-info")
	if !ok {
		return
	}
	// 采集客户端指纹信息，用于查询当前设备是否已投票。
	resp, err := h.svcCtx.Business.VoteRecord(r.Context(), &tkv1.VoteRecordRequest{
		LotteryInfoId: id,
		// DeviceID 优先从请求头读取，兼容 query 兜底。
		DeviceId: getDeviceID(r),
		// IP 用于风控与审计。
		ClientIp: getClientIP(r),
		// UA 参与指纹计算（DeviceID 缺失时）。
		UserAgent: strings.TrimSpace(r.UserAgent()),
	})
	writeRPCReply(w, resp, err)
}

// Vote 提交投票：入参只接受 option_id。
func (h *PublicHandler) Vote(w http.ResponseWriter, r *http.Request) {
	// 提取图纸 ID。
	id, ok := mustPathID(w, r, "lottery-info")
	if !ok {
		return
	}

	var reqBody struct {
		OptionID uint64 `json:"option_id"`
	}
	// 读取请求体并限制最大大小，避免异常大包攻击。
	body, err := io.ReadAll(io.LimitReader(r.Body, 1<<20))
	if err != nil {
		httpresp.Fail(w, http.StatusBadRequest, codes.InvalidRequestBody, "invalid request body")
		return
	}
	// option_id 必填且必须是正整数。
	if len(body) == 0 || json.Unmarshal(body, &reqBody) != nil || reqBody.OptionID == 0 {
		httpresp.Fail(w, http.StatusBadRequest, codes.OptionIDRequired, "option_id is required")
		return
	}

	// 将 HTTP 请求转换为 gRPC 投票请求。
	resp, err := h.svcCtx.Business.Vote(r.Context(), &tkv1.VoteRequest{
		LotteryInfoId: id,
		OptionId:      reqBody.OptionID,
		// 指纹参数交给业务域做限流与去重。
		DeviceId:  getDeviceID(r),
		ClientIp:  getClientIP(r),
		UserAgent: strings.TrimSpace(r.UserAgent()),
	})
	writeRPCReply(w, resp, err)
}
