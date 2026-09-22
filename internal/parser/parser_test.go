package parser

import (
	"testing"

	"github.com/mehaotian/xuyu/internal/ast"
	"github.com/mehaotian/xuyu/internal/lexer"
)

func TestParseAssignment(t *testing.T) {
	tokens, err := lexer.Lex("设置 名字 = 1")
	if err != nil {
		t.Fatalf("词法分析失败：%v", err)
	}

	program, err := Parse(tokens)
	if err != nil {
		t.Fatalf("语法分析失败：%v", err)
	}

	if got := program.String(); got != "程序[设置(名字, 整数(1))]" {
		t.Fatalf("AST 输出不正确：%q", got)
	}

	assignment, ok := program.Statements[0].(ast.Assignment)
	if !ok {
		t.Fatalf("第一条语句不是赋值语句：%T", program.Statements[0])
	}
	if assignment.Name != "名字" {
		t.Fatalf("变量名不正确：%q", assignment.Name)
	}
}

func TestParseMultipleAssignments(t *testing.T) {
	tokens, err := lexer.Lex("设置 名字 = 1\n设置 年龄 = 18")
	if err != nil {
		t.Fatalf("词法分析失败：%v", err)
	}

	program, err := Parse(tokens)
	if err != nil {
		t.Fatalf("语法分析失败：%v", err)
	}

	if len(program.Statements) != 2 {
		t.Fatalf("语句数量不正确：得到 %d，想要 2", len(program.Statements))
	}
	if got := program.String(); got != "程序[设置(名字, 整数(1)), 设置(年龄, 整数(18))]" {
		t.Fatalf("多条语句 AST 不正确：%q", got)
	}
}

func TestParseVariableReference(t *testing.T) {
	tokens, err := lexer.Lex("设置 总数 = 名字")
	if err != nil {
		t.Fatalf("词法分析失败：%v", err)
	}

	program, err := Parse(tokens)
	if err != nil {
		t.Fatalf("语法分析失败：%v", err)
	}

	if got := program.String(); got != "程序[设置(总数, 变量(名字))]" {
		t.Fatalf("变量引用 AST 不正确：%q", got)
	}
}

func TestParseArithmeticPrecedence(t *testing.T) {
	tokens, err := lexer.Lex("设置 结果 = 1 + 2 * 3 - 4 / 2")
	if err != nil {
		t.Fatalf("词法分析失败：%v", err)
	}

	program, err := Parse(tokens)
	if err != nil {
		t.Fatalf("语法分析失败：%v", err)
	}

	want := "程序[设置(结果, 二元(-, 二元(+, 整数(1), 二元(*, 整数(2), 整数(3))), 二元(/, 整数(4), 整数(2))))]"
	if got := program.String(); got != want {
		t.Fatalf("表达式优先级不正确：得到 %q，想要 %q", got, want)
	}
}

func TestParseParenthesizedExpression(t *testing.T) {
	tokens, err := lexer.Lex("设置 结果 = (1 + 2) * 3")
	if err != nil {
		t.Fatalf("词法分析失败：%v", err)
	}

	program, err := Parse(tokens)
	if err != nil {
		t.Fatalf("语法分析失败：%v", err)
	}

	want := "程序[设置(结果, 二元(*, 二元(+, 整数(1), 整数(2)), 整数(3)))]"
	if got := program.String(); got != want {
		t.Fatalf("括号表达式不正确：得到 %q，想要 %q", got, want)
	}
}

func TestParseIgnoresBlankLines(t *testing.T) {
	tokens, err := lexer.Lex("\n设置 名字 = 1\n\n")
	if err != nil {
		t.Fatalf("词法分析失败：%v", err)
	}

	program, err := Parse(tokens)
	if err != nil {
		t.Fatalf("空行不应该导致解析失败：%v", err)
	}
	if len(program.Statements) != 1 {
		t.Fatalf("空行被错误解析成语句：%d", len(program.Statements))
	}
}

func TestParseReportsMissingValue(t *testing.T) {
	tokens, err := lexer.Lex("设置 名字 =")
	if err != nil {
		t.Fatalf("词法分析失败：%v", err)
	}

	_, err = Parse(tokens)
	if err == nil {
		t.Fatal("缺少整数值应该返回错误")
	}
	if got := err.Error(); got != "语法错误：期望整数、变量或左括号，实际文件结束，位置 1:8" {
		t.Fatalf("错误信息不正确：%q", got)
	}
}

func TestParseReportsMissingRightParenthesis(t *testing.T) {
	tokens, err := lexer.Lex("设置 结果 = (1 + 2")
	if err != nil {
		t.Fatalf("词法分析失败：%v", err)
	}

	_, err = Parse(tokens)
	if err == nil {
		t.Fatal("缺少右括号应该返回错误")
	}
	if got := err.Error(); got != "语法错误：期望右括号，实际文件结束，位置 1:15" {
		t.Fatalf("右括号错误不正确：%q", got)
	}
}

func TestParseRejectsUnaryMinus(t *testing.T) {
	tokens, err := lexer.Lex("设置 数字 = -1")
	if err != nil {
		t.Fatalf("词法分析失败：%v", err)
	}

	_, err = Parse(tokens)
	if err == nil {
		t.Fatal("当前阶段不应该接受一元负号")
	}
	if got := err.Error(); got != "语法错误：期望整数、变量或左括号，实际减号，位置 1:9" {
		t.Fatalf("一元负号错误不正确：%q", got)
	}
}

func TestParseRejectsExtraToken(t *testing.T) {
	tokens, err := lexer.Lex("设置 名字 = 1 设置 年龄 = 18")
	if err != nil {
		t.Fatalf("词法分析失败：%v", err)
	}

	_, err = Parse(tokens)
	if err == nil {
		t.Fatal("没有换行分隔的语句应该返回错误")
	}
}
