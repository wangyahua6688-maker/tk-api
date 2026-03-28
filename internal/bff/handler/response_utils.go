package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/wangyahua6688-maker/tk-common/utils/codes"
	"github.com/wangyahua6688-maker/tk-common/utils/httpresp"
	tkv1 "github.com/wangyahua6688-maker/tk-proto/gen/go/tk/v1"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

// writeRPCReply 统一处理 gRPC 返回并转成 HTTP 响应。
// 处理策略：
// 1. gRPC 调用异常 -> 502；
// 2. 业务码非 0 -> 200 + {code,msg}；
// 3. 业务码为 0 -> 解 data_json 并输出 {code:0,data}。
func writeRPCReply(w http.ResponseWriter, resp *tkv1.JsonDataReply, err error) {
	// gRPC 连接失败/超时等网络级错误。
	if err != nil {
		// 调用httpresp.Fail完成当前处理。
		httpresp.Fail(w, http.StatusBadGateway, codes.UpstreamUnavailable, "upstream service unavailable")
		// 返回当前处理结果。
		return
	}
	// 下游空响应兜底。
	if resp == nil {
		// 调用httpresp.Fail完成当前处理。
		httpresp.Fail(w, http.StatusBadGateway, codes.UpstreamEmptyReply, "empty upstream response")
		// 返回当前处理结果。
		return
	}
	// 业务错误：保持 HTTP 200，仅透传业务码与业务文案。
	if resp.GetCode() != 0 {
		// 调用httpresp.BizFail完成当前处理。
		httpresp.BizFail(w, int(resp.GetCode()), strings.TrimSpace(resp.GetMsg()))
		// 返回当前处理结果。
		return
	}

	// data_json 约定为 JSON 对象。
	data := map[string]interface{}{}
	// 定义并初始化当前变量。
	raw := strings.TrimSpace(resp.GetDataJson())
	// 判断条件并进入对应分支逻辑。
	if raw != "" {
		// 下游返回格式异常时，明确抛网关级错误。
		if err := json.Unmarshal([]byte(raw), &data); err != nil {
			// 调用httpresp.Fail完成当前处理。
			httpresp.Fail(w, http.StatusBadGateway, codes.UpstreamBadPayload, "invalid upstream payload")
			// 返回当前处理结果。
			return
		}
	}
	// 标准成功响应。
	httpresp.OK(w, data)
}

// writeTypedProtoReply 统一处理 typed gRPC response 并转成 HTTP 响应。
func writeTypedProtoReply(w http.ResponseWriter, code int32, msg string, payload proto.Message, err error) {
	if err != nil {
		httpresp.Fail(w, http.StatusBadGateway, codes.UpstreamUnavailable, "upstream service unavailable")
		return
	}
	if code != 0 {
		httpresp.BizFail(w, int(code), strings.TrimSpace(msg))
		return
	}

	data := map[string]interface{}{}
	if payload != nil {
		raw, marshalErr := protojson.MarshalOptions{
			UseProtoNames:   true,
			EmitUnpopulated: true,
		}.Marshal(payload)
		if marshalErr != nil {
			httpresp.Fail(w, http.StatusBadGateway, codes.UpstreamBadPayload, "invalid upstream payload")
			return
		}
		if len(raw) > 0 && string(raw) != "null" {
			if err := json.Unmarshal(raw, &data); err != nil {
				httpresp.Fail(w, http.StatusBadGateway, codes.UpstreamBadPayload, "invalid upstream payload")
				return
			}
		}
	}

	httpresp.OK(w, data)
}
