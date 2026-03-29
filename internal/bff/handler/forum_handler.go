package handler

import (
	"context"
	"net/http"
	"strings"

	tkv1 "github.com/wangyahua6688-maker/tk-proto/gen/go/tk/v1"
	"google.golang.org/grpc"
)

// forumUserClient 定义论坛模块依赖的用户域接口。
type forumUserClient interface {
	ForumTopics(ctx context.Context, in *tkv1.ForumTopicsRequest, opts ...grpc.CallOption) (*tkv1.JsonDataReply, error)
	ForumTopicDetail(ctx context.Context, in *tkv1.ForumTopicDetailRequest, opts ...grpc.CallOption) (*tkv1.JsonDataReply, error)
	ForumAuthorHistory(ctx context.Context, in *tkv1.ForumAuthorHistoryRequest, opts ...grpc.CallOption) (*tkv1.JsonDataReply, error)
}

// ForumHandler 负责论坛列表、详情、作者历史等接口。
type ForumHandler struct {
	user forumUserClient
}

// NewForumHandler 创建论坛模块处理器。
func NewForumHandler(user forumUserClient) *ForumHandler {
	return &ForumHandler{user: user}
}

// TopicList 用户域帖子列表接口，统一使用 /public/user/topics 路径。
func (h *ForumHandler) TopicList(w http.ResponseWriter, r *http.Request) {
	// 定义并初始化当前变量。
	limit := parseIntOrDefault(r.URL.Query().Get("limit"), 20)
	// 定义并初始化当前变量。
	feed := strings.TrimSpace(r.URL.Query().Get("feed"))
	// 定义并初始化当前变量。
	keyword := strings.TrimSpace(r.URL.Query().Get("keyword"))
	// 定义并初始化当前变量。
	issue := strings.TrimSpace(r.URL.Query().Get("issue"))
	// 定义并初始化当前变量。
	year := parseIntOrDefault(r.URL.Query().Get("year"), 0)
	// 判断条件并进入对应分支逻辑。
	if strings.TrimSpace(r.URL.Query().Get("year")) == "" {
		// 更新当前变量或字段值。
		year = 0
	}
	// 定义并初始化当前变量。
	resp, err := h.user.ForumTopics(r.Context(), &tkv1.ForumTopicsRequest{
		// 调用int32完成当前处理。
		Limit: int32(limit),
		// 处理当前语句逻辑。
		Feed: feed,
		// 处理当前语句逻辑。
		Keyword: keyword,
		// 处理当前语句逻辑。
		Issue: issue,
		// 调用int32完成当前处理。
		Year: int32(year),
	})
	// 调用writeRPCReply完成当前处理。
	writeRPCReply(w, resp, err)
}

// TopicDetail 论坛帖子详情接口，统一使用 /public/user/topics/:id/detail 路径。
func (h *ForumHandler) TopicDetail(w http.ResponseWriter, r *http.Request) {
	// 统一解析帖子ID。
	postID, ok := mustPathID(w, r, "topics")
	// 判断条件并进入对应分支逻辑。
	if !ok {
		// 返回当前处理结果。
		return
	}

	// 转发用户域 RPC，获取详情聚合数据。
	resp, err := h.user.ForumTopicDetail(r.Context(), &tkv1.ForumTopicDetailRequest{
		// 处理当前语句逻辑。
		PostId: postID,
	})
	// 调用writeRPCReply完成当前处理。
	writeRPCReply(w, resp, err)
}

// AuthorHistory 论坛作者历史贴列表接口，统一使用 /public/user/users/:id/history-topics 路径。
func (h *ForumHandler) AuthorHistory(w http.ResponseWriter, r *http.Request) {
	// 从路径中解析用户ID。
	userID, ok := mustPathID(w, r, "users")
	// 判断条件并进入对应分支逻辑。
	if !ok {
		// 返回当前处理结果。
		return
	}

	// 读取查询参数：limit/issue/year。
	limit := parseIntOrDefault(r.URL.Query().Get("limit"), 30)
	// 定义并初始化当前变量。
	issue := strings.TrimSpace(r.URL.Query().Get("issue"))
	// 定义并初始化当前变量。
	yearRaw := strings.TrimSpace(r.URL.Query().Get("year"))
	// 定义并初始化当前变量。
	year := parseIntOrDefault(yearRaw, 0)
	// 判断条件并进入对应分支逻辑。
	if yearRaw == "" {
		// 更新当前变量或字段值。
		year = 0
	}

	// 转发到用户域 RPC。
	resp, err := h.user.ForumAuthorHistory(r.Context(), &tkv1.ForumAuthorHistoryRequest{
		// 处理当前语句逻辑。
		UserId: userID,
		// 调用int32完成当前处理。
		Limit: int32(limit),
		// 处理当前语句逻辑。
		Issue: issue,
		// 调用int32完成当前处理。
		Year: int32(year),
	})
	// 调用writeRPCReply完成当前处理。
	writeRPCReply(w, resp, err)
}
