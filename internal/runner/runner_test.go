package runner

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRunFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "最小赋值.中")
	if err := os.WriteFile(path, []byte("设置 名字 = 1\n"), 0o600); err != nil {
		t.Fatalf("创建测试源码失败：%v", err)
	}

	environment, err := RunFile(path)
	if err != nil {
		t.Fatalf("运行文件失败：%v", err)
	}

	if got := environment.String(); got != "环境[名字=1]" {
		t.Fatalf("运行结果不正确：%q", got)
	}
}

func TestRunFileReportsMissingFile(t *testing.T) {
	_, err := RunFile(filepath.Join(t.TempDir(), "不存在.中"))
	if err == nil {
		t.Fatal("文件不存在时应该返回错误")
	}
}

func TestRunFileExecutesMultipleAssignments(t *testing.T) {
	path := filepath.Join(t.TempDir(), "两条赋值.中")
	if err := os.WriteFile(path, []byte("设置 名字 = 1\n设置 年龄 = 18\n"), 0o600); err != nil {
		t.Fatalf("创建测试源码失败：%v", err)
	}

	environment, err := RunFile(path)
	if err != nil {
		t.Fatalf("运行文件失败：%v", err)
	}

	if got := environment.String(); got != "环境[名字=1, 年龄=18]" {
		t.Fatalf("多条赋值结果不正确：%q", got)
	}
}

func TestRunFileExecutesArithmeticExpression(t *testing.T) {
	path := filepath.Join(t.TempDir(), "整数表达式.中")
	if err := os.WriteFile(path, []byte("设置 基数 = 10\n设置 总数 = (基数 + 2) * 3 - 4 / 2\n"), 0o600); err != nil {
		t.Fatalf("创建测试源码失败：%v", err)
	}

	environment, err := RunFile(path)
	if err != nil {
		t.Fatalf("运行表达式文件失败：%v", err)
	}

	if got := environment.String(); got != "环境[基数=10, 总数=34]" {
		t.Fatalf("表达式文件结果不正确：%q", got)
	}
}
