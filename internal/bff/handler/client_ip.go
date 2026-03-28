package handler

import (
	"context"
	"net"
	"net/http"
	"strings"

	"google.golang.org/grpc/metadata"
)

// clientIPKey 是注入 context 的客户端 IP 的 key 类型（防止 key 碰撞）。
type clientIPKey struct{}

const clientIPMetadataKey = "x-client-ip"

// withClientIP 将客户端真实 IP 写入 context，供下游 gRPC 调用链使用。
// 真实 IP 解析优先级：X-Forwarded-For → X-Real-IP → RemoteAddr
func withClientIP(r *http.Request) context.Context {
	ip := resolveClientIP(r)
	return context.WithValue(r.Context(), clientIPKey{}, ip)
}

// resolveClientIP 从请求中解析客户端真实 IP。
// 在反向代理（Nginx）后部署时，真实 IP 通常在 X-Forwarded-For 头的第一项。
func resolveClientIP(r *http.Request) string {
	// 1. X-Forwarded-For（可能包含多个，取第一个）
	if xff := strings.TrimSpace(r.Header.Get("X-Forwarded-For")); xff != "" {
		parts := strings.Split(xff, ",")
		if ip := strings.TrimSpace(parts[0]); isValidIP(ip) {
			return ip
		}
	}

	// 2. X-Real-IP（Nginx 常用配置）
	if xri := strings.TrimSpace(r.Header.Get("X-Real-IP")); isValidIP(xri) {
		return xri
	}

	// 3. 直连时的 RemoteAddr（格式为 ip:port）
	if host, _, err := net.SplitHostPort(r.RemoteAddr); err == nil {
		return host
	}
	return r.RemoteAddr
}

// isValidIP 判断字符串是否为合法 IP 地址（IPv4 或 IPv6）。
func isValidIP(ip string) bool {
	return net.ParseIP(ip) != nil
}

// ClientIPFromContext 从 context 中读取客户端 IP。
// 注：此方法供 gRPC 服务端（tk-user）从 context 中读取，
//
//	BFF 层在转发 gRPC 请求前先调用 withClientIP 将 IP 写入 context，
//	或通过 gRPC metadata 传递（见下方 WithClientIPMetadata）。
func ClientIPFromContext(ctx context.Context) string {
	if ip, ok := ctx.Value(clientIPKey{}).(string); ok {
		return ip
	}
	// 兼容通过字符串 key 注入的方式（tk-user repo 层读取）
	if ip, ok := ctx.Value("client_ip").(string); ok {
		return ip
	}
	return ""
}

// withClientIPOutgoingContext 将客户端 IP 注入 gRPC outgoing metadata。
func withClientIPOutgoingContext(ctx context.Context) context.Context {
	ip := strings.TrimSpace(ClientIPFromContext(ctx))
	if ip == "" {
		return ctx
	}
	return metadata.AppendToOutgoingContext(ctx, clientIPMetadataKey, ip)
}

// injectClientIPMiddleware 是一个 HTTP 中间件，将客户端 IP 注入请求 context。
// 使用方式：在 RegisterHandlers 中包装需要 IP 限流的 handler。
//
// 示例：
//
//	{Method: http.MethodPost, Path: "/api/v1/public/user/auth/sms-code",
//	    Handler: withIPContext(h.SendSMSCode)},
func withIPContext(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 将真实 IP 写入 context，使用约定的字符串 key 便于跨包读取
		ctx := context.WithValue(r.Context(), "client_ip", resolveClientIP(r))
		next(w, r.WithContext(ctx))
	}
}
