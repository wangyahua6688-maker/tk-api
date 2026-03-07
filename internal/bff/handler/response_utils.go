package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"tk-api/internal/shared/httpresp"
	tkv1 "tk-proto/tk/v1"

	"github.com/zeromicro/go-zero/rest/httpx"
)

// writeRPCReply 统一处理 gRPC 返回并转成 HTTP 响应。
// 处理策略：
// 1. gRPC 调用异常 -> 502；
// 2. 业务码非 0 -> 200 + {code,msg}；
// 3. 业务码为 0 -> 解 data_json 并输出 {code:0,data}。
func writeRPCReply(w http.ResponseWriter, resp *tkv1.JsonDataReply, err error) {
	// gRPC 连接失败/超时等网络级错误。
	if err != nil {
		httpresp.Fail(w, http.StatusBadGateway, 50201, "upstream service unavailable")
		return
	}
	// 下游空响应兜底。
	if resp == nil {
		httpresp.Fail(w, http.StatusBadGateway, 50202, "empty upstream response")
		return
	}
	// 业务错误：保持 HTTP 200，仅透传业务码与业务文案。
	if resp.GetCode() != 0 {
		httpx.OkJson(w, httpresp.Envelope{
			Code: int(resp.GetCode()),
			Msg:  strings.TrimSpace(resp.GetMsg()),
		})
		return
	}

	// data_json 约定为 JSON 对象。
	data := map[string]interface{}{}
	raw := strings.TrimSpace(resp.GetDataJson())
	if raw != "" {
		// 下游返回格式异常时，明确抛网关级错误。
		if err := json.Unmarshal([]byte(raw), &data); err != nil {
			httpresp.Fail(w, http.StatusBadGateway, 50203, "invalid upstream payload")
			return
		}
	}
	// 标准成功响应。
	httpresp.OK(w, data)
}
