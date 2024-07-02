package algorithm

import (
	"fmt"
	"math"
	"sort"
	"strings"
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

/*
字符串转换整数 (atoi)
函数 myAtoi(string s) 的算法如下：

空格：读入字符串并丢弃无用的前导空格（" "）
符号：检查下一个字符（假设还未到字符末尾）为 '-' 还是 '+'。如果两者都不存在，则假定结果为正。
转换：通过跳过前置零来读取该整数，直到遇到非数字字符或到达字符串的结尾。如果没有读取数字，则结果为0。
舍入：如果整数数超过 32 位有符号整数范围 [−231,  231 − 1] ，需要截断这个整数，使其保持在这个范围内。具体来说，小于 −231 的整数应该被舍入为 −231 ，大于 231 − 1 的整数应该被舍入为 231 − 1 。
返回整数作为最终结果
*/
func myAtoi(s string) (r int) {
	t := strings.TrimSpace(s)
	symbol := 1
	if len(t) == 0 {
		return 0
	}
	switch t[0] {
	case '-':
		symbol = -1
		t = t[1:]
		break
	case '+':
		t = t[1:]
		break
	}
	for i := 0; i < len(t); i++ {
		n := t[i]
		if n < '0' || n > '9' {
			r = r * symbol
			return
		}
		r = r*10 + int(n-'0')
		m := r * symbol
		if m < math.MinInt32 {
			r = math.MinInt32
			return
		}
		if m > math.MaxInt32 {
			r = math.MaxInt32
			return
		}
	}
	r = r * symbol
	return
}

/*
回文数
给你一个整数 x ，如果 x 是一个回文整数，返回 true ；否则，返回 false 。

回文数是指正序（从左向右）和倒序（从右向左）读都是一样的整数。例如，121 是回文，而 123 不是。
*/
func isPalindromeV0(x int) bool {
	if x < 0 {
		return false
	}
	if x < 9 {
		return true
	}
	result := x
	temp := 0
	for result > 0 {
		temp = temp*10 + result%10
		result = result / 10
	}
	if temp == x {
		return true
	}
	return false
}
func isPalindrome(x int) bool {
	// 特殊情况：
	// 如上所述，当 x < 0 时，x 不是回文数。
	// 同样地，如果数字的最后一位是 0，为了使该数字为回文，
	// 则其第一位数字也应该是 0
	// 只有 0 满足这一属性
	if x < 0 || (x%10 == 0 && x != 0) {
		return false
	}

	revertedNumber := 0
	for x > revertedNumber {
		revertedNumber = revertedNumber*10 + x%10
		x /= 10
	}

	// 当数字长度为奇数时，我们可以通过 revertedNumber/10 去除处于中位的数字。
	// 例如，当输入为 12321 时，在 while 循环的末尾我们可以得到 x = 12，revertedNumber = 123，
	// 由于处于中位的数字不影响回文（它总是与自己相等），所以我们可以简单地将其去除。
	return x == revertedNumber || x == revertedNumber/10
}

/*
LCP 61. 气温变化趋势
*/
func temperatureTrend(temperatureA []int, temperatureB []int) int {
	n := len(temperatureA)
	ans, arr := 0, 0
	for i := 1; i < n; i++ {
		ta := getTrend(temperatureA[i-1], temperatureA[i])
		tb := getTrend(temperatureB[i-1], temperatureB[i])
		if ta == tb {
			arr++
			ans = max(ans, arr)
		} else {
			arr = 0
		}
	}
	return ans
}

func getTrend(x, y int) int {
	if x == y {
		return 0
	}
	if x < y {
		return -1
	}
	return 1
}

// LCP 01. 猜数字
func game(guess []int, answer []int) int {
	flag := 0
	for i := 0; i < 3; i++ {
		if guess[i] == answer[i] {
			flag++
		}
	}
	return flag
}

// LCP 02. 分式化简
func fraction(cont []int) []int {
	res := make([]int, 2)
	n := len(cont)
	res[0], res[1] = cont[n-1], 1
	for i := n - 2; i >= 0; i-- {
		res[1], res[0] = res[0], cont[i]*res[0]+res[1]
	}
	return res
}

// LCP 03. 机器人大冒险
// command命令 u x移动R y移动。
// obstacles障碍物坐标
// x，y平面大小
func robot(command string, obstacles [][]int, x int, y int) bool {
	// 如果目标点不在路径上，返回失败
	if !isOnThePath(command, x, y) {
		return false
	}
	for _, o := range obstacles {
		// 判断有效的故障点是否在路径上（故障的步数大于等于目标的点，视为无效故障）
		if (x+y > o[0]+o[1]) && isOnThePath(command, o[0], o[1]) {
			return false
		}
	}
	return true
}

func isOnThePath(command string, x int, y int) bool {
	uNum := strings.Count(command, "U")*((x+y)/len(command)) + strings.Count(command[0:(x+y)%len(command)], "U")
	rNum := strings.Count(command, "R")*((x+y)/len(command)) + strings.Count(command[0:(x+y)%len(command)], "R")
	if uNum == y && rNum == x {
		return true
	}
	return false
}

// 使数组元素全部相等的最少操作次数
// 排序+前缀和+二分查找
func minOperations(nums, queries []int) []int64 {
	n := len(nums)
	sort.Ints(nums)
	sum := make([]int, n+1) // 前缀和
	for i, x := range nums {
		sum[i+1] = sum[i] + x
	}
	ans := make([]int64, len(queries))
	for i, q := range queries {
		j := sort.SearchInts(nums, q)
		left := q*j - sum[j]               // 蓝色面积
		right := sum[n] - sum[j] - q*(n-j) // 绿色面积
		ans[i] = int64(left + right)
	}
	return ans
}

// 2741. 特别的排列
const MOD int64 = 1000000007

func specialPerm(nums []int) int {
	n := len(nums)
	// 1<<n   2的n次方
	f := make([][]int64, 1<<n)
	for i := range f {
		f[i] = make([]int64, n)
	}
	for i := 0; i < n; i++ {
		f[1<<i][i] = 1
	}

	for state := 1; state < (1 << n); state++ {
		for i := 0; i < n; i++ {
			if state>>i&1 == 0 {
				continue
			}
			for j := 0; j < n; j++ {
				if i == j || state>>j&1 == 0 {
					continue
				}
				x := nums[i]
				y := nums[j]
				if x%y != 0 && y%x != 0 {
					continue
				}
				f[state][i] = (f[state][i] + f[state^(1<<i)][j]) % MOD
			}
		}
	}

	var sum int64
	for i := 0; i < n; i++ {
		sum = (sum + f[(1<<n)-1][i]) % MOD
	}
	return int(sum)
}

// [困难] 10. 正则表达式匹配
func isMatch(s string, p string) bool {
	m, n := len(s), len(p)
	matches := func(i, j int) bool {
		if i == 0 {
			return false
		}
		if p[j-1] == '.' {
			return true
		}
		return s[i-1] == p[j-1]
	}

	f := make([][]bool, m+1)
	for i := 0; i < len(f); i++ {
		f[i] = make([]bool, n+1)
	}
	f[0][0] = true
	for i := 0; i <= m; i++ {
		for j := 1; j <= n; j++ {
			if p[j-1] == '*' {
				f[i][j] = f[i][j] || f[i][j-2]
				if matches(i, j-1) {
					f[i][j] = f[i][j] || f[i-1][j]
				}
			} else if matches(i, j) {
				f[i][j] = f[i][j] || f[i-1][j-1]
			}
		}
	}
	return f[m][n]
}

// 递归解法
func isMatchV1(s string, p string) bool {
	return dfs(s, p, len(s)-1, len(p)-1)
}

func dfs(s, p string, i, j int) bool {
	if j < 0 {
		return i < 0
	}
	if p[j] == '*' {
		if i < 0 || (p[j-1] != '.' && p[j-1] != s[i]) {
			return dfs(s, p, i, j-2)
		}
		return dfs(s, p, i-1, j) || dfs(s, p, i, j-2)
	}
	if i < 0 {
		return false
	}
	if p[j] == '.' || s[i] == p[j] {
		return dfs(s, p, i-1, j-1)
	}
	return false
}

// 盛最多水的容器
func maxArea(height []int) int {
	l := 0
	r := len(height) - 1
	t := 0
	for l < r {
		t = max(min(height[l], height[r])*(r-l), t)
		if height[l] < height[r] {
			l++
		} else {
			r--
		}
	}
	return t
}

// 整数转罗马数字
var valueSymbols = []struct {
	value  int
	symbol string
}{
	{1000, "M"},
	{900, "CM"},
	{500, "D"},
	{400, "CD"},
	{100, "C"},
	{90, "XC"},
	{50, "L"},
	{40, "XL"},
	{10, "X"},
	{9, "IX"},
	{5, "V"},
	{4, "IV"},
	{1, "I"},
}

func intToRoman(num int) string {
	roman := []byte{}
	for _, vs := range valueSymbols {
		for num >= vs.value {
			num -= vs.value
			roman = append(roman, vs.symbol...)
		}
		if num == 0 {
			break
		}
	}
	return string(roman)
}

// 2742. 给墙壁刷油漆

// 给墙壁刷油漆
// 给你两个长度为 n 下标从 0 开始的整数数组 cost 和 time ，分别表示给 n 堵不同的墙刷油漆需要的开销和时间。你有两名油漆匠：
//   - 一位需要 付费 的油漆匠，刷第 i 堵墙需要花费 time[i] 单位的时间，开销为 cost[i] 单位的钱。
//   - 一位 免费 的油漆匠，刷 任意 一堵墙的时间为 1 单位，开销为 0 。但是必须在付费油漆匠 工作 时，免费油漆匠才会工作。
//
// 请你返回刷完 n 堵墙最少开销为多少。
// https://leetcode.cn/problems/painting-the-walls
func paintWalls(cost []int, time []int) int {

	return paintWalls2(cost, time)
}

func paintWalls1(cost []int, time []int) int {

	// dp在每个位置做出方案 在结尾的时候判断方案是否合理 + memo偏移

	// 免费的油漆匠依赖付费的油漆匠, 但不依赖具体的, 只依赖的时间, 只要付费的时间 > 免费的时间即可
	// dfs(i, payTime, freeTime) 三个参数可以优化为2个 deltaTime = payTime-freeTime,只要最后的时候 deltaTime>=0即可

	n := len(cost)
	inf := math.MaxInt / 2 // 有加法, inf适当减小

	memo := make([][]int, n)
	for i := range memo {
		memo[i] = make([]int, 2*n+1) // trick!! 只要 delta>=n,那么可以免费刷所有的墙,j绝不会超过n,但 payTime[i]可能超过,这种情况下,不让其访问到memo,直接返回
		for j := range memo[i] {
			memo[i][j] = -1
		}
	}

	offset := n
	var dfs func(i, delta int) int
	dfs = func(i, delta int) int {
		if i == n {
			// 注意!! 方案合不合法在最后的位置判断, delta>=0才合法,
			if delta >= offset {
				return 0
			} else {
				return inf
			}
		}

		if delta-offset >= n-i { // trick!! 剪枝,因为免费不增加cost,所以能免费时肯定免费
			return 0
		}

		if memo[i][delta] != -1 {
			return memo[i][delta]
		}

		// 选择付费
		res := dfs(i+1, delta+time[i]) + cost[i]

		// 选择免费
		res = min(res, dfs(i+1, delta-1))

		memo[i][delta] = res
		return res
	}
	ans := dfs(0, offset) // trick!! 因为需要使用memo保存中间结果,但delta可能是负数(前面全选择免费,中间一个选择付费),所以偏移n,让memo的deltaIdx不会为负数
	return ans
}

func paintWalls2(cost []int, time []int) int {

	// 0-1背包
	// 付费个数 + 免费个数 = n
	// 付费时间 >= 免费时间=免费个数
	// 付费时间 >= n-付费个数
	// (付费时间+1) >=n
	// time[i]+1 为物品体积 cost[i]为物品价值, 体积>=n的情况下,价值最小是多少

	n := len(cost)
	f := make([]int, n+1) // j: 0->i号物品,选于不选(付费还是免费)的情况下, 体积>=j时,能获取的最少收益
	for j := range f {
		f[j] = math.MaxInt / 2 // 一个物品不选时,体检不可能>j,所以设置为极大值(防止后面加法溢出 inf/2)
	}
	f[0] = 0

	s := 0
	for i, c := range cost {
		t := time[i] + 1
		s += t
		for j := min(s, n); j > 0; j-- {
			f[j] = min(f[j], f[max(0, j-t)]+c) // 选与不选
		}
	}
	return f[n]
}

// 2766. 重新放置石块
func relocateMarbles(nums []int, moveFrom []int, moveTo []int) []int {
	set := map[int]interface{}{}
	for _, i := range nums {
		set[i] = 0
	}
	for k, i := range moveFrom {
		delete(set, i)
		set[moveTo[k]] = 0
	}
	res := make([]int, 0, len(set))
	for k := range set {
		res = append(res, k)
	}
	sort.Ints(res)
	return res

}

// 88. 合并两个有序数组
func merge(nums1 []int, m int, nums2 []int, n int) {
	l := len(nums1) - 1
	for n > 0 {
		if m > 0 && nums1[m-1] > nums2[n-1] {
			nums1[l] = nums1[m-1]
			m--
			l--
		} else {
			nums1[l] = nums2[n-1]
			l--
			n--
		}
	}
}

// 27. 移除元素
func removeElement(nums []int, val int) int {
	n := len(nums)
	left := 0
	for right := 0; right < n; right++ {
		if nums[right] != val {
			nums[left] = nums[right]
			left++
		}
	}
	return left
}

// 26. 删除有序数组中的重复项
func removeDuplicatesV0(nums []int) (l int) {
	ans := nums[0] - 1
	nums2 := make([]int, len(nums))
	for _, v := range nums {
		if ans == v {
			continue
		}
		ans = v
		nums2[l] = v
		l++
	}
	nums = nums2
	fmt.Println(nums)
	return
}
func removeDuplicates(nums []int) int {
	n := len(nums)
	if n == 0 {
		return 0
	}
	slow := 1
	for fast := 1; fast < n; fast++ {
		if nums[fast] != nums[fast-1] {
			nums[slow] = nums[fast]
			slow++
		}
	}
	return slow
}

// 80. 删除有序数组中的重复项 II
