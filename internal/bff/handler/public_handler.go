package handler

import "tk-api/internal/bff/svc"

// PublicHandler 封装所有对外公开接口处理器。
// 说明：
// - BFF 层只做入参解析与协议适配；
// - 真实业务计算在下游微服务完成；
// - 统一通过 writeRPCReply 输出标准响应结构。
type PublicHandler struct {
	// 处理当前语句逻辑。
	svcCtx *svc.ServiceContext
}

// NewPublicHandler 创建公共路由处理器实例。
func NewPublicHandler(svcCtx *svc.ServiceContext) *PublicHandler {
	// 返回当前处理结果。
	return &PublicHandler{svcCtx: svcCtx}
}
