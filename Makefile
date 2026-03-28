# 声明所有目标为伪目标（不对应实际文件）
.PHONY: run run-bff run-business run-user run-local run-local-bff run-local-business run-local-user smoke-local tidy fmt test

# 运行主应用（默认启动 BFF 服务）
run: run-bff

# 启动 BFF (Backend for Frontend) 服务
# 使用生产环境配置文件 tk-api.yaml
run-bff:
	go run . -f etc/tk-api.yaml

# 启动业务服务
# 切换到 tk-business 目录并使用业务服务配置
run-business:
	cd ../tk-business && go run . -f etc/business.yaml

# 启动用户服务
# 切换到 tk-user 目录并使用用户服务配置
run-user:
	cd ../tk-user && go run . -f etc/user.yaml

# 运行本地开发环境主应用（默认启动本地 BFF 服务）
run-local: run-local-bff

# 启动本地 BFF 服务
# 使用本地环境配置文件 tk-api.local.yaml
run-local-bff:
	go run . -f etc/tk-api.local.yaml

# 启动本地业务服务
# 切换到 tk-business 目录并使用本地业务服务配置
run-local-business:
	cd ../tk-business && go run . -f etc/business.local.yaml

# 启动本地用户服务
# 切换到 tk-user 目录并使用本地用户服务配置
run-local-user:
	cd ../tk-user && go run . -f etc/user.local.yaml

# 执行本地环境的冒烟测试
# 调用 smoke_local.sh 脚本，传入本地服务器地址
smoke-local:
	./scripts/smoke_local.sh http://127.0.0.1:8088

# 整理 Go 模块依赖
# 移除未使用的依赖并更新 go.mod 文件
tidy:
	go mod tidy

# 格式化所有 Go 代码文件
fmt:
	gofmt -w ./

# 运行所有测试用例
test:
	go test ./...
