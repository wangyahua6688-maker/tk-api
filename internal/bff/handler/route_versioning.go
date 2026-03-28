package handler

import (
	"net/http"
	"strings"
)

const routeSunsetAt = "Thu, 31 Dec 2026 23:59:59 GMT"

// markDeprecatedRoute 为旧路由返回迁移提示头，便于调用方灰度切换到新规范路径。
func markDeprecatedRoute(w http.ResponseWriter, successorPath string) {
	w.Header().Set("Deprecation", "true")
	w.Header().Set("Sunset", routeSunsetAt)
	w.Header().Set("Link", "<"+successorPath+">; rel=\"successor-version\"")
}

// isLegacyBusinessRoute 判断当前请求是否仍命中旧的未带 business 标记的业务域路径。
func isLegacyBusinessRoute(path string) bool {
	return !strings.Contains(path, "/public/business/")
}
