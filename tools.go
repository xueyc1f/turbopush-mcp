package main

import (
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func registerTools(s *server.MCPServer, c *Client) {
	registerPlatformTools(s, c)
	registerAccountTools(s, c)
	registerSettingTools(s, c)
	registerContentTools(s, c)
	registerPublishTools(s, c)
	registerRecordTools(s, c)
}

// jsonResult 将 apiResp.Data 包装为 MCP text result
func jsonResult(resp *apiResp) *mcp.CallToolResult {
	if resp == nil || resp.Data == nil {
		return mcp.NewToolResultText("{}")
	}
	return mcp.NewToolResultText(string(resp.Data))
}
