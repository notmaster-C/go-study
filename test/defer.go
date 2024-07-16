package test

import (
	"errors"
	"fmt"
)

/*
这个Go语言的测试函数TestDefer演示了如何使用defer语句与recover来捕获和处理panic。这里是函数的执行流程：

函数TestDefer开始执行，并声明了一个error类型的变量err。

紧接着，一个匿名函数被声明并立即执行。这个匿名函数使用defer关键字，因此它会被推迟执行，直到TestDefer函数即将返回之前。

在匿名函数内部，使用recover来捕获在TestDefer函数中可能发生的panic。如果recover捕获到panic，它会返回panic的值（这里是字符串"error"），否则返回nil。

如果recover返回了非nil值，使用fmt.Sprintf将panic的值格式化为字符串，并使用errors.New创建一个新的错误。

raisePanic函数被调用，并执行panic("error")，这会导致当前的执行流程被中断，并开始搜索最近的defer来处理这个panic。

defer的匿名函数捕获到panic，并将其转换为一个错误，赋值给err变量。

raisePanic函数中的panic被处理后，TestDefer函数继续执行并返回err。

由于err被赋予了由panic信息创建的新错误，函数返回的错误将包含字符串"error"。

由于TestDefer是一个测试函数，它可能会被一个测试框架调用，并且通常测试框架会打印出返回的错误。如果没有其他的日志输出或格式化，这个函数最后会输出一个错误，其内容是字符串"error"。

请注意，TestDefer函数的返回值是一个error接口，通常测试框架会检查这个错误是否为nil来确定测试是否通过。如果错误不是nil，测试框架可能会记录这个错误信息。如果需要具体的输出格式，可能需要结合测试框架的用法。
*/
func TestDefer() error {
	var err error
	defer func() {
		if r := recover(); r != nil {
			err = errors.New(fmt.Sprintf("%s", r))
		}
	}()
	raisePanic()
	return err
}
func raisePanic() {
	panic("error")
}
