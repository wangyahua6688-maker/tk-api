package handler

import (
	"fmt"
	"net"
	"net/http"
	"strconv"
	"strings"

	"tk-api/internal/shared/httpresp"
)

// mustPathID 从 REST 路径中提取 ID（例如 /lottery-info/:id/detail）。
// 失败时直接返回标准错误响应，调用方无需重复处理。
func mustPathID(w http.ResponseWriter, r *http.Request, prefix string) (uint64, bool) {
	// 统一解析路径中的业务 ID。
	id, err := parsePathID(r.URL.Path, prefix)
	if err != nil {
		// 参数不合法时返回业务错误码，前端可直接提示。
		httpresp.Fail(w, http.StatusBadRequest, 40011, err.Error())
		return 0, false
	}
	return id, true
}

// parsePathID 按给定路径段前缀读取紧随其后的数字 ID。
func parsePathID(path string, prefix string) (uint64, error) {
	// 以 "/" 分段遍历路径，寻找 prefix 后的下一个段作为 ID。
	parts := strings.Split(strings.Trim(path, "/"), "/")
	for idx := range parts {
		if parts[idx] == prefix && idx+1 < len(parts) {
			// 只接受正整数 ID。
			id, err := strconv.ParseUint(parts[idx+1], 10, 64)
			if err != nil || id == 0 {
				return 0, fmt.Errorf("invalid id")
			}
			return id, nil
		}
	}
	return 0, fmt.Errorf("invalid id")
}

// parseIntOrDefault 将 query 参数转为整数，失败则回退默认值。
func parseIntOrDefault(raw string, fallback int) int {
	// 清理空白后解析。
	v, err := strconv.Atoi(strings.TrimSpace(raw))
	if err != nil || v <= 0 {
		// 非法输入一律使用默认值，避免污染后端查询。
		return fallback
	}
	return v
}

// getDeviceID 提取设备标识：优先头部，其次 query 参数。
func getDeviceID(r *http.Request) string {
	// 先读请求头，便于 App/Web 统一传递稳定设备标识。
	deviceID := strings.TrimSpace(r.Header.Get("X-Device-ID"))
	if deviceID != "" {
		return deviceID
	}
	// 兼容旧调用：允许从 query 读取。
	return strings.TrimSpace(r.URL.Query().Get("device_id"))
}

// getClientIP 提取客户端 IP，按代理链路优先级依次回退。
func getClientIP(r *http.Request) string {
	// 代理场景优先读取 X-Forwarded-For 的首个真实源地址。
	if xff := strings.TrimSpace(r.Header.Get("X-Forwarded-For")); xff != "" {
		parts := strings.Split(xff, ",")
		if len(parts) > 0 && strings.TrimSpace(parts[0]) != "" {
			return strings.TrimSpace(parts[0])
		}
	}
	// 其次读取 X-Real-IP。
	if rip := strings.TrimSpace(r.Header.Get("X-Real-IP")); rip != "" {
		return rip
	}
	// 最后回退到连接层地址。
	host, _, err := net.SplitHostPort(strings.TrimSpace(r.RemoteAddr))
	if err == nil {
		return host
	}
	return strings.TrimSpace(r.RemoteAddr)
}
