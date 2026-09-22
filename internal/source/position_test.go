package source

import "testing"

func TestPositionString(t *testing.T) {
	position := Position{Offset: 4, Line: 2, Column: 5}

	if got := position.String(); got != "2:5" {
		t.Fatalf("位置格式不正确：%q", got)
	}
}

func TestSpanUsesHalfOpenEnd(t *testing.T) {
	span := Span{
		Start: Position{Offset: 0, Line: 1, Column: 1},
		End:   Position{Offset: 6, Line: 1, Column: 3},
	}

	if got := span.String(); got != "1:1-1:3" {
		t.Fatalf("范围格式不正确：%q", got)
	}
}
