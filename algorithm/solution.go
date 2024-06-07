package algorithm

import (
	"fmt"
	"math"
	"sort"
)

/*
@question: 1.两数之和

给定一个整数数组 nums 和一个整数目标值 target，请你在该数组中找出 和为目标值 target  的那 两个 整数，并返回它们的数组下标。

你可以假设每种输入只会对应一个答案。但是，数组中同一个元素在答案里不能重复出现。

你可以按任意顺序返回答案。
*/
func twoSum(nums []int, target int) []int {
	hashtable := map[int]int{}
	for i, num := range nums {
		targetNum := target - num
		if p, ok := hashtable[targetNum]; ok {
			return []int{p, i}
		}
		hashtable[num] = i
	}
	return nil
}

/*
2.两数相加

给你两个 非空 的链表，表示两个非负的整数。它们每位数字都是按照 逆序 的方式存储的，并且每个节点只能存储 一位 数字。

请你将两个数相加，并以相同形式返回一个表示和的链表。

你可以假设除了数字 0 之外，这两个数都不会以 0 开头。
*/
func addTwoNumbers(l1 *ListNode, l2 *ListNode) (head *ListNode) {
	var l3 *ListNode
	temp := 0
	for l1 != nil || l2 != nil {
		if l1 != nil {
			temp = temp + l1.Val
			l1 = l1.Next
		}
		if l2 != nil {
			temp = temp + l2.Val
			l2 = l2.Next
		}
		sum := temp
		temp = sum / 10
		sum = sum % 10
		if head == nil {
			head = &ListNode{Val: sum}
			l3 = head
		} else {
			l3.Next = &ListNode{Val: sum}
			l3 = l3.Next
		}
	}
	if temp > 0 {
		l3.Next = &ListNode{Val: temp}
	}
	return
}

type ListNode struct {
	Val  int
	Next *ListNode
}

/*
@question: 3. 无重复字符的最长子串

给定一个字符串 s ，请你找出其中不含有重复字符的 最长子串的长度。
*/
func lengthOfLongestSubstring(s string) int {

	rk, ans := -1, 0
	n := len(s)
	m := map[byte]int{}

	for i := 0; i < n; i++ {
		if i != 0 {
			delete(m, s[i-1])
		}
		for rk+1 < n && m[s[rk+1]] == 0 {
			val, ok := m[s[rk+1]]
			if !ok {
				// 如果该键不存在于映射中，则需要初始化
				m[s[rk+1]] = 1
			} else {
				// 如果键存在，则递增其值
				m[s[rk+1]] = val + 1
			}

			rk++
		}
		ans = max(ans, rk+1-i)
	}
	return ans
}

/*
@question:4.寻找两个正序数组的中位数

给定两个大小分别为 m 和 n 的正序（从小到大）数组 nums1 和 nums2。请你找出并返回这两个正序数组的 中位数 。

算法的时间复杂度应该为 O(log (m+n)) 。

	//var nums1, nums2 = []int{1, 2}, []int{3, 4}
	//var nums1, nums2 = []int{1, 3}, []int{2}
	//findMedianSortedArrays(nums1, nums2)
*/
func findMedianSortedArraysV1(nums1 []int, nums2 []int) float64 {
	nums3 := append(nums1, nums2...)
	sort.Ints(nums3)
	fmt.Println(nums3)
	l := len(nums3)
	i := l / 2
	//todo: l <>0
	if l%2 == 0 {
		i2 := (nums3[i-1] + nums3[i])
		return float64(i2) / 2
	}
	i2 := nums3[i]
	fmt.Println(i2)
	return float64(i2)
}
func findMedianSortedArrays(nums1 []int, nums2 []int) float64 {
	tot := len(nums1) + len(nums2)
	if tot%2 == 0 {
		left := find(nums1, 0, nums2, 0, tot/2)
		right := find(nums1, 0, nums2, 0, tot/2+1)
		return float64(left+right) / 2.0
	} else {
		return float64(find(nums1, 0, nums2, 0, tot/2+1))
	}
}

