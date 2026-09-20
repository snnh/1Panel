package ctl_conf

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const testConf = `#!/bin/bash
# 1Panel 命令行工具

BASE_DIR=/opt
LANGUAGE=zh
ORIGINAL_PORT="8080"
ORIGINAL_VERSION=v2.2.4
ORIGINAL_USERNAME="admin"
ORIGINAL_ENTRANCE=entrance
ORIGINAL_PASSWORD=
QUOTED_EMPTY=""
`

func newTestFile(t *testing.T, content string) string {
	t.Helper()
	filePath := filepath.Join(t.TempDir(), "1pctl")
	if err := os.WriteFile(filePath, []byte(content), 0o755); err != nil {
		t.Fatalf("write test conf failed, err: %v", err)
	}
	return filePath
}

func TestLoadFromFile(t *testing.T) {
	filePath := newTestFile(t, testConf)

	tests := []struct {
		key  string
		want string
	}{
		{key: "BASE_DIR", want: "/opt"},
		{key: "LANGUAGE", want: "zh"},
		{key: "ORIGINAL_PORT", want: "8080"},
		{key: "ORIGINAL_VERSION", want: "v2.2.4"},
		{key: "ORIGINAL_USERNAME", want: "admin"},
		{key: "ORIGINAL_ENTRANCE", want: "entrance"},
		{key: "ORIGINAL_PASSWORD", want: ""},
		{key: "QUOTED_EMPTY", want: `""`},
	}

	for _, tt := range tests {
		got, err := LoadFromFile(filePath, tt.key)
		if err != nil {
			t.Errorf("LoadFromFile(%s) returned an error: %v", tt.key, err)
			continue
		}
		if got != tt.want {
			t.Errorf("LoadFromFile(%s) = %q, want %q", tt.key, got, tt.want)
		}
	}
}

func TestLoadFromFileKeepsTrailingQuotesOnlyWhenPaired(t *testing.T) {
	filePath := newTestFile(t, strings.Join([]string{
		`HALF_QUOTED="admin`,
		`INNER_QUOTED=a"b`,
	}, "\n"))

	got, err := LoadFromFile(filePath, "HALF_QUOTED")
	if err != nil {
		t.Fatalf("LoadFromFile(HALF_QUOTED) returned an error: %v", err)
	}
	if got != `"admin` {
		t.Errorf("LoadFromFile(HALF_QUOTED) = %q, want %q", got, `"admin`)
	}

	got, err = LoadFromFile(filePath, "INNER_QUOTED")
	if err != nil {
		t.Fatalf("LoadFromFile(INNER_QUOTED) returned an error: %v", err)
	}
	if got != `a"b` {
		t.Errorf("LoadFromFile(INNER_QUOTED) = %q, want %q", got, `a"b`)
	}
}

func TestLoadFromFileMissingKey(t *testing.T) {
	filePath := newTestFile(t, testConf)

	if _, err := LoadFromFile(filePath, "NOT_EXIST"); err == nil {
		t.Error("LoadFromFile(NOT_EXIST) should return an error")
	}
	if _, err := LoadFromFile(filepath.Join(t.TempDir(), "not-exist"), "BASE_DIR"); err == nil {
		t.Error("LoadFromFile with a missing file should return an error")
	}
}
