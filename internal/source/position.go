// Package source 保存源码中的位置和范围。
//
// 这些类型会被 Token、AST 和错误信息共同使用。Lexer 负责计算行列位置，
// 其他阶段只保存和传递范围，不重复解释偏移量。
package source

import "fmt"

// Position 表示源码中的一个位置。
//
// Offset 使用字节偏移，Line 和 Column 从 1 开始。Offset 适合在原始
// UTF-8 字节中切片，Line 和 Column 适合展示给使用者。
type Position struct {
	Offset int
	Line   int
	Column int
}

// String 返回适合诊断和测试使用的短格式。
func (p Position) String() string {
	return fmt.Sprintf("%d:%d", p.Line, p.Column)
}

// Span 表示源码中的一段范围，结束位置采用半开区间。
//
// 例如从第 1 列开始、长度为 2 个字符的内容，可以表示为 Start=1:1、
// End=1:3。半开区间可以避免相邻 Token 共享结束位置时产生歧义。
type Span struct {
	Start Position
	End   Position
}

// String 返回范围的短格式。
func (s Span) String() string {
	return fmt.Sprintf("%s-%s", s.Start, s.End)
}
