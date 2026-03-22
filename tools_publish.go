package main

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func registerPublishTools(s *server.MCPServer, c *Client) {
	s.AddTool(
		mcp.NewTool("publish_article",
			mcp.WithDescription(`发布文章到指定平台账号。需要先通过 create_article 创建文章内容，然后调用此工具发布。
postAccounts 数组中每个元素需要包含：id(账号ID)、platName(平台名)、settings(平台配置，必须包含 platType 字段)。
发布是同步操作，会等待所有账号发布完成后返回结果汇总。`),
			mcp.WithNumber("article_id", mcp.Required(), mcp.Description("文章ID（通过 create_article 创建后获得）")),
			mcp.WithBoolean("syncDraft", mcp.Description("是否仅同步草稿（不直接发布），默认 false")),
			mcp.WithBoolean("headless", mcp.Description("是否使用无头浏览器模式，默认 false")),
			mcp.WithArray("postAccounts", mcp.Required(), mcp.Description("发布目标账号数组，每个元素包含 id(账号ID)、platName(平台名)、settings(平台配置对象，含 platType)")),
		),
		publishHandler(c, "/sse/article/%d"),
	)

	s.AddTool(
		mcp.NewTool("publish_graph_text",
			mcp.WithDescription(`发布图文到指定平台账号。需要先通过 create_graph_text 创建图文内容，然后调用此工具发布。
参数说明同 publish_article。`),
			mcp.WithNumber("article_id", mcp.Required(), mcp.Description("图文ID（通过 create_graph_text 创建后获得）")),
			mcp.WithBoolean("syncDraft", mcp.Description("是否仅同步草稿")),
			mcp.WithBoolean("headless", mcp.Description("是否使用无头浏览器模式")),
			mcp.WithArray("postAccounts", mcp.Required(), mcp.Description("发布目标账号数组")),
		),
		publishHandler(c, "/sse/graphText/%d"),
	)

	s.AddTool(
		mcp.NewTool("publish_video",
			mcp.WithDescription(`发布视频到指定平台账号。需要先通过 create_video 创建视频内容，然后调用此工具发布。
参数说明同 publish_article。`),
			mcp.WithNumber("article_id", mcp.Required(), mcp.Description("视频ID（通过 create_video 创建后获得）")),
			mcp.WithBoolean("syncDraft", mcp.Description("是否仅同步草稿")),
			mcp.WithBoolean("headless", mcp.Description("是否使用无头浏览器模式")),
			mcp.WithArray("postAccounts", mcp.Required(), mcp.Description("发布目标账号数组")),
		),
		publishHandler(c, "/sse/video/%d"),
	)
}

func publishHandler(c *Client, pathTemplate string) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := request.GetArguments()
		rid, err := request.RequireFloat("article_id")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		body := map[string]any{
			"postAccounts": args["postAccounts"],
		}
		if v, ok := args["syncDraft"].(bool); ok {
			body["syncDraft"] = v
		}
		if v, ok := args["headless"].(bool); ok {
			body["headless"] = v
		}

		path := fmt.Sprintf(pathTemplate, int(rid))
		events, err := c.PostSSE(path, body)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("发布请求失败: %s", err)), nil
		}

		return mcp.NewToolResultText(summarizeEvents(events)), nil
	}
}

func summarizeEvents(events []sseEvent) string {
	var (
		successes []string
		errors    []string
		result    string
	)

	for _, ev := range events {
		switch ev.Event {
		case "success":
			successes = append(successes, ev.Data)
		case "error":
			errors = append(errors, ev.Data)
		case "finish":
			var finish struct {
				Msg string `json:"msg"`
				Res []bool `json:"res"`
			}
			if json.Unmarshal([]byte(ev.Data), &finish) == nil {
				total := len(finish.Res)
				successCount := 0
				for _, r := range finish.Res {
					if r {
						successCount++
					}
				}
				result = fmt.Sprintf("%s (成功: %d/%d)", finish.Msg, successCount, total)
			} else {
				result = ev.Data
			}
		case "wait":
			return fmt.Sprintf("发布等待: %s", ev.Data)
		case "vip":
			return fmt.Sprintf("会员限制: %s", ev.Data)
		}
	}

	var sb strings.Builder
	if result != "" {
		sb.WriteString("## 发布结果\n")
		sb.WriteString(result)
		sb.WriteString("\n\n")
	}
	if len(successes) > 0 {
		sb.WriteString("### 成功\n")
		for _, s := range successes {
			sb.WriteString("- ")
			sb.WriteString(s)
			sb.WriteString("\n")
		}
	}
	if len(errors) > 0 {
		sb.WriteString("\n### 错误\n")
		for _, e := range errors {
			sb.WriteString("- ")
			sb.WriteString(e)
			sb.WriteString("\n")
		}
	}
	if sb.Len() == 0 {
		return "发布完成，无详细事件信息"
	}
	return sb.String()
}
