package console

import (
	"fmt"
	"os"

	"github.com/mgutz/ansi"
)

// Success 打印一条成功消息，绿色输出
func Success(msg string) {
	colorOut(msg, "green")
}

// Error打印一条错误消息，红色输出
func Error(msg string) {
	colorOut(msg, "red")
}

// Warning 打印一条提示消息，黄色输出
func Waring(msg string) {
	colorOut(msg, "yellow")
}

// Exit 打印一条报错消息，并退出os.Exit(1)
func Exit(msg string) {
	Error(msg)
	os.Exit(1)
}

// ExitIf 语法糖，自带err!=nil 判断
func ExitIf(err error) {
	if err != nil {
		Exit(err.Error())
	}
}

// colorOut 内部使用，设置高亮颜色
func colorOut(message, color string) {
	fmt.Fprintln(os.Stdout, ansi.Color(message, color))
}
