// Package token 定义词法分析器和后续语法分析器之间传递的最小数据。
package token

import (
	"fmt"

	"github.com/mehaotian/xuyu/internal/source"
)

// Kind 表示 Token 的类别。
//
// 当前只保留下一步实验需要的类别。新增类别时，应同时补充测试和语言
// 规范说明，不要为了预留未来功能一次性加入完整关键字表。
type Kind uint8

const (
	KindInvalid Kind = iota
	KindEOF
	KindKeyword
	KindIdentifier
	KindInteger
	KindEqual
	KindNewline
	KindPlus
	KindMinus
	KindStar
	KindSlash
	KindLeftParen
	KindRightParen
)

// String 返回 Token 类别的中文名称，便于调试输出和错误信息使用。
func (k Kind) String() string {
	switch k {
	case KindInvalid:
		return "无效"
	case KindEOF:
		return "文件结束"
	case KindKeyword:
		return "关键字"
	case KindIdentifier:
		return "标识符"
	case KindInteger:
		return "整数"
	case KindEqual:
		return "等号"
	case KindNewline:
		return "换行"
	case KindPlus:
		return "加号"
	case KindMinus:
		return "减号"
	case KindStar:
		return "乘号"
	case KindSlash:
		return "除号"
	case KindLeftParen:
		return "左括号"
	case KindRightParen:
		return "右括号"
	default:
		return "未知"
	}
}

// Token 是词法分析器产出的最小单元。
//
// Lexeme 保留源码中的原文，例如标识符“名字”或整数“123”。后续如果
// 需要 SymbolID 或数值转换，应在更明确的阶段完成，不在 Token 中猜测。
type Token struct {
	Kind   Kind
	Lexeme string
	Span   source.Span
}

// String 返回一条适合人阅读的调试表示。
func (t Token) String() string {
	return fmt.Sprintf("%s(%q)@%s", t.Kind, t.Lexeme, t.Span)
}
