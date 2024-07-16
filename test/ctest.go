package test

import (
	"fmt"
	"sync"
)

// 面试题

func CTest() {
	a := [3]int{1, 2, 3}
	// a := []int{1, 2, 3}
	for k, v := range a {
		if k == 0 {
			a[0], a[1] = 100, 200
			fmt.Print(a)
		}
		a[k] = 100 + v
	}
	fmt.Print(a)
}
func CTest2() {
	a := [3]int{1, 2, 3}
	for k, v := range a {
		if k == 0 {
			a[0], a[1] = 100, 200
			fmt.Print(a)
		}
		a[k] = 100 + v
	}
	fmt.Print(a)
}

// 打印 cat dog fish 一百次 要求有顺序
func CTest3_my() {
	var wg sync.WaitGroup
	var mu sync.Mutex // 使用互斥锁来保证顺序

	wg.Add(30)
	for i := 0; i < 100; i++ {
		go func() {
			defer wg.Done()
			mu.Lock()         // 锁定互斥锁
			defer mu.Unlock() // 确保互斥锁被释放
			Cat()
			Dog()
			Fish()
		}()
	}
	wg.Wait()
}
func Dog() {
	fmt.Println("dog")
}
func Cat() {
	fmt.Println("cat")
}
func Fish() {
	fmt.Println("fish")
}
func CTest3() {
	var wg sync.WaitGroup
	wg.Add(100)

	// 创建三个带缓存的channel
	ch1 := make(chan struct{}, 1)
	ch2 := make(chan struct{}, 1)
	ch3 := make(chan struct{}, 1)

	// 初始化第一个信号
	ch1 <- struct{}{}

	for i := 0; i < 100; i++ {
		go func() {
			defer wg.Done()
			<-ch1             // 等待前一个循环的Cat打印完成
			Cat_o(ch1, ch2)   // Cat打印
			ch3 <- struct{}{} // 通知下一个循环的Fish可以开始
		}()
		go func() {
			defer wg.Done()
			<-ch2           // 等待上一个Dog打印完成
			Dog_o(ch2, ch3) // Dog打印
		}()
		go func() {
			defer wg.Done()
			<-ch3            // 等待上一个Fish打印完成
			Fish_o(ch3, ch1) // Fish打印，然后循环回到Cat
		}()
	}

	wg.Wait()
}

func Dog_o(rev, send chan struct{}) {
	<-rev
	fmt.Println("dog")
	send <- struct{}{}
}
func Cat_o(rev, send chan struct{}) {
	<-rev
	fmt.Println("cat")
	send <- struct{}{}
}
func Fish_o(rev, send chan struct{}) {
	<-rev
	fmt.Println("fish")
	send <- struct{}{}
}
