// Package ast 定义源码经过语法分析后的最小结构。
//
// 当前支持程序、多条赋值语句，以及整数和变量引用两种表达式。
// Statement 和 Expression 的边界保持独立，后续可以在不改动现有节点
// 含义的前提下加入其他语句和表达式。
package ast

import (
	"fmt"
	"strings"

	"github.com/mehaotian/xuyu/internal/source"
)

// Node 是所有 AST 节点共同的接口。
type Node interface {
	Span() source.Span
}

// Statement 表示可以直接执行的一条语句。
type Statement interface {
	Node
	statementNode()
}

// Expression 表示可以产生一个值的结构。
type Expression interface {
	Node
	expressionNode()
}

// Program 是一个源码文件的最小表示。
type Program struct {
	Statements []Statement
}

// Span 返回整个程序的范围。
func (p Program) Span() source.Span {
	if len(p.Statements) == 0 {
		return source.Span{}
	}
	return source.Span{
		Start: p.Statements[0].Span().Start,
		End:   p.Statements[len(p.Statements)-1].Span().End,
	}
}

// String 返回便于测试和阅读的紧凑结构表示。
func (p Program) String() string {
	parts := make([]string, 0, len(p.Statements))
	for _, statement := range p.Statements {
		parts = append(parts, fmt.Sprint(statement))
	}
	return "程序[" + strings.Join(parts, ", ") + "]"
}

// Assignment 表示“设置 名字 = 值”形式的赋值语句。
type Assignment struct {
	Name  string
	Value Expression
	Range source.Span
}

func (Assignment) statementNode() {}

// Span 返回赋值语句的源码范围。
func (a Assignment) Span() source.Span {
	return a.Range
}

// String 返回赋值语句的紧凑结构表示。
func (a Assignment) String() string {
	return fmt.Sprintf("设置(%s, %s)", a.Name, a.Value)
}

// IntegerLiteral 表示一个整数值。
type IntegerLiteral struct {
	Raw   string
	Value int64
	Range source.Span
}

func (IntegerLiteral) expressionNode() {}

// Span 返回整数的源码范围。
func (i IntegerLiteral) Span() source.Span {
	return i.Range
}

// String 返回整数的结构表示。
func (i IntegerLiteral) String() string {
	return fmt.Sprintf("整数(%d)", i.Value)
}

// VariableReference 表示对已经存在的变量进行读取。
type VariableReference struct {
	Name  string
	Range source.Span
}

func (VariableReference) expressionNode() {}

// Span 返回变量引用的源码范围。
func (v VariableReference) Span() source.Span {
	return v.Range
}

// String 返回变量引用的结构表示。
func (v VariableReference) String() string {
	return fmt.Sprintf("变量(%s)", v.Name)
}
