package main

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func registerAccountTools(s *server.MCPServer, c *Client) {
	s.AddTool(
		mcp.NewTool("list_accounts",
			mcp.WithDescription("获取所有平台账号列表，包含账号ID、名称、所属平台、登录状态等信息"),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			resp, err := c.Get("/account/list")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			return jsonResult(resp), nil
		},
	)

	s.AddTool(
		mcp.NewTool("list_logged_accounts",
			mcp.WithDescription("获取所有已登录的平台账号列表，仅返回当前处于登录状态的账号"),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			resp, err := c.Get("/account/logged")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			return jsonResult(resp), nil
		},
	)
}
