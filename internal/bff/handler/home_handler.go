package handler

import (
	"net/http"
	"strings"

	tkv1 "tk-proto/gen/go/tk/v1"
)

// HomeOverview 首页聚合接口：返回 banner、广播、切换彩种、外链等首页结构化数据。
func (h *PublicHandler) HomeOverview(w http.ResponseWriter, r *http.Request) {
	// 首页能力已并入 business 服务，BFF 直接调用业务域 RPC。
	resp, err := h.svcCtx.Business.HomeOverview(r.Context(), &tkv1.HomeOverviewRequest{})
	// 调用writeRPCReply完成当前处理。
	writeRPCReply(w, resp, err)
}

// LotteryCategories 图库分类搜索接口：支持 keyword 模糊筛选分类。
func (h *PublicHandler) LotteryCategories(w http.ResponseWriter, r *http.Request) {
	// 定义并初始化当前变量。
	keyword := strings.TrimSpace(r.URL.Query().Get("keyword"))
	// 定义并初始化当前变量。
	resp, err := h.svcCtx.Business.LotteryCategories(r.Context(), &tkv1.CategoryLibraryRequest{
		// 处理当前语句逻辑。
		Keyword: keyword,
	})
	// 调用writeRPCReply完成当前处理。
	writeRPCReply(w, resp, err)
}

// LiveScenePage 开奖现场聚合接口：一次返回整页渲染数据。
func (h *PublicHandler) LiveScenePage(w http.ResponseWriter, r *http.Request) {
	// 支持可选 query：special_lottery_id。
	sid := parseIntOrDefault(strings.TrimSpace(r.URL.Query().Get("special_lottery_id")), 0)
	// 定义并初始化当前变量。
	resp, err := h.svcCtx.Business.LiveScenePage(r.Context(), &tkv1.IDRequest{
		// 调用uint64完成当前处理。
		Id: uint64(sid),
	})
	// 调用writeRPCReply完成当前处理。
	writeRPCReply(w, resp, err)
}
