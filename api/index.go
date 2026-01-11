package main

import (
	"context"
	"log"
	"net/http"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// MCPHandler Vercel API Route handler for MCP
func MCPHandler(w http.ResponseWriter, r *http.Request) {
	// 设置 CORS 头
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	// 创建 MCP 服务器
	server := mcp.NewServer(&mcp.Implementation{
		Name:    "git-commit-mcp",
		Version: "v1.0.0",
	}, nil)

	// 注册所有工具
	registerTools(server)

	// 创建流式传输
	transport := &mcp.StreamableServerTransport{
		SessionID: "git-commit-mcp-vercel", // Vercel 环境下的会话ID
		Stateless: true,                    // Vercel serverless 推荐无状态
	}

	// 连接服务器
	ctx := context.Background()
	_, err := server.Connect(ctx, transport, nil)
	if err != nil {
		log.Printf("Failed to connect server: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// 使用传输处理 HTTP 请求
	transport.ServeHTTP(w, r)
}

// 注册工具的函数（需要从 main.go 导入或复制）
func registerTools(server *mcp.Server) {
	// 这里需要复制 main.go 中的工具注册代码
	// 由于 Go 的包限制，可能需要重构代码结构
}
