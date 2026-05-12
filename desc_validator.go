package main

import (
	"fmt"
	"strings"
)

// desc 字段语法问题分类。
const (
	issueEmptyTopic    = "empty_topic"
	issueUnclosedTopic = "unclosed_topic"
	issueEmptyAtUser   = "empty_at_user"
	issueAtUserNoSpace = "at_user_no_space"
)

// syntaxIssue 表示 desc 中检测到的一处语法问题。
// position 是 UTF-8 字节偏移，便于定位但对非 ASCII 文本不等于字符序号。
type syntaxIssue struct {
	kind     string
	message  string
	position int
}

// descSpan 描述 desc 中一段被解析出的范围（半开区间 [start, end)）。
type descSpan struct {
	start int
	end   int
	value string
}

// parseTopics 解析 desc 中所有完整话题 #名称#。
// 完整定义：两个 # 之间包含至少一个字符，且不出现空格、换行或另一个 #。
//
// 实现使用字节扫描，依赖 UTF-8 多字节序列的后续字节恒 >= 0x80
// 与 ASCII 控制字符（'#'、' '、'\n'）不可能冲突，因此对中文等内容安全。
func parseTopics(text string) []descSpan {
	var out []descSpan
	n := len(text)
	for i := 0; i < n; {
		if text[i] != '#' {
			i++
			continue
		}
		j := i + 1
		valid := false
		for j < n {
			c := text[j]
			if c == '#' {
				if j > i+1 {
					valid = true
				}
				break
			}
			if c == ' ' || c == '\n' {
				break
			}
			j++
		}
		if valid {
			out = append(out, descSpan{start: i, end: j + 1, value: text[i+1 : j]})
			i = j + 1
			continue
		}
		i++
	}
	return out
}

// parseAtUsers 解析 desc 中所有 @用户名 片段。
// 用户名以 @ 起始，到空格、换行、下一个 @ 或下一个 # 之前结束；
// 文本末尾也视为结束（末尾未带空格的会被另行报告为 at_user_no_space）。
func parseAtUsers(text string) []descSpan {
	var out []descSpan
	n := len(text)
	for i := 0; i < n; {
		if text[i] != '@' {
			i++
			continue
		}
		j := i + 1
		for j < n {
			c := text[j]
			if c == ' ' || c == '\n' || c == '@' || c == '#' {
				break
			}
			j++
		}
		if j > i+1 {
			out = append(out, descSpan{start: i, end: j, value: text[i+1 : j]})
			i = j
			continue
		}
		i++
	}
	return out
}

// inSpans 判断 byte 偏移 pos 是否落在任一区间内。
func inSpans(pos int, spans []descSpan) bool {
	for _, s := range spans {
		if pos >= s.start && pos < s.end {
			return true
		}
	}
	return false
}

// detectSyntaxIssues 扫描 desc，返回所有可识别的话题与 @ 用户语法问题。
func detectSyntaxIssues(text string) []syntaxIssue {
	if text == "" {
		return nil
	}
	topics := parseTopics(text)
	users := parseAtUsers(text)

	var issues []syntaxIssue
	n := len(text)

	// 检测空话题与未闭合话题。
	for i := 0; i < n; {
		if text[i] != '#' || inSpans(i, topics) {
			i++
			continue
		}
		if i+1 >= n {
			// 末尾孤立 #，不视为问题（用户可能尚未输入完毕）。
			i++
			continue
		}
		switch text[i+1] {
		case ' ', '\n':
			// 孤立 #，跳过。
			i++
		case '#':
			issues = append(issues, syntaxIssue{
				kind:     issueEmptyTopic,
				message:  "话题标签中间不能为空，请填写话题名称",
				position: i,
			})
			i += 2
		default:
			issues = append(issues, syntaxIssue{
				kind:     issueUnclosedTopic,
				message:  "发现未闭合的话题标签，请使用 #话题名# 格式（前后各一个 #）",
				position: i,
			})
			i++
		}
	}

	// 检测空 @（@ 后无有效用户名）。
	for i := range n {
		if text[i] != '@' {
			continue
		}
		if inSpans(i, users) || inSpans(i, topics) {
			continue
		}
		issues = append(issues, syntaxIssue{
			kind:     issueEmptyAtUser,
			message:  "@ 后面请填写用户名，用户名后加空格结束，如 @小红书 ",
			position: i,
		})
	}

	// 检测 @用户名 位于文本末尾且未跟空格。
	for _, u := range users {
		if u.end == n {
			issues = append(issues, syntaxIssue{
				kind:     issueAtUserNoSpace,
				message:  fmt.Sprintf("@%s 后面请加一个空格，再输入其他内容或下一个 @用户", u.value),
				position: u.start,
			})
		}
	}

	return issues
}

// validateDescFormat 校验 desc 中话题（#话题名#）与提及用户（@用户名 ）的语法。
// 返回 nil 表示通过，否则返回包含所有问题的聚合 error。
func validateDescFormat(desc string) error {
	issues := detectSyntaxIssues(desc)
	if len(issues) == 0 {
		return nil
	}
	msgs := make([]string, 0, len(issues))
	for _, iss := range issues {
		msgs = append(msgs, fmt.Sprintf("[位置 %d] %s", iss.position, iss.message))
	}
	return fmt.Errorf("desc 格式校验失败:\n  - %s", strings.Join(msgs, "\n  - "))
}
