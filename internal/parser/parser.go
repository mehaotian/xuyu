// Package parser 将 Token 组成 AST。
//
// 当前实现多条赋值语句，右值支持整数、变量和四则运算。加减、乘除和
// 括号分别由不同层级处理，保持运算优先级明确。
package parser

import (
	"fmt"
	"strconv"

	"github.com/mehaotian/xuyu/internal/ast"
	"github.com/mehaotian/xuyu/internal/source"
	"github.com/mehaotian/xuyu/internal/token"
)

// Error 表示一次带位置的语法错误。
type Error struct {
	Message string
	Span    source.Span
}

// Error 返回适合直接显示给使用者的错误文本。
func (e *Error) Error() string {
	return fmt.Sprintf("语法错误：%s，位置 %s", e.Message, e.Span.Start)
}

// Parse 将 Token 列表解析成一个程序。
//
// 输入可以包含多条由换行分隔的“设置 名字 = 值”语句，并以文件结束
// Token 收尾。
func Parse(tokens []token.Token) (ast.Program, error) {
	parser := Parser{tokens: tokens}
	parser.skipNewlines()
	if parser.current().Kind == token.KindEOF {
		return ast.Program{}, parser.errorAt(parser.current(), "至少需要一条赋值语句")
	}

	var statements []ast.Statement
	for parser.current().Kind != token.KindEOF {
		assignment, err := parser.parseAssignment()
		if err != nil {
			return ast.Program{}, err
		}
		statements = append(statements, assignment)

		if parser.current().Kind == token.KindEOF {
			break
		}
		if parser.current().Kind != token.KindNewline {
			return ast.Program{}, parser.errorAt(parser.current(), "两条语句之间需要换行")
		}
		parser.skipNewlines()
	}

	return ast.Program{Statements: statements}, nil
}

type Parser struct {
	tokens []token.Token
	index  int
}

func (p *Parser) parseAssignment() (ast.Assignment, error) {
	keyword, err := p.expect(token.KindKeyword, "关键字“设置”")
	if err != nil {
		return ast.Assignment{}, err
	}
	if keyword.Lexeme != "设置" {
		return ast.Assignment{}, p.errorAt(keyword, "关键字必须是“设置”")
	}

	name, err := p.expect(token.KindIdentifier, "标识符")
	if err != nil {
		return ast.Assignment{}, err
	}

	if _, err := p.expect(token.KindEqual, "等号"); err != nil {
		return ast.Assignment{}, err
	}

	value, err := p.parseExpression()
	if err != nil {
		return ast.Assignment{}, err
	}

	return ast.Assignment{
		Name:  name.Lexeme,
		Value: value,
		Range: source.Span{Start: keyword.Span.Start, End: value.Span().End},
	}, nil
}

func (p *Parser) parseExpression() (ast.Expression, error) {
	return p.parseAdditive()
}

func (p *Parser) parseAdditive() (ast.Expression, error) {
	left, err := p.parseMultiplicative()
	if err != nil {
		return nil, err
	}

	for {
		operator := p.current()
		if operator.Kind != token.KindPlus && operator.Kind != token.KindMinus {
			return left, nil
		}
		p.index++

		right, err := p.parseMultiplicative()
		if err != nil {
			return nil, err
		}
		left = ast.BinaryExpression{
			Left:     left,
			Operator: binaryOperator(operator.Kind),
			Right:    right,
			Range:    source.Span{Start: left.Span().Start, End: right.Span().End},
		}
	}
}

func (p *Parser) parseMultiplicative() (ast.Expression, error) {
	left, err := p.parsePrimary()
	if err != nil {
		return nil, err
	}

	for {
		operator := p.current()
		if operator.Kind != token.KindStar && operator.Kind != token.KindSlash {
			return left, nil
		}
		p.index++

		right, err := p.parsePrimary()
		if err != nil {
			return nil, err
		}
		left = ast.BinaryExpression{
			Left:     left,
			Operator: binaryOperator(operator.Kind),
			Right:    right,
			Range:    source.Span{Start: left.Span().Start, End: right.Span().End},
		}
	}
}

func (p *Parser) parsePrimary() (ast.Expression, error) {
	current := p.current()
	switch current.Kind {
	case token.KindInteger:
		p.index++
		value, err := strconv.ParseInt(current.Lexeme, 10, 64)
		if err != nil {
			return nil, p.errorAt(current, "整数超出可表示范围")
		}
		return ast.IntegerLiteral{
			Raw:   current.Lexeme,
			Value: value,
			Range: current.Span,
		}, nil
	case token.KindIdentifier:
		p.index++
		return ast.VariableReference{
			Name:  current.Lexeme,
			Range: current.Span,
		}, nil
	case token.KindLeftParen:
		p.index++
		expression, err := p.parseExpression()
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(token.KindRightParen, "右括号"); err != nil {
			return nil, err
		}
		return expression, nil
	default:
		return nil, p.errorAt(current, fmt.Sprintf("期望整数、变量或左括号，实际%s", current.Kind))
	}
}

func binaryOperator(kind token.Kind) ast.BinaryOperator {
	switch kind {
	case token.KindPlus:
		return ast.OperatorAdd
	case token.KindMinus:
		return ast.OperatorSubtract
	case token.KindStar:
		return ast.OperatorMultiply
	case token.KindSlash:
		return ast.OperatorDivide
	default:
		panic("不是二元运算符")
	}
}

func (p *Parser) expect(kind token.Kind, description string) (token.Token, error) {
	current := p.current()
	if current.Kind != kind {
		return token.Token{}, p.errorAt(current, fmt.Sprintf("期望%s，实际%s", description, current.Kind))
	}
	p.index++
	return current, nil
}

func (p *Parser) current() token.Token {
	if p.index >= len(p.tokens) {
		return token.Token{Kind: token.KindEOF}
	}
	return p.tokens[p.index]
}

func (p *Parser) skipNewlines() {
	for p.current().Kind == token.KindNewline {
		p.index++
	}
}

func (p *Parser) errorAt(current token.Token, message string) error {
	return &Error{Message: message, Span: current.Span}
}
