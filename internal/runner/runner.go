// Package runner 负责把一个源码文件交给现有的前端和解释器。
//
// 这里不定义新的语言规则，只负责连接文件读取、Lexer、Parser 和
// Interpreter，形成一条可以从命令行调用的最小路径。
package runner

import (
	"fmt"
	"os"

	"github.com/mehaotian/xuyu/internal/interpreter"
	"github.com/mehaotian/xuyu/internal/lexer"
	"github.com/mehaotian/xuyu/internal/parser"
)

// RunFile 读取并执行一个中文源码文件。
func RunFile(path string) (*interpreter.Environment, error) {
	source, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("读取文件失败：%w", err)
	}

	tokens, err := lexer.Lex(string(source))
	if err != nil {
		return nil, err
	}

	program, err := parser.Parse(tokens)
	if err != nil {
		return nil, err
	}

	return interpreter.Run(program)
}
