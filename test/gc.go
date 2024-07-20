package test

import (
	"fmt"
	"net/http"
	"sync"
)

type Student struct {
	Name string
	Age  int
}

func StudentRegister(name string, age int) *Student {
	s := new(Student) // 局部变量s逃逸到堆中
	s.Name = name
	s.Age = age
	return s
}
func testGC() {
	// 指针逃逸
	// 函数StudentRegister()内部的s为局部变量，其值通过函数返回值返回，s本身为一个指针，其指向的内存地址不会是栈而是堆。
	StudentRegister("Jim", 18)

	// 栈空间不足逃逸
	//Slice函数分配了一个长度为1000的切片，是否逃逸取决于栈的空间是否足够大。当切片长度不断增加到10000时就会发生逃逸
	// 实际上当栈的空间不足以存放当前对象或无法判断当前切片长度时会将对象分配到堆中
	Slice()

	// 动态类型逃逸
	// 很多函数的参数为interface类型，比如fmt.Println(a …interface{})，编译期间很难确定其参数具体类型，也会产生逃逸。
	s := "Escape"
	fmt.Println(s)

	// 闭包引用对象逃逸
	// 该函数返回一个闭包，闭包使用了函数的局部变量a,b，使用时通过该函数获取闭包，然后每次执行闭包都会一次输出Fibonacci数列。
	// Fibonacci()函数中原本属于局部变量的a和b由于闭包的引用，不得不将两者放到堆上，以至于产生了逃逸。
	f := Fibonacci()
	for i := 0; i < 10; i++ {
		fmt.Printf("Fibonacci: %d\n", f())
	}

	//总结
	// 栈上分配内存比在堆上分配内存有更高的效率
	// 栈上分配的内存不需要GC处理
	// 堆上分配的内存使用完毕会交给GC处理
	// 逃逸分析的目的是决定分配地址是栈还是堆
	// 逃逸分析在编译阶段完成
}
func Fibonacci() func() int {
	a, b := 0, 1
	return func() int {
		a, b = b, a+b
		return a
	}
}
func Slice() {
	s := make([]int, 1000, 1000)
	for index, _ := range s {
		s[index] = index
	}
}

// 开启pprof，监听请求
func testpprof() {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		ip := "0.0.0.0:6060"
		if err := http.ListenAndServe(ip, nil); err != nil {
			fmt.Printf("start pprof failed on %s\n", ip)
		}
	}()
	fmt.Println("continue~")
	wg.Wait()
}
