package token

import (
	"testing"

	"github.com/mehaotian/xuyu/internal/source"
)

func TestKindString(t *testing.T) {
	tests := []struct {
		kind Kind
		want string
	}{
		{kind: KindKeyword, want: "关键字"},
		{kind: KindIdentifier, want: "标识符"},
		{kind: KindInteger, want: "整数"},
		{kind: KindEqual, want: "等号"},
		{kind: KindEOF, want: "文件结束"},
	}

	for _, test := range tests {
		if got := test.kind.String(); got != test.want {
			t.Fatalf("类别 %d 的名称不正确：得到 %q，想要 %q", test.kind, got, test.want)
		}
	}
}

func TestTokenString(t *testing.T) {
	token := Token{
		Kind:   KindKeyword,
		Lexeme: "设置",
		Span: source.Span{
			Start: source.Position{Offset: 0, Line: 1, Column: 1},
			End:   source.Position{Offset: 6, Line: 1, Column: 3},
		},
	}

	if got := token.String(); got != "关键字(\"设置\")@1:1-1:3" {
		t.Fatalf("Token 格式不正确：%q", got)
	}
}
