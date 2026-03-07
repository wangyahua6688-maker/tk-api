# tk-api（BFF 网关）

`tk-api` 现在只负责客户端 HTTP 接口（BFF），内部通过 gRPC 调用同级目录的两个领域服务：

- `../tk-business`：业务域（首页聚合、分类库、开奖、详情、投票、开奖现场聚合）
- `../tk-user`：用户域（帖子、评论分组，后续可扩展登录与用户中心）

## 使用框架

- BFF：`go-zero`（`rest` + `zrpc`）
- gRPC 协议：`tk-proto/proto/tk/v1/*.proto`
- 共享模型：`tk-shared/models`

## 目录

```text
tk-api/
  main.go
  etc/tk-api.yaml
  internal/bff/
    config/
    svc/
    handler/
```

## 端口建议

- `tk-business`: `9102`
- `tk-user`: `9103`
- `tk-api(BFF)`: `8088`

## 启动顺序

1. `cd ../tk-user && go run . -f etc/user.yaml`
2. `cd ../tk-business && go run . -f etc/business.yaml`
3. `cd ../tk-api && go run . -f etc/tk-api.yaml`

## 协议与模型拆分说明

- `tk-proto` 独立的原因：
  - 一个仓库集中管理 gRPC 契约，避免不同服务各自维护 proto 导致版本漂移；
  - `tk-api/tk-business/tk-user` 都依赖同一份接口定义，升级与回滚更可控；
  - 后续如新增 `tk-admin-api` 或外部消费方，可直接复用协议模块。
- `tk-shared` 独立的原因：
  - `tk-business` 与 `tk-admin` 使用相同表模型，放公共模块可避免重复维护；
  - 字段变更只改一处，降低“模型不一致”引发的线上问题。

## 对外接口（保持 tk-web 兼容）

- `GET /health`
- `GET /api/v1/public/home`
- `GET /api/v1/public/live-scene`
- `GET /api/v1/public/lottery-categories`
- `GET /api/v1/public/user/topics`（主路径）
- `GET /api/v1/public/forum/topics`（兼容别名）
- `GET /api/v1/public/lottery-cards`
- `GET /api/v1/public/special-lotteries/:id/dashboard`
- `GET /api/v1/public/lottery-info/:id/detail`
- `GET /api/v1/public/lottery-info/:id/history`
- `GET /api/v1/public/lottery-info/:id/results`
- `GET /api/v1/public/lottery-info/:id/vote-record`
- `POST /api/v1/public/lottery-info/:id/vote`

## 数据表前缀

客户端相关表前缀已统一为 `tk_`。