func find(nums1 []int, i int, nums2 []int, j int, k int) int {
	if len(nums1)-i > len(nums2)-j {
		return find(nums2, j, nums1, i, k)
	}
	if k == 1 {
		//递归找到k=1时，表示找到了该k位数
		if i == len(nums1) {
			return nums2[j]
		} else {
			return min(nums1[i], nums2[j])
		}
	}
	if len(nums1) == i {
		return nums2[j+k-1]
	}
	si := min(len(nums1), i+k/2)
	sj := j + k - k/2
	//通过比较nums1[si-1]和nums[sj-1]，可以删除最小的部分
	if nums1[si-1] > nums2[sj-1] {
		//nums2左侧数据全部抛弃(都比两个数小)，剩下数组查找第k-(sj-j)位数
		return find(nums1, i, nums2, sj, k-(sj-j))
	} else {
		//nums1左侧数据全部抛弃(都比两个数小)，剩下数组查找第k-(si-i)位数
		return find(nums1, si, nums2, j, k-(si-i))
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

/*
@question:5. 最长回文子串
给你一个字符串 s，找到 s 中最长的回文子串。

示例 1：

输入：s = "babad"
输出："bab"
解释："aba" 同样是符合题意的答案。
示例 2：

输入：s = "cbbd"
输出："bb"

*/

func longestPalindromeV1(s string) string {
	dp := make([][]bool, len(s))
	result := s[0:1] //初始化结果(最小的回文就是单个字符)
	for i := 0; i < len(s); i++ {
		dp[i] = make([]bool, len(s))
		dp[i][i] = true // 根据case 1 初始数据
	}
	for length := 2; length <= len(s); length++ { //长度固定，不断移动起点
		for start := 0; start < len(s)-length+1; start++ { //长度固定，不断移动起点
			end := start + length - 1
			if s[start] != s[end] { //首尾不同则不可能为回文
				continue
			} else if length < 3 {
				dp[start][end] = true //即case 2的判断
			} else {
				dp[start][end] = dp[start+1][end-1] //状态转移
			}
			if dp[start][end] && (end-start+1) > len(result) { //记录最大值
				result = s[start : end+1]
			}
		}
	}
	return result
}
func longestPalindrome(s string) string {
	if s == "" {
		return ""
	}
	start, end := 0, 0
	for i := 0; i < len(s); i++ {
		left1, right1 := expandAroundCenter(s, i, i)
		left2, right2 := expandAroundCenter(s, i, i+1)
		if right1-left1 > end-start {
			start, end = left1, right1
		}
		if right2-left2 > end-start {
			start, end = left2, right2
		}
	}
	return s[start : end+1]
}

func expandAroundCenter(s string, left, right int) (int, int) {
	for ; left >= 0 && right < len(s) && s[left] == s[right]; left, right = left-1, right+1 {
	}
	return left + 1, right - 1
}

/*
将一个给定字符串 s 根据给定的行数 numRows ，以从上往下、从左到右进行 Z 字形排列。

比如输入字符串为 "PAYPALISHIRING" 行数为 3 时，排列如下：

P   A   H   N
A P L S I I G
Y   I   R
之后，你的输出需要从左往右逐行读取，产生出一个新的字符串，比如："PAHNAPLSIIGYIR"。

请你实现这个将字符串进行指定行数变换的函数：

string convert(string s, int numRows);
*/
func convertV0(s string, numRows int) string {
	if numRows < 2 {
		return s
	}
	return ""
}
func convert(s string, numRows int) string {
	n, r := len(s), numRows
	if r == 1 || r >= n {
		return s
	}
	t := r*2 - 2
	ans := make([]byte, 0, n)
	for i := 0; i < r; i++ { // 枚举矩阵的行
		for j := 0; j+i < n; j += t { // 枚举每个周期的起始下标
			ans = append(ans, s[j+i]) // 当前周期的第一个字符
			if 0 < i && i < r-1 && j+t-i < n {
				ans = append(ans, s[j+t-i]) // 当前周期的第二个字符
			}
		}
	}
	return string(ans)
}

/*
整数反转

给你一个 32 位的有符号整数 x ，返回将 x 中的数字部分反转后的结果。

如果反转后整数超过 32 位的有符号整数的范围 [−231,  231 − 1] ，就返回 0。

假设环境不允许存储 64 位整数（有符号或无符号）。

示例 1：

输入：x = 123
输出：321
示例 2：

输入：x = -123
输出：-321
示例 3：

输入：x = 120
输出：21
示例 4：

输入：x = 0
输出：0
*/
func reverse(x int) int {
	var ans int
	for x != 0 {
		ans = ans*10 + x%10
		x /= 10
	}
	if ans < math.MinInt32 || ans > math.MaxInt32 {
		return 0
	}
	return ans
}
func Test() {
	//fmt.Println(convert("PAYPALISHIRING", 2))
	fmt.Println(reverse(123456))
}
