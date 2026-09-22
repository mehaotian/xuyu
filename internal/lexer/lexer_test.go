package lexer

import (
	"testing"

	"github.com/mehaotian/xuyu/internal/token"
)

func TestLexesMinimalAssignment(t *testing.T) {
	tokens, err := Lex("设置 名字 = 123")
	if err != nil {
		t.Fatalf("词法分析失败：%v", err)
	}

	want := []struct {
		kind   token.Kind
		lexeme string
	}{
		{kind: token.KindKeyword, lexeme: "设置"},
		{kind: token.KindIdentifier, lexeme: "名字"},
		{kind: token.KindEqual, lexeme: "="},
		{kind: token.KindInteger, lexeme: "123"},
		{kind: token.KindEOF, lexeme: ""},
	}

	if len(tokens) != len(want) {
		t.Fatalf("Token 数量不正确：得到 %d，想要 %d", len(tokens), len(want))
	}

	for index, expected := range want {
		got := tokens[index]
		if got.Kind != expected.kind || got.Lexeme != expected.lexeme {
			t.Fatalf("第 %d 个 Token 不正确：得到 %s，想要 %s(%q)", index, got, expected.kind, expected.lexeme)
		}
	}
}

func TestLexesLineBreakAndUnicodePosition(t *testing.T) {
	tokens, err := Lex("设置\n名字")
	if err != nil {
		t.Fatalf("词法分析失败：%v", err)
	}

	if got := tokens[0].Span.String(); got != "1:1-1:3" {
		t.Fatalf("关键字范围不正确：%q", got)
	}
	if tokens[1].Kind != token.KindNewline {
		t.Fatalf("第二个 Token 应该是换行：%s", tokens[1])
	}
	if got := tokens[1].Span.String(); got != "1:3-2:1" {
		t.Fatalf("换行范围不正确：%q", got)
	}
	if got := tokens[2].Span.String(); got != "2:1-2:3" {
		t.Fatalf("第二行标识符范围不正确：%q", got)
	}
}

func TestLexesArithmeticOperators(t *testing.T) {
	tokens, err := Lex("设置 总数 = (名字 + 2) * 3 - 4 / 2")
	if err != nil {
		t.Fatalf("词法分析失败：%v", err)
	}

	want := []struct {
		kind   token.Kind
		lexeme string
	}{
		{kind: token.KindKeyword, lexeme: "设置"},
		{kind: token.KindIdentifier, lexeme: "总数"},
		{kind: token.KindEqual, lexeme: "="},
		{kind: token.KindLeftParen, lexeme: "("},
		{kind: token.KindIdentifier, lexeme: "名字"},
		{kind: token.KindPlus, lexeme: "+"},
		{kind: token.KindInteger, lexeme: "2"},
		{kind: token.KindRightParen, lexeme: ")"},
		{kind: token.KindStar, lexeme: "*"},
		{kind: token.KindInteger, lexeme: "3"},
		{kind: token.KindMinus, lexeme: "-"},
		{kind: token.KindInteger, lexeme: "4"},
		{kind: token.KindSlash, lexeme: "/"},
		{kind: token.KindInteger, lexeme: "2"},
		{kind: token.KindEOF, lexeme: ""},
	}

	if len(tokens) != len(want) {
		t.Fatalf("Token 数量不正确：得到 %d，想要 %d", len(tokens), len(want))
	}
	for index, expected := range want {
		got := tokens[index]
		if got.Kind != expected.kind || got.Lexeme != expected.lexeme {
			t.Fatalf("第 %d 个 Token 不正确：得到 %s，想要 %s(%q)", index, got, expected.kind, expected.lexeme)
		}
	}
}

func TestKeywordMustMatchTheWholeIdentifier(t *testing.T) {
	tokens, err := Lex("设置值")
	if err != nil {
		t.Fatalf("词法分析失败：%v", err)
	}

	if tokens[0].Kind != token.KindIdentifier {
		t.Fatalf("连续标识符不应该被拆成关键字：%s", tokens[0])
	}
}

func TestLexesASCIIIdentifierAndUnderscore(t *testing.T) {
	tokens, err := Lex("player_1")
	if err != nil {
		t.Fatalf("词法分析失败：%v", err)
	}

	if tokens[0].Kind != token.KindIdentifier || tokens[0].Lexeme != "player_1" {
		t.Fatalf("英文标识符不正确：%s", tokens[0])
	}
}

func TestUnknownCharacterReportsPosition(t *testing.T) {
	_, err := Lex("设置 @")
	if err == nil {
		t.Fatal("无法识别的字符应该返回错误")
	}

	if got := err.Error(); got != "无法识别的字符 \"@\"，位置 1:4" {
		t.Fatalf("错误信息不正确：%q", got)
	}
}
