#!/usr/bin/env bash
set -euo pipefail

# API 基础 URL，可通过第一个参数传入，默认为 http://127.0.0.1:8088
BASE_URL="${1:-http://127.0.0.1:8088}"

# 检查 jq 命令是否已安装，未安装则退出
if ! command -v jq >/dev/null 2>&1; then
  echo "jq is required to run smoke_local.sh" >&2
  exit 1
fi

# 存储上一次 HTTP 请求的响应体临时文件路径
LAST_BODY=""

# 发送 HTTP 请求并输出响应结果
# 参数:
#   name - 请求名称标识
#   method - HTTP 方法 (GET/POST 等)
#   path - 请求路径
#   body - 请求体数据 (可选，仅用于非 GET 请求)
# 输出:
#   打印请求信息和响应状态码
#   使用 jq 解析并输出响应的 code、msg 和 data 类型
# 副作用:
#   将响应内容保存到临时文件，并将文件路径存入 LAST_BODY 变量
hit() {
  local name="$1"
  local method="$2"
  local path="$3"
  local body="${4:-}"
  local tmp
  local http_code

  tmp="$(mktemp)"
  if [[ "$method" == "GET" ]]; then
    http_code="$(curl -sS -o "$tmp" -w '%{http_code}' "${BASE_URL}${path}")"
  else
    http_code="$(curl -sS -o "$tmp" -w '%{http_code}' -H 'Content-Type: application/json' -X "$method" -d "$body" "${BASE_URL}${path}")"
  fi

  printf '\n[%s] %s %s -> HTTP %s\n' "$name" "$method" "$path" "$http_code"
  jq '{code,msg,data_type:(.data|type)}' "$tmp"
  LAST_BODY="$tmp"
}

# 获取首页信息
hit home GET "/api/v1/public/home"
# 从首页响应中提取特别抽奖 ID
special_id="$(jq -r '.data.special_lotteries[0].id // 0' "$LAST_BODY")"

# 获取彩票分类列表
hit categories GET "/api/v1/public/lottery-categories"
# 获取九小类别的彩票卡片列表
hit cards GET "/api/v1/public/lottery-cards?category=jiuxiao"
# 从卡片响应中提取第一个彩票信息 ID
lottery_info_id="$(jq -r '.data.items[0].id // 0' "$LAST_BODY")"

# 获取特别抽奖仪表板数据
hit dashboard GET "/api/v1/public/special-lotteries/${special_id}/dashboard"
# 获取特别抽奖历史记录
hit draw_history GET "/api/v1/public/special-lotteries/${special_id}/history?limit=5"
# 从历史记录中提取第一个开奖记录 ID
draw_record_id="$(jq -r '.data.items[0].id // 0' "$LAST_BODY")"

# 获取直播场景信息
hit live_scene GET "/api/v1/public/live-scene?special_lottery_id=${special_id}"
# 获取彩票详细信息
hit lottery_detail GET "/api/v1/public/lottery-info/${lottery_info_id}/detail"
# 从彩票详情中提取当前开奖记录 ID
detail_draw_record_id="$(jq -r '.data.current.draw_record_id // 0' "$LAST_BODY")"
# 如果详情中存在有效的开奖记录 ID，则更新使用该 ID
if [[ "$detail_draw_record_id" != "0" ]]; then
  draw_record_id="$detail_draw_record_id"
fi

# 获取彩票历史记录
hit lottery_history GET "/api/v1/public/lottery-info/${lottery_info_id}/history"
# 获取彩票投票记录
hit vote_record GET "/api/v1/public/lottery-info/${lottery_info_id}/vote-record"
# 获取开奖记录详情
hit draw_detail GET "/api/v1/public/draw-records/${draw_record_id}/detail"

# 获取用户主题帖子列表
hit forum_topics GET "/api/v1/public/user/topics?limit=5"
# 从帖子列表中提取第一个帖子 ID
post_id="$(jq -r '.data.items[0].id // 0' "$LAST_BODY")"
# 从帖子列表中提取第一个帖子的用户 ID
user_id="$(jq -r '.data.items[0].user_id // 0' "$LAST_BODY")"

# 如果存在有效的帖子 ID，则获取帖子详情
if [[ "$post_id" != "0" ]]; then
  hit forum_detail GET "/api/v1/public/user/topics/${post_id}/detail"
fi

# 如果存在有效的用户 ID，则获取作者的历史帖子
if [[ "$user_id" != "0" ]]; then
  hit author_history GET "/api/v1/public/user/users/${user_id}/history-topics?limit=5"
fi

# 获取专家板块列表
hit expert GET "/api/v1/public/user/expert-boards?limit=5"
# 发送获取短信验证码的请求
hit sms_code POST "/api/v1/public/user/auth/sms-code" '{"phone":"13800138000","purpose":"login"}'
