package main

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func registerRecordTools(s *server.MCPServer, c *Client) {
	s.AddTool(
		mcp.NewTool("list_records",
			mcp.WithDescription("获取发布记录列表，可按状态和类型筛选"),
			mcp.WithNumber("status", mcp.Description("状态筛选：1=发布中，2=全部失败，3=部分成功，4=全部成功")),
			mcp.WithNumber("type", mcp.Description("类型筛选：1=文章，2=图文，3=视频")),
			mcp.WithNumber("current", mcp.Description("页码")),
			mcp.WithNumber("size", mcp.Description("每页条数")),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			args := request.GetArguments()
			path := "/record/list?"
			if v, ok := args["status"].(float64); ok {
				path += fmt.Sprintf("status=%d&", int(v))
			}
			if v, ok := args["type"].(float64); ok {
				path += fmt.Sprintf("type=%d&", int(v))
			}
			if v, ok := args["current"].(float64); ok {
				path += fmt.Sprintf("current=%d&", int(v))
			}
			if v, ok := args["size"].(float64); ok {
				path += fmt.Sprintf("size=%d&", int(v))
			}
			resp, err := c.Get(path)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			return jsonResult(resp), nil
		},
	)

	s.AddTool(
		mcp.NewTool("get_record_info",
			mcp.WithDescription("获取发布记录的详细信息，包括每个账号的发布结果、耗时、失败原因等"),
			mcp.WithNumber("record_id", mcp.Required(), mcp.Description("发布记录ID")),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			prid, err := request.RequireFloat("record_id")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			resp, err := c.Get(fmt.Sprintf("/record/info/%d", int(prid)))
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			return jsonResult(resp), nil
		},
	)
}
