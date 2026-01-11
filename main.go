package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// 提交类型定义
type CommitType struct {
	Emoji string
	Name  string
	Desc  string
}

var CommitTypes = []CommitType{
	{Emoji: "✨", Name: "feat", Desc: "新增功能"},
	{Emoji: "🐛", Name: "fix", Desc: "修复 Bug"},
	{Emoji: "📝", Name: "docs", Desc: "文档变更"},
	{Emoji: "💄", Name: "style", Desc: "代码格式"},
	{Emoji: "♻️", Name: "refactor", Desc: "重构代码"},
	{Emoji: "⚡️", Name: "perf", Desc: "性能优化"},
	{Emoji: "✅", Name: "test", Desc: "增加测试"},
	{Emoji: "🔧", Name: "chore", Desc: "构建/工具变动"},
	{Emoji: "📦", Name: "build", Desc: "构建系统变动"},
	{Emoji: "👷", Name: "ci", Desc: "CI 配置变动"},
	{Emoji: "⏪", Name: "revert", Desc: "回退代码"},
	{Emoji: "🎉", Name: "init", Desc: "项目初始化"},
	{Emoji: "🎨", Name: "ui", Desc: "更新 UI 样式"},
	{Emoji: "⚙️", Name: "config", Desc: "配置文件修改"},
	{Emoji: "🔀", Name: "merge", Desc: "合并分支"},
}

// ============================================
// 工具参数定义
// ============================================

// PathParam Git 仓库路径参数
type PathParam struct {
	Path string `json:"path" jsonschema:"Git 仓库路径，默认为当前目录"`
}

// CommitMessageParam 提交信息参数
type CommitMessageParam struct {
	CommitType string   `json:"commit_type" jsonschema:"提交类型: feat/fix/docs/style/refactor/perf/test/chore/build/ci/revert/init/ui/config/merge"`
	ShortDesc  string   `json:"short_desc" jsonschema:"简短描述（不超过50字符）"`
	Details    []string `json:"details" jsonschema:"详细描述列表，每项一个变更点"`
}

// GitCommitParam Git 提交参数
type GitCommitParam struct {
	Message string `json:"message" jsonschema:"提交信息"`
	Path    string `json:"path,omitempty" jsonschema:"Git 仓库路径，默认为当前目录"`
}

// GitLogParam Git 日志参数
type GitLogParam struct {
	Count *int32 `json:"count,omitempty" jsonschema:"显示的提交数量，默认10条"`
	Path  string `json:"path,omitempty" jsonschema:"Git 仓库路径，默认为当前目录"`
}

// ============================================
// 工具实现
// ============================================

// GitStatusOutput Git 状态输出
type GitStatusOutput struct {
	Status string `json:"status" jsonschema:"Git 状态信息"`
}

// GitStatus 获取 Git 仓库状态
func GitStatus(ctx context.Context, req *mcp.CallToolRequest, param PathParam) (
	*mcp.CallToolResult,
	GitStatusOutput,
	error,
) {
	repoPath := param.Path
	if repoPath == "" {
		repoPath = "."
	}

	// 检查是否是 Git 仓库
	if !isGitRepo(repoPath) {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: "❌ 当前目录不是 Git 仓库"},
			},
		}, GitStatusOutput{Status: "Not a git repository"}, nil
	}

	// 执行 git status --porcelain
	cmd := exec.CommandContext(ctx, "git", "status", "--porcelain")
	cmd.Dir = repoPath
	output, err := cmd.Output()
	if err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: fmt.Sprintf("❌ 获取状态失败: %v", err)},
			},
		}, GitStatusOutput{Status: "Error"}, nil
	}

	outputStr := strings.TrimSpace(string(output))
	if outputStr == "" {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: "✅ 工作区干净，没有变更"},
			},
		}, GitStatusOutput{Status: "Clean"}, nil
	}

	// 解析状态并格式化输出
	lines := strings.Split(outputStr, "\n")
	var result strings.Builder
	result.WriteString("📊 变更导图：\n\n")

	for _, line := range lines {
		if len(line) < 3 {
			continue
		}

		status := line[:2]
		path := line[3:]

		var icon, statusStr string
		switch {
		case strings.Contains(status, "??"):
			icon, statusStr = "➕", "新增"
		case strings.Contains(status, "M "):
			icon, statusStr = "📝", "修改"
		case strings.Contains(status, " M"):
			icon, statusStr = "📝", "修改"
		case strings.Contains(status, "D "):
			icon, statusStr = "➖", "删除"
		case strings.Contains(status, " D"):
			icon, statusStr = "➖", "删除"
		case strings.Contains(status, "A "):
			icon, statusStr = "➕", "新增"
		default:
			continue
		}

		result.WriteString(fmt.Sprintf("%s %s %s\n", icon, statusStr, path))
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: result.String()},
		},
	}, GitStatusOutput{Status: result.String()}, nil
}

// CommitMessageOutput 提交信息输出
type CommitMessageOutput struct {
	Message string `json:"message" jsonschema:"生成的提交信息"`
}

