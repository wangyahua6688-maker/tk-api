package handler

import (
	"net/http"

	"github.com/wangyahua6688-maker/tk-common/utils/codes"
	"github.com/wangyahua6688-maker/tk-common/utils/httpresp"
	"github.com/wangyahua6688-maker/tk-common/utils/reqx"
)

// mustPathID 从 REST 路径中提取 ID（例如 /lottery-info/:id/detail）。
// 失败时直接返回标准错误响应，调用方无需重复处理。
func mustPathID(w http.ResponseWriter, r *http.Request, prefix string) (uint64, bool) {
	// 统一解析路径中的业务 ID。
	id, err := reqx.ParsePathID(r.URL.Path, prefix)
	// 判断条件并进入对应分支逻辑。
	if err != nil {
		// 参数不合法时返回业务错误码，前端可直接提示。
		httpresp.Fail(w, http.StatusBadRequest, codes.InvalidID, err.Error())
		// 返回当前处理结果。
		return 0, false
	}
	// 返回当前处理结果。
	return id, true
}

// parseIntOrDefault 将 query 参数转为整数，失败则回退默认值。
func parseIntOrDefault(raw string, fallback int) int {
	// 返回当前处理结果。
	return reqx.ParseIntOrDefault(raw, fallback)
}

// getDeviceID 提取设备标识：优先头部，其次 query 参数。
func getDeviceID(r *http.Request) string {
	// 返回当前处理结果。
	return reqx.DeviceID(r)
}

// getClientIP 提取客户端 IP，按代理链路优先级依次回退。
func getClientIP(r *http.Request) string {
	// 返回当前处理结果。
	return reqx.ClientIP(r)
}
