// Package lexer 将源码切分为 Token。
//
// 当前实现覆盖空白、换行、关键字、标识符、整数、赋值号、四则运算符
// 和括号。Parser 负责组织结构，Interpreter 负责执行语义，Lexer 不在
// 这两层之间猜测含义。
package lexer

import (
	"fmt"
	"unicode"
	"unicode/utf8"

	"github.com/mehaotian/xuyu/internal/source"
	"github.com/mehaotian/xuyu/internal/token"
)

// Error 表示一次带源码范围的词法错误。
type Error struct {
	Message string
	Span    source.Span
}

// Error 返回适合直接显示给使用者的错误文本。
func (e *Error) Error() string {
	return fmt.Sprintf("%s，位置 %s", e.Message, e.Span.Start)
}

// Lexer 持有一段尚未扫描完的源码。
type Lexer struct {
	input  string
	offset int
	line   int
	column int
}

// Lex 扫描整段源码并返回 Token。
//
// 当前采用严格策略：遇到无法识别的字符立即返回错误，不生成猜测性的
// Token。返回的 Token 中包含文件结束标记，方便后续 Parser 统一处理。
func Lex(input string) ([]token.Token, error) {
	lexer := Lexer{
		input:  input,
		line:   1,
		column: 1,
	}

	var tokens []token.Token
	for {
		item, err := lexer.next()
		if err != nil {
			return tokens, err
		}

		tokens = append(tokens, item)
		if item.Kind == token.KindEOF {
			return tokens, nil
		}
	}
}

func (l *Lexer) next() (token.Token, error) {
	l.skipHorizontalWhitespace()

	start := l.position()
	if l.offset >= len(l.input) {
		return token.Token{Kind: token.KindEOF, Span: source.Span{Start: start, End: start}}, nil
	}

	r, size := utf8.DecodeRuneInString(l.input[l.offset:])
	if r == utf8.RuneError && size == 1 {
		return token.Token{}, l.errorAt(start, "源码包含无效的 UTF-8 字节")
	}

	switch {
	case r == '\n':
		l.advance()
		return token.Token{
			Kind:   token.KindNewline,
			Lexeme: "\n",
			Span:   source.Span{Start: start, End: l.position()},
		}, nil
	case isIdentifierStart(r):
		return l.readIdentifier(), nil
	case unicode.IsDigit(r):
		return l.readInteger(), nil
	case r == '=':
		return l.readSingleToken(start, r, token.KindEqual), nil
	case r == '+':
		return l.readSingleToken(start, r, token.KindPlus), nil
	case r == '-':
		return l.readSingleToken(start, r, token.KindMinus), nil
	case r == '*':
		return l.readSingleToken(start, r, token.KindStar), nil
	case r == '/':
		return l.readSingleToken(start, r, token.KindSlash), nil
	case r == '(':
		return l.readSingleToken(start, r, token.KindLeftParen), nil
	case r == ')':
		return l.readSingleToken(start, r, token.KindRightParen), nil
	default:
		l.advance()
		return token.Token{}, l.errorAt(start, fmt.Sprintf("无法识别的字符 %q", string(r)))
	}
}

func (l *Lexer) readSingleToken(start source.Position, r rune, kind token.Kind) token.Token {
	l.advance()
	return token.Token{
		Kind:   kind,
		Lexeme: string(r),
		Span:   source.Span{Start: start, End: l.position()},
	}
}

func (l *Lexer) readIdentifier() token.Token {
	start := l.position()
	startOffset := l.offset

	for l.offset < len(l.input) {
		r, size := utf8.DecodeRuneInString(l.input[l.offset:])
		if r == utf8.RuneError && size == 1 {
			break
		}
		if !isIdentifierPart(r) {
			break
		}
		l.advance()
	}

	lexeme := l.input[startOffset:l.offset]
	kind := token.KindIdentifier
	if lexeme == "设置" {
		kind = token.KindKeyword
	}

	return token.Token{
		Kind:   kind,
		Lexeme: lexeme,
		Span:   source.Span{Start: start, End: l.position()},
	}
}

func (l *Lexer) readInteger() token.Token {
	start := l.position()
	startOffset := l.offset

	for l.offset < len(l.input) {
		r, size := utf8.DecodeRuneInString(l.input[l.offset:])
		if r == utf8.RuneError && size == 1 {
			break
		}
		if !unicode.IsDigit(r) {
			break
		}
		l.advance()
	}

	return token.Token{
		Kind:   token.KindInteger,
		Lexeme: l.input[startOffset:l.offset],
		Span:   source.Span{Start: start, End: l.position()},
	}
}

func (l *Lexer) skipHorizontalWhitespace() {
	for l.offset < len(l.input) {
		r, size := utf8.DecodeRuneInString(l.input[l.offset:])
		if r == utf8.RuneError && size == 1 {
			return
		}
		if !unicode.IsSpace(r) || r == '\n' {
			return
		}
		l.advance()
	}
}

func (l *Lexer) advance() (rune, int) {
	r, size := utf8.DecodeRuneInString(l.input[l.offset:])
	l.offset += size
	if r == '\n' {
		l.line++
		l.column = 1
	} else {
		l.column++
	}
	return r, size
}

func (l *Lexer) position() source.Position {
	return source.Position{Offset: l.offset, Line: l.line, Column: l.column}
}

func (l *Lexer) errorAt(start source.Position, message string) error {
	return &Error{
		Message: message,
		Span:    source.Span{Start: start, End: l.position()},
	}
}

func isIdentifierStart(r rune) bool {
	return r == '_' || isChinese(r) || isASCIILetter(r)
}

func isIdentifierPart(r rune) bool {
	return isIdentifierStart(r) || unicode.IsDigit(r)
}

func isChinese(r rune) bool {
	return unicode.Is(unicode.Han, r)
}

func isASCIILetter(r rune) bool {
	return r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z'
}
