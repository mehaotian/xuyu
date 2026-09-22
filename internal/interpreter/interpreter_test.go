package interpreter

import (
	"testing"

	"github.com/mehaotian/xuyu/internal/ast"
	"github.com/mehaotian/xuyu/internal/lexer"
	"github.com/mehaotian/xuyu/internal/parser"
)

func TestRunAssignment(t *testing.T) {
	program := parseSource(t, "设置 名字 = 1")

	environment, err := Run(program)
	if err != nil {
		t.Fatalf("执行失败：%v", err)
	}

	value, ok := environment.Get("名字")
	if !ok {
		t.Fatal("执行后没有找到变量“名字”")
	}
	if value != 1 {
		t.Fatalf("变量值不正确：得到 %d，想要 1", value)
	}

	if got := environment.String(); got != "环境[名字=1]" {
		t.Fatalf("环境输出不正确：%q", got)
	}
}

func TestGetUnknownVariable(t *testing.T) {
	environment := NewEnvironment()

	if _, ok := environment.Get("不存在"); ok {
		t.Fatal("空环境不应该包含未知变量")
	}
}

func TestRunVariableReference(t *testing.T) {
	program := parseSource(t, "设置 名字 = 1\n设置 总数 = 名字")

	environment, err := Run(program)
	if err != nil {
		t.Fatalf("执行变量引用失败：%v", err)
	}

	if got := environment.String(); got != "环境[名字=1, 总数=1]" {
		t.Fatalf("变量引用结果不正确：%q", got)
	}
}

func TestRunUndefinedVariableReportsError(t *testing.T) {
	program := parseSource(t, "设置 总数 = 未知")

	_, err := Run(program)
	if err == nil {
		t.Fatal("读取未定义变量应该返回错误")
	}
	if got := err.Error(); got != "执行错误：变量“未知”尚未定义，位置 1:9" {
		t.Fatalf("未定义变量错误不正确：%q", got)
	}
}

func parseSource(t *testing.T, input string) ast.Program {
	t.Helper()

	tokens, err := lexer.Lex(input)
	if err != nil {
		t.Fatalf("词法分析失败：%v", err)
	}

	program, err := parser.Parse(tokens)
	if err != nil {
		t.Fatalf("语法分析失败：%v", err)
	}

	return program
}
