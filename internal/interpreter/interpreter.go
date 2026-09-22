// Package interpreter 执行当前阶段已经支持的 AST。
//
// 现在只执行整数赋值和整数四则运算。解释器先保持简单，作为后续编译
// 后端的行为参考。
package interpreter

import (
	"fmt"
	"sort"
	"strings"

	"github.com/mehaotian/xuyu/internal/ast"
	"github.com/mehaotian/xuyu/internal/source"
)

// Error 表示执行阶段遇到的错误。
type Error struct {
	Message string
	Span    source.Span
}

// Error 返回适合直接显示给使用者的错误文本。
func (e *Error) Error() string {
	return fmt.Sprintf("执行错误：%s，位置 %s", e.Message, e.Span.Start)
}

// Environment 保存程序执行过程中的整数变量。
type Environment struct {
	values map[string]int64
}

// NewEnvironment 创建一个空的变量环境。
func NewEnvironment() *Environment {
	return &Environment{values: make(map[string]int64)}
}

// Set 写入一个整数变量。
func (e *Environment) Set(name string, value int64) {
	e.values[name] = value
}

// Get 读取一个变量。第二个返回值表示变量是否存在。
func (e *Environment) Get(name string) (int64, bool) {
	value, ok := e.values[name]
	return value, ok
}

// String 返回稳定的调试表示，便于测试和查看执行结果。
func (e *Environment) String() string {
	names := make([]string, 0, len(e.values))
	for name := range e.values {
		names = append(names, name)
	}
	sort.Strings(names)

	items := make([]string, 0, len(names))
	for _, name := range names {
		items = append(items, fmt.Sprintf("%s=%d", name, e.values[name]))
	}
	return "环境[" + strings.Join(items, ", ") + "]"
}

// Run 执行一个已经完成语法分析的程序。
//
// 当前程序只能包含整数赋值。解释器不重新解析源码，也不猜测未知节点
// 的含义；遇到尚未支持的 AST 会直接返回错误。
func Run(program ast.Program) (*Environment, error) {
	environment := NewEnvironment()

	for _, statement := range program.Statements {
		assignment, ok := statement.(ast.Assignment)
		if !ok {
			return nil, &Error{
				Message: "当前只支持整数赋值",
				Span:    statement.Span(),
			}
		}

		value, err := evaluateExpression(assignment.Value, environment)
		if err != nil {
			return nil, err
		}

		environment.Set(assignment.Name, value)
	}

	return environment, nil
}

func evaluateExpression(expression ast.Expression, environment *Environment) (int64, error) {
	switch value := expression.(type) {
	case ast.IntegerLiteral:
		return value.Value, nil
	case ast.VariableReference:
		result, ok := environment.Get(value.Name)
		if !ok {
			return 0, &Error{
				Message: fmt.Sprintf("变量“%s”尚未定义", value.Name),
				Span:    value.Range,
			}
		}
		return result, nil
	case ast.BinaryExpression:
		left, err := evaluateExpression(value.Left, environment)
		if err != nil {
			return 0, err
		}
		right, err := evaluateExpression(value.Right, environment)
		if err != nil {
			return 0, err
		}

		switch value.Operator {
		case ast.OperatorAdd:
			return left + right, nil
		case ast.OperatorSubtract:
			return left - right, nil
		case ast.OperatorMultiply:
			return left * right, nil
		case ast.OperatorDivide:
			if right == 0 {
				return 0, &Error{
					Message: "除数不能为零",
					Span:    value.Right.Span(),
				}
			}
			return left / right, nil
		default:
			return 0, &Error{
				Message: "当前不支持这种运算符",
				Span:    value.Span(),
			}
		}
	default:
		return 0, &Error{
			Message: "当前不支持这种表达式",
			Span:    expression.Span(),
		}
	}
}
