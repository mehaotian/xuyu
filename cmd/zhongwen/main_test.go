package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunWithoutArgumentsShowsHelp(t *testing.T) {
	output := run(nil)

	if output == "" {
		t.Fatal("没有参数时应该显示帮助信息")
	}

	if !strings.HasPrefix(output, "中文") {
		t.Fatalf("帮助信息开头不正确：%q", output)
	}
}

func TestRunVersion(t *testing.T) {
	if got := run([]string{"版本"}); got != version+"\n" {
		t.Fatalf("版本信息不正确：%q", got)
	}
}

func TestRunUnknownCommand(t *testing.T) {
	output := run([]string{"未知命令"})

	if output == "" {
		t.Fatal("未知命令应该给出提示")
	}

	if !strings.HasPrefix(output, "暂不支持命令") {
		t.Fatalf("未知命令提示不正确：%q", output)
	}
}

func TestRunSourceFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "最小赋值.中")
	if err := os.WriteFile(path, []byte("设置 名字 = 1\n"), 0o600); err != nil {
		t.Fatalf("创建测试源码失败：%v", err)
	}

	output := run([]string{"运行", path})
	if output != "环境[名字=1]\n" {
		t.Fatalf("运行输出不正确：%q", output)
	}
}

func TestRunWithoutPathReportsUsage(t *testing.T) {
	output := run([]string{"运行"})
	if !strings.HasPrefix(output, "运行命令需要一个源码文件路径") {
		t.Fatalf("缺少文件路径时提示不正确：%q", output)
	}
}
