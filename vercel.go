package main

import (
	"context"
	"log"
	"net/http"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// VercelHandler 创建一个兼容 Vercel 的 HTTP handler
func VercelHandler() http.Handler {
	// 创建 MCP 服务器
	server := mcp.NewServer(&mcp.Implementation{
		Name:    "git-commit-mcp",
		Version: "v1.0.0",
	}, nil)

	// 注册所有工具（复用 main.go 中的工具注册逻辑）
	registerTools(server)

	// 创建流式传输
	transport := &mcp.StreamableServerTransport{
		SessionID: "git-commit-mcp-session", // 可以使用随机生成的 ID
		Stateless: false,                    // Vercel 可以支持有状态
	}

	// 连接服务器
	ctx := context.Background()
	_, err := server.Connect(ctx, transport, nil)
	if err != nil {
		log.Fatalf("Failed to connect server: %v", err)
	}

	// 返回 HTTP handler
	return transport
}

// registerTools 注册所有工具（从 main.go 提取）
func registerTools(server *mcp.Server) {
	// 添加 git status 工具
	mcp.AddTool(server, &mcp.Tool{
		Name:        "git_status",
		Description: "获取 Git 仓库状态，显示所有变更文件（新增、修改、删除）",
	}, GitStatus)

	// 添加生成提交信息工具
	mcp.AddTool(server, &mcp.Tool{
		Name:        "generate_commit_message",
		Description: "根据提交类型和描述生成符合规范的 Git 提交信息",
	}, GenerateCommitMessage)

	// 添加 git commit 工具
	mcp.AddTool(server, &mcp.Tool{
		Name:        "git_commit",
		Description: "执行 git add 和 git commit，使用指定的提交信息",
	}, GitCommit)

	// 添加列出提交类型工具
	mcp.AddTool(server, &mcp.Tool{
		Name:        "list_commit_types",
		Description: "获取所有支持的提交类型及其说明",
	}, ListCommitTypes)

	// 添加 git log 工具
	mcp.AddTool(server, &mcp.Tool{
		Name:        "git_log",
		Description: "查看最近的 Git 提交历史",
	}, GitLog)

	// 添加 git branch 工具
	mcp.AddTool(server, &mcp.Tool{
		Name:        "git_branch",
		Description: "查看当前所在的 Git 分支",
	}, GitBranch)
}

// VercelResponse Vercel API 响应格式
type VercelResponse struct {
	Body       string            `json:"body"`
	Headers    map[string]string `json:"headers"`
	StatusCode int               `json:"statusCode"`
}

// APIHandler Vercel API Route handler
func APIHandler(w http.ResponseWriter, r *http.Request) {
	handler := VercelHandler()
	handler.ServeHTTP(w, r)
}

// 为了兼容 Vercel 的 serverless 函数，可以导出这个函数
func Handler(w http.ResponseWriter, r *http.Request) {
	APIHandler(w, r)
}
