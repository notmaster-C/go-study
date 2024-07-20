package algorithm

import "fmt"

func Test() {
	//fmt.Println(convert("PAYPALISHIRING", 2))
	// fmt.Println(robot("URR", [][]int{{2, 2}}, 3, 2))
	solu_122()
}
func solu_122() {
	prices := []int{7, 1, 5, 3, 6, 4}
	s := maxProfit(prices)
	fmt.Println(s)
}
