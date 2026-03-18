package main

import (
	"tk-api/internal/bff"
)

// main 启动程序入口。
func main() {
	// tk-api 仅作为 BFF 入口：负责 HTTP 协议适配与下游 RPC 聚合。
	// 真实业务计算全部在 tk-business / tk-user 中完成。
	// 启动后会监听 etc/tk-api.yaml 中配置的 Host/Port。
	bff.Run()
}
