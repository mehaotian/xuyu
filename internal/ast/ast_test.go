package ast

import (
	"testing"

	"github.com/mehaotian/xuyu/internal/source"
)

func TestAssignmentString(t *testing.T) {
	assignment := Assignment{
		Name: "名字",
		Value: IntegerLiteral{
			Raw:   "1",
			Value: 1,
		},
	}

	if got := assignment.String(); got != "设置(名字, 整数(1))" {
		t.Fatalf("赋值结构输出不正确：%q", got)
	}
}

func TestProgramSpanCoversStatements(t *testing.T) {
	program := Program{
		Statements: []Statement{
			Assignment{
				Range: source.Span{
					Start: source.Position{Line: 1, Column: 1},
					End:   source.Position{Line: 1, Column: 8},
				},
			},
		},
	}

	if got := program.Span().String(); got != "1:1-1:8" {
		t.Fatalf("程序范围不正确：%q", got)
	}
}

func TestVariableReferenceString(t *testing.T) {
	reference := VariableReference{Name: "名字"}

	if got := reference.String(); got != "变量(名字)" {
		t.Fatalf("变量引用输出不正确：%q", got)
	}
}
