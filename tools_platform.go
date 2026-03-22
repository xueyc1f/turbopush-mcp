package main

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func registerPlatformTools(s *server.MCPServer, c *Client) {
	s.AddTool(
		mcp.NewTool("list_platforms",
			mcp.WithDescription("获取所有支持的发布平台列表，包含平台ID、名称、支持的内容类型（文章/图文/视频）"),
			mcp.WithBoolean("enable", mcp.Description("仅返回已启用的平台")),
			mcp.WithBoolean("article", mcp.Description("筛选支持文章的平台")),
			mcp.WithBoolean("graph_text", mcp.Description("筛选支持图文的平台")),
			mcp.WithBoolean("video", mcp.Description("筛选支持视频的平台")),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			args := request.GetArguments()
			path := "/platform/list?"
			if v, ok := args["enable"].(bool); ok && v {
				path += "enable=true&"
			}
			if v, ok := args["article"].(bool); ok && v {
				path += "article=true&"
			}
			if v, ok := args["graph_text"].(bool); ok && v {
				path += "graph_text=true&"
			}
			if v, ok := args["video"].(bool); ok && v {
				path += "video=true&"
			}
			resp, err := c.Get(path)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			return jsonResult(resp), nil
		},
	)
}
