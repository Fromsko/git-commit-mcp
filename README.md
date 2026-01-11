# Git Commit MCP Server

一个基于 Model Context Protocol (MCP) 的 Git 操作服务器，提供完整的 Git 工作流程支持，包括状态查看、规范提交信息生成和提交执行。

## 功能特性

- **`git_status`**: 获取 Git 仓库状态，显示所有变更文件（新增、修改、删除）
- **`generate_commit_message`**: 根据提交类型和描述生成符合规范的 Git 提交信息
- **`git_commit`**: 执行 git add 和 git commit，使用指定的提交信息
- **`list_commit_types`**: 获取所有支持的提交类型及其说明
- **`git_log`**: 查看最近的 Git 提交历史
- **`git_branch`**: 查看当前所在的 Git 分支

## 提交规范

服务器支持以下提交类型，每种类型都有对应的表情符号：

| Type     | Emoji | 说明                                   |
| -------- | ----- | -------------------------------------- |
| feat     | ✨     | 新增功能                               |
| fix      | 🐛     | 修复 Bug                               |
| docs     | 📝     | 文档变更                               |
| style    | 💄     | 代码格式（不影响代码运行的变动）       |
| refactor | ♻️     | 重构（既不是新增功能，也不是修复 Bug） |
| perf     | ⚡️     | 性能优化                               |
| test     | ✅     | 增加测试                               |
| chore    | 🔧     | 构建过程或辅助工具的变动               |
| build    | 📦     | 构建系统或外部依赖变动                 |
| ci       | 👷     | CI 配置文件和脚本变动                  |
| revert   | ⏪     | 回退代码                               |
| init     | 🎉     | 项目初始化                             |
| ui       | 🎨     | 更新 UI 和样式文件                     |
| config   | ⚙️     | 配置文件修改                           |
| merge    | 🔀     | 合并分支                               |

## 提交信息格式

```
<emoji> <type>: <简短描述>

详细描述：
- 变更点 1
- 变更点 2
- 变更点 3
```

## 安装和使用

1. 安装依赖:
```bash
go mod tidy
```

2. 编译项目:
```bash
go build -o git-commit-mcp main.go
```

3. 运行服务器:
```bash
go run main.go
# 或者
./git-commit-mcp
```

## MCP 客户端配置

在支持 MCP 的客户端中添加此服务器配置：

```json
{
  "mcpServers": {
    "git-commit-mcp": {
      "command": "go",
      "args": ["run", "/path/to/git-commit-mcp/main.go"]
    }
  }
}
```

## 使用流程

### 推荐的 Git 提交流程

1. **查看变更状态**
   ```
   调用 git_status 工具查看当前变更
   ```

2. **生成提交信息**
   ```
   调用 generate_commit_message 工具生成规范的提交信息
   ```

3. **确认并提交**
   ```
   调用 git_commit 工具执行提交
   ```

4. **查看提交历史**（可选）
   ```
   调用 git_log 工具查看最近的提交
   ```

### 工具使用示例

#### 1. 查看状态
```json
{
  "name": "git_status",
  "arguments": {
    "path": "."  // 可选，默认当前目录
  }
}
```

#### 2. 生成提交信息
```json
{
  "name": "generate_commit_message",
  "arguments": {
    "commit_type": "feat",
    "short_desc": "实现用户登录功能",
    "details": [
      "添加登录表单组件",
      "实现 JWT token 认证",
      "添加登录状态管理"
    ]
  }
}
```

#### 3. 执行提交
```json
{
  "name": "git_commit",
  "arguments": {
    "message": "✨ feat: 实现用户登录功能\n\n详细描述：\n- 添加登录表单组件\n- 实现 JWT token 认证\n- 添加登录状态管理",
    "path": "."  // 可选，默认当前目录
  }
}
```

## 开发

此项目使用 [Model Context Protocol Go SDK](https://github.com/modelcontextprotocol/go-sdk) 构建。

### 项目结构

- `main.go`: 主服务器文件，定义 MCP 工具和服务器逻辑
- `go.mod`: Go 模块文件
- `README.md`: 项目说明文档

### 特性

- ✅ 完整的 Git 工作流程支持
- ✅ 中文友好的输出格式
- ✅ 符合规范的提交信息生成
- ✅ 错误处理和状态检查
- ✅ 类型安全的参数定义
- ✅ 清晰的变更状态展示

### 技术栈

- Go 1.21+
- Model Context Protocol Go SDK v1.2.0
- Git 命令行工具

## 许可证

MIT License
