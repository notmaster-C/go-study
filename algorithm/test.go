package algorithm

import "fmt"

func Test() {
	//fmt.Println(convert("PAYPALISHIRING", 2))
	// fmt.Println(robot("URR", [][]int{{2, 2}}, 3, 2))
	nums := []int{1, 1, 1, 2, 2, 3}
	rotate(nums, 2)
	fmt.Println(nums)
}