// GenerateCommitMessage 生成符合规范的 Git 提交信息
func GenerateCommitMessage(ctx context.Context, req *mcp.CallToolRequest, param CommitMessageParam) (
	*mcp.CallToolResult,
	CommitMessageOutput,
	error,
) {
	var typeInfo *CommitType
	for i, t := range CommitTypes {
		if t.Name == param.CommitType {
			typeInfo = &CommitTypes[i]
			break
		}
	}
	if typeInfo == nil {
		typeInfo = &CommitTypes[0] // 默认使用 feat
	}

	var details strings.Builder
	for i, d := range param.Details {
		if i > 0 {
			details.WriteString("\n")
		}
		details.WriteString(fmt.Sprintf("- %s", d))
	}

	commitMsg := fmt.Sprintf("%s %s: %s\n\n详细描述：\n%s",
		typeInfo.Emoji, typeInfo.Name, param.ShortDesc, details.String())

	result := fmt.Sprintf("📝 生成的提交信息：\n\n```\n%s\n```", commitMsg)

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: result},
		},
	}, CommitMessageOutput{Message: commitMsg}, nil
}

// GitCommitOutput Git 提交输出
type GitCommitOutput struct {
	Result string `json:"result" jsonschema:"提交结果"`
}

// GitCommit 执行 Git 提交
func GitCommit(ctx context.Context, req *mcp.CallToolRequest, param GitCommitParam) (
	*mcp.CallToolResult,
	GitCommitOutput,
	error,
) {
	repoPath := param.Path
	if repoPath == "" {
		repoPath = "."
	}

	// git add .
	addCmd := exec.CommandContext(ctx, "git", "add", ".")
	addCmd.Dir = repoPath
	if output, err := addCmd.CombinedOutput(); err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: fmt.Sprintf("❌ git add 失败: %s", string(output))},
			},
		}, GitCommitOutput{Result: "Failed"}, nil
	}

	// git commit
	commitCmd := exec.CommandContext(ctx, "git", "commit", "-m", param.Message)
	commitCmd.Dir = repoPath
	if output, err := commitCmd.CombinedOutput(); err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: fmt.Sprintf("❌ git commit 失败: %s", string(output))},
			},
		}, GitCommitOutput{Result: "Failed"}, nil
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: "✅ 提交成功！\n\n💡 如需推送，请执行: git push"},
		},
	}, GitCommitOutput{Result: "Success"}, nil
}

// CommitTypesOutput 提交类型列表输出
type CommitTypesOutput struct {
	Types string `json:"types" jsonschema:"支持的提交类型列表"`
}

// ListCommitTypes 获取支持的提交类型列表
func ListCommitTypes(ctx context.Context, req *mcp.CallToolRequest, param struct{}) (
	*mcp.CallToolResult,
	CommitTypesOutput,
	error,
) {
	var result strings.Builder
	result.WriteString("📋 支持的提交类型：\n\n")
	result.WriteString("| Type | Emoji | 说明 |\n")
	result.WriteString("|------|-------|------|\n")

	for _, t := range CommitTypes {
		result.WriteString(fmt.Sprintf("| %s | %s | %s |\n", t.Name, t.Emoji, t.Desc))
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: result.String()},
		},
	}, CommitTypesOutput{Types: result.String()}, nil
}

// GitLogOutput Git 日志输出
type GitLogOutput struct {
	Log string `json:"log" jsonschema:"Git 提交历史"`
}

// GitLog 查看 Git 提交历史
func GitLog(ctx context.Context, req *mcp.CallToolRequest, param GitLogParam) (
	*mcp.CallToolResult,
	GitLogOutput,
	error,
) {
	repoPath := param.Path
	if repoPath == "" {
		repoPath = "."
	}

	n := "10"
	if param.Count != nil {
		n = fmt.Sprintf("%d", *param.Count)
	}

	cmd := exec.CommandContext(ctx, "git", "log", "--oneline", "-n", n)
	cmd.Dir = repoPath
	output, err := cmd.Output()
	if err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: fmt.Sprintf("❌ 获取日志失败: %v", err)},
			},
		}, GitLogOutput{Log: "Error"}, nil
	}

	result := fmt.Sprintf("📜 最近 %s 条提交：\n\n%s", n, string(output))

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: result},
		},
	}, GitLogOutput{Log: result}, nil
}

// GitBranchOutput Git 分支输出
type GitBranchOutput struct {
	Branch string `json:"branch" jsonschema:"当前分支"`
}

// GitBranch 查看当前分支
func GitBranch(ctx context.Context, req *mcp.CallToolRequest, param PathParam) (
	*mcp.CallToolResult,
	GitBranchOutput,
	error,
) {
	repoPath := param.Path
	if repoPath == "" {
		repoPath = "."
	}

	cmd := exec.CommandContext(ctx, "git", "branch", "--show-current")
	cmd.Dir = repoPath
	output, err := cmd.Output()
	if err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: fmt.Sprintf("❌ 获取分支失败: %v", err)},
			},
		}, GitBranchOutput{Branch: "Error"}, nil
	}

	branch := strings.TrimSpace(string(output))
	result := fmt.Sprintf("🌿 当前分支: %s", branch)

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: result},
		},
	}, GitBranchOutput{Branch: branch}, nil
}

// ============================================
// 辅助函数
// ============================================

// isGitRepo 检查是否是 Git 仓库
func isGitRepo(path string) bool {
	gitDir := filepath.Join(path, ".git")
	if stat, err := os.Stat(gitDir); err == nil {
		return stat.IsDir() || (stat.Mode().Perm()&0111 != 0) // 可能是 git file
	}
	return false
}

func main() {
	// 创建 MCP 服务器实例
	server := mcp.NewServer(&mcp.Implementation{
		Name:    "git-commit-mcp",
		Version: "v1.0.0",
	}, nil)

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

	// 启动服务器，通过 stdio 传输
	log.Println("Starting Git Commit MCP Server...")
	if err := server.Run(context.Background(), &mcp.StreamableServerTransport{}); err != nil {
		log.Fatal(err)
	}
}
