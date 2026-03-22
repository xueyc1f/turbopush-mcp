package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Client struct {
	baseURL    string
	auth       string
	httpClient *http.Client
}

func NewClient(port, auth string) *Client {
	return &Client{
		baseURL: fmt.Sprintf("http://127.0.0.1:%s", port),
		auth:    auth,
		httpClient: &http.Client{
			Timeout: time.Minute * 10,
		},
	}
}

// apiResp turbo_push 统一响应格式
type apiResp struct {
	Code int             `json:"code"`
	Msg  string          `json:"msg"`
	Data json.RawMessage `json:"data"`
}

func (c *Client) do(method, path string, body any) (*apiResp, error) {
	var reqBody io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("marshal request body: %w", err)
		}
		reqBody = bytes.NewReader(data)
	}
	req, err := http.NewRequest(method, c.baseURL+path, reqBody)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Authorization", c.auth)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http request: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusUnauthorized {
		return nil, fmt.Errorf("认证失败，请检查 TURBO_PUSH_AUTH 环境变量")
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}
	var result apiResp
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("unmarshal response: %w", err)
	}
	if result.Code != 200 {
		return &result, fmt.Errorf("API 错误 [%d]: %s", result.Code, result.Msg)
	}
	return &result, nil
}

func (c *Client) Get(path string) (*apiResp, error) {
	return c.do(http.MethodGet, path, nil)
}

func (c *Client) Post(path string, body any) (*apiResp, error) {
	return c.do(http.MethodPost, path, body)
}

func (c *Client) Delete(path string) (*apiResp, error) {
	return c.do(http.MethodDelete, path, nil)
}

// SSE 相关

type sseEvent struct {
	Event string `json:"event"`
	Data  string `json:"data"`
}

// PostSSE 发送 POST 请求并消费 SSE 流，收集所有事件后返回汇总
func (c *Client) PostSSE(path string, body any) ([]sseEvent, error) {
	var reqBody io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("marshal request body: %w", err)
		}
		reqBody = bytes.NewReader(data)
	}
	req, err := http.NewRequest(http.MethodPost, c.baseURL+path, reqBody)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Authorization", c.auth)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "text/event-stream")

	client := &http.Client{Timeout: time.Hour}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http request: %w", err)
	}
	defer resp.Body.Close()

	var events []sseEvent
	buf := make([]byte, 0, 4096)
	tmp := make([]byte, 1024)
	var currentEvent, currentData string

	for {
		n, err := resp.Body.Read(tmp)
		if n > 0 {
			buf = append(buf, tmp[:n]...)
			// 按行解析 SSE
			for {
				idx := bytes.IndexByte(buf, '\n')
				if idx < 0 {
					break
				}
				line := string(buf[:idx])
				buf = buf[idx+1:]

				if line == "" {
					// 空行表示一个事件结束
					if currentEvent != "" || currentData != "" {
						events = append(events, sseEvent{Event: currentEvent, Data: currentData})
						if currentEvent == "finish" {
							return events, nil
						}
						currentEvent = ""
						currentData = ""
					}
					continue
				}
				if len(line) > 6 && line[:6] == "event:" {
					currentEvent = trimPrefix(line[6:])
				} else if len(line) > 5 && line[:5] == "data:" {
					currentData = trimPrefix(line[5:])
				}
			}
		}
		if err != nil {
			if err == io.EOF {
				// 流结束，返回已收集的事件
				if currentEvent != "" || currentData != "" {
					events = append(events, sseEvent{Event: currentEvent, Data: currentData})
				}
				return events, nil
			}
			return events, fmt.Errorf("read SSE stream: %w", err)
		}
	}
}

func trimPrefix(s string) string {
	if len(s) > 0 && s[0] == ' ' {
		return s[1:]
	}
	return s
}
