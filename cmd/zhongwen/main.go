// 中文语言实验器的命令行入口。
//
// 命令行只负责参数校验、文件运行和用户可见输出，具体的语言处理
// 由 internal/runner 及其下游模块完成。
package main

import (
	"fmt"
	"os"

	"github.com/mehaotian/xuyu/internal/runner"
)

const version = "0.0.1-dev"

func main() {
	fmt.Print(run(os.Args[1:]))
}

// run 将命令行参数转换成一段输出，方便用测试验证入口行为。
func run(args []string) string {
	if len(args) == 0 {
		return helpText()
	}

	switch args[0] {
	case "帮助", "--帮助", "-h", "--help":
		return helpText()
	case "版本", "--版本", "-v", "--version":
		return version + "\n"
	case "运行", "--运行", "run":
		return runFile(args[1:])
	default:
		return fmt.Sprintf("暂不支持命令：%s\n\n%s", args[0], helpText())
	}
}

func runFile(args []string) string {
	if len(args) == 0 {
		return "运行命令需要一个源码文件路径\n"
	}
	if len(args) > 1 {
		return "运行命令暂时只接受一个源码文件路径\n"
	}

	environment, err := runner.RunFile(args[0])
	if err != nil {
		return "运行失败：" + err.Error() + "\n"
	}
	return environment.String() + "\n"
}

func helpText() string {
	return "中文语言实验器 " + version + "\n" +
		"用法：中文 [命令]\n\n" +
		"当前命令：\n" +
		"  帮助    显示帮助信息\n" +
		"  版本    显示版本信息\n" +
		"  运行    运行一个 .中 源码文件\n"
}
