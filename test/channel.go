package test

import "fmt"

/*
 */
func testChannel() {
	ch := make(chan int, 1)
	defer close(ch)
	go func() { ch <- 3 + 4 }()
	v := <-ch
	fmt.Println(v)
	// 往一个已经被close的channel中继续发送数据会导致run-time panic。
}
