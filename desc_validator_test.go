package main

import (
	"strings"
	"testing"
)

func TestParseTopics(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want []descSpan
	}{
		{
			name: "empty",
			in:   "",
			want: nil,
		},
		{
			name: "single topic",
			in:   "#每日发文#",
			want: []descSpan{{start: 0, end: len("#每日发文#"), value: "每日发文"}},
		},
		{
			name: "two adjacent topics",
			in:   "#a##b#",
			want: []descSpan{
				{start: 0, end: 3, value: "a"},
				{start: 3, end: 6, value: "b"},
			},
		},
		{
			name: "unclosed - no second hash",
			in:   "#unclosed",
			want: nil,
		},
		{
			name: "empty topic ##",
			in:   "##",
			want: nil,
		},
		{
			name: "with space inside breaks topic",
			in:   "#a b#",
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseTopics(tt.in)
			if !equalSpans(got, tt.want) {
				t.Errorf("parseTopics(%q) = %+v, want %+v", tt.in, got, tt.want)
			}
		})
	}
}

func TestParseAtUsers(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want []descSpan
	}{
		{
			name: "empty",
			in:   "",
			want: nil,
		},
		{
			name: "single user with trailing space",
			in:   "@luster ",
			want: []descSpan{{start: 0, end: 7, value: "luster"}},
		},
		{
			name: "user at end of text",
			in:   "@luster",
			want: []descSpan{{start: 0, end: 7, value: "luster"}},
		},
		{
			name: "lone @",
			in:   "@",
			want: nil,
		},
		{
			name: "@ followed by space",
			in:   "@ hi",
			want: nil,
		},
		{
			name: "two users separated by space",
			in:   "@a @b ",
			want: []descSpan{
				{start: 0, end: 2, value: "a"},
				{start: 3, end: 5, value: "b"},
			},
		},
		{
			name: "chinese user",
			in:   "@小红书 ",
			want: []descSpan{{start: 0, end: 1 + len("小红书"), value: "小红书"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseAtUsers(tt.in)
			if !equalSpans(got, tt.want) {
				t.Errorf("parseAtUsers(%q) = %+v, want %+v", tt.in, got, tt.want)
			}
		})
	}
}

func TestDetectSyntaxIssues(t *testing.T) {
	tests := []struct {
		name     string
		in       string
		wantKind []string
	}{
		{
			name:     "valid mix",
			in:       "#每日发文##解放生产力#@TurboPush @luster 内容",
			wantKind: nil,
		},
		{
			name:     "empty desc",
			in:       "",
			wantKind: nil,
		},
		{
			name:     "unclosed topic",
			in:       "#未闭合 后面文字",
			wantKind: []string{issueUnclosedTopic},
		},
		{
			name:     "empty topic",
			in:       "##",
			wantKind: []string{issueEmptyTopic},
		},
		{
			name:     "empty at user",
			in:       "@ 内容",
			wantKind: []string{issueEmptyAtUser},
		},
		{
			name:     "at user no space at end",
			in:       "前文 @luster",
			wantKind: []string{issueAtUserNoSpace},
		},
		{
			name:     "isolated hash with space is not issue",
			in:       "数学 # 表示集合",
			wantKind: nil,
		},
		{
			name:     "trailing hash alone is not issue",
			in:       "结尾 #",
			wantKind: nil,
		},
		{
			name:     "multiple issues",
			in:       "##  @ 内容 #未闭合 末尾@user",
			wantKind: []string{issueEmptyTopic, issueUnclosedTopic, issueEmptyAtUser, issueAtUserNoSpace},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := detectSyntaxIssues(tt.in)
			gotKinds := make([]string, len(got))
			for i, iss := range got {
				gotKinds[i] = iss.kind
			}
			if !equalStringSets(gotKinds, tt.wantKind) {
				t.Errorf("detectSyntaxIssues(%q) kinds = %v, want %v", tt.in, gotKinds, tt.wantKind)
			}
		})
	}
}

func TestValidateDescFormat(t *testing.T) {
	tests := []struct {
		name    string
		in      string
		wantErr bool
		wantSub string
	}{
		{
			name:    "valid",
			in:      "#话题# @user 描述",
			wantErr: false,
		},
		{
			name:    "empty desc",
			in:      "",
			wantErr: false,
		},
		{
			name:    "unclosed topic",
			in:      "#未闭合 文字",
			wantErr: true,
			wantSub: "未闭合的话题标签",
		},
		{
			name:    "at user no space at end",
			in:      "嗨 @luster",
			wantErr: true,
			wantSub: "@luster",
		},
		{
			name:    "empty at user",
			in:      "前 @ 后",
			wantErr: true,
			wantSub: "@ 后面请填写用户名",
		},
		{
			name:    "empty topic",
			in:      "## 内容",
			wantErr: true,
			wantSub: "话题标签中间不能为空",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateDescFormat(tt.in)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("validateDescFormat(%q) = nil, want error", tt.in)
				}
				if tt.wantSub != "" && !strings.Contains(err.Error(), tt.wantSub) {
					t.Errorf("validateDescFormat(%q) error %q does not contain %q", tt.in, err.Error(), tt.wantSub)
				}
				return
			}
			if err != nil {
				t.Errorf("validateDescFormat(%q) = %v, want nil", tt.in, err)
			}
		})
	}
}

func equalSpans(a, b []descSpan) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func equalStringSets(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
