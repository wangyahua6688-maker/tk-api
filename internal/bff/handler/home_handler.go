package handler

import (
	"context"
	"net/http"
	"strings"

	tkv1 "github.com/wangyahua6688-maker/tk-proto/gen/go/tk/v1"
	"google.golang.org/grpc"
)

// homeBusinessClient 定义首页模块依赖的业务域接口。
type homeBusinessClient interface {
	HomeOverview(ctx context.Context, in *tkv1.HomeOverviewRequest, opts ...grpc.CallOption) (*tkv1.JsonDataReply, error)
	LiveScenePage(ctx context.Context, in *tkv1.IDRequest, opts ...grpc.CallOption) (*tkv1.JsonDataReply, error)
	LotteryCategories(ctx context.Context, in *tkv1.CategoryLibraryRequest, opts ...grpc.CallOption) (*tkv1.JsonDataReply, error)
}

// HomeHandler 负责首页与首页相关聚合接口。
type HomeHandler struct {
	business homeBusinessClient
}

// NewHomeHandler 创建首页模块处理器。
func NewHomeHandler(business homeBusinessClient) *HomeHandler {
	return &HomeHandler{business: business}
}

// HomeOverview 首页聚合接口：返回 banner、广播、切换彩种、外链等首页结构化数据。
func (h *HomeHandler) HomeOverview(w http.ResponseWriter, r *http.Request) {
	if isLegacyBusinessRoute(r.URL.Path) {
		markDeprecatedRoute(w, "/api/v1/public/business/home")
	}
	// 首页能力已并入 business 服务，BFF 直接调用业务域 RPC。
	resp, err := h.business.HomeOverview(r.Context(), &tkv1.HomeOverviewRequest{})
	// 调用writeRPCReply完成当前处理。
	writeRPCReply(w, resp, err)
}

// LotteryCategories 图库分类搜索接口：支持 keyword 模糊筛选分类。
func (h *HomeHandler) LotteryCategories(w http.ResponseWriter, r *http.Request) {
	if isLegacyBusinessRoute(r.URL.Path) {
		markDeprecatedRoute(w, "/api/v1/public/business/lottery-categories")
	}
	// 定义并初始化当前变量。
	keyword := strings.TrimSpace(r.URL.Query().Get("keyword"))
	// 定义并初始化当前变量。
	resp, err := h.business.LotteryCategories(r.Context(), &tkv1.CategoryLibraryRequest{
		// 处理当前语句逻辑。
		Keyword: keyword,
	})
	// 调用writeRPCReply完成当前处理。
	writeRPCReply(w, resp, err)
}

// LiveScenePage 开奖现场聚合接口：一次返回整页渲染数据。
func (h *HomeHandler) LiveScenePage(w http.ResponseWriter, r *http.Request) {
	if isLegacyBusinessRoute(r.URL.Path) {
		markDeprecatedRoute(w, "/api/v1/public/business/live-scene")
	}
	// 支持可选 query：special_lottery_id。
	sid := parseIntOrDefault(strings.TrimSpace(r.URL.Query().Get("special_lottery_id")), 0)
	// 定义并初始化当前变量。
	resp, err := h.business.LiveScenePage(r.Context(), &tkv1.IDRequest{
		// 调用uint64完成当前处理。
		Id: uint64(sid),
	})
	// 调用writeRPCReply完成当前处理。
	writeRPCReply(w, resp, err)
}
