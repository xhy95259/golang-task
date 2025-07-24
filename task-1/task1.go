package main

import (
	"fmt"
	"sort"
)

/*
------只出现一次的数字
给定一个非空整数数组，除了某个元素只出现一次以外，
其余每个元素均出现两次。找出那个只出现了一次的元素。
可以使用 for 循环遍历数组，结合 if 条件判断和 map 数据结构来解决，
例如通过 map 记录每个元素出现的次数，然后再遍历 map
找到出现次数为1的元素
*/
func singleNumber(nums []int) int {
	single := 0 // 只用了一个变量
	for _, num := range nums {
		single ^= num // 原地修改
	}
	return single
}

// 判断一个整数是否是回文数
func isPalindrome(x int) bool {
	// 负数或末尾是0的正数（除了0本身）不可能是回文数
	if x < 0 || (x%10 == 0 && x > 0) {
		return false
	}
	rev := 0
	// 反转数字的后半部分，直到反转数 >= 剩余的前半部分
	for rev < x {
		rev = rev*10 + x%10
		x /= 10
	}
	// 比较前半部分和反转后的后半部分
	return rev == x || rev == x/10
}

/*
给定一个只包括 '('，')'，'{'，'}'，'['，']' 的字符串，判断字符串是否有效
*/
func isValid(s string) bool {
	// 创建一个栈用于存储左括号
	var stack []rune
	// 定义括号对应关系
	pairs := map[rune]rune{
		')': '(',
		'}': '{',
		']': '[',
	}
	// 遍历字符串中的每个字符
	for _, char := range s {
		// 如果是右括号
		if matching, ok := pairs[char]; ok {
			// 如果栈为空或者栈顶元素与当前右括号不匹配，返回false
			if len(stack) == 0 || stack[len(stack)-1] != matching {
				return false
			}
			// 匹配成功，弹出栈顶元素
			stack = stack[:len(stack)-1]
		} else {
			// 如果是左括号，压入栈中
			stack = append(stack, char)
		}
	}
	// 如果栈为空，说明所有括号都匹配成功
	return len(stack) == 0
}

// 查找字符串数组中的最长公共前缀
func longestCommonPrefix(strs []string) string {
	if len(strs) == 0 {
		return ""
	}
	// 以第一个字符串为基准
	prefix := strs[0]
	// 逐个比较字符串
	for i := 1; i < len(strs); i++ {
		// 找出当前prefix与当前字符串的公共前缀
		j := 0
		for j < len(prefix) && j < len(strs[i]) && prefix[j] == strs[i][j] {
			j++
		}
		// 更新prefix
		prefix = prefix[:j]
		// 如果prefix已经为空，直接返回
		if prefix == "" {
			return ""
		}
	}
	return prefix
}

// 给定一个由整数组成的非空数组所表示的非负整数，在该数的基础上加一
func plusOne(digits []int) []int {
	n := len(digits)
	// 从最低位开始处理
	for i := n - 1; i >= 0; i-- {
		// 如果当前位小于9，直接加1并返回
		if digits[i] < 9 {
			digits[i]++
			return digits
		}
		// 当前位是9，加1后变成0，进位
		digits[i] = 0
	}
	// 如果所有位都是9，需要在最高位前添加1
	result := make([]int, n+1)
	result[0] = 1
	return result
}

/*
------删除有序数组中的重复项
删除有序数组中的重复项：给你一个有序数组 nums ，
请你原地删除重复出现的元素，使每个元素只出现一次，
返回删除后数组的新长度。不要使用额外的数组空间，
你必须在原地修改输入数组并在使用 O(1) 额外空间的条件下完成。
可以使用双指针法，一个慢指针 i 用于记录不重复元素的位置，
一个快指针 j 用于遍历数组，当 nums[i] 与 nums[j] 不相等时，将 nums[j] 赋值给 nums[i + 1]，
并将 i 后移一位。
*/
func removeDuplicates(nums []int) int {
	if len(nums) == 0 {
		return 0
	}
	// 慢指针i表示不重复元素的位置
	i := 0
	// 快指针j用于遍历数组
	for j := 1; j < len(nums); j++ {
		// 如果当前元素与慢指针指向的元素不同
		if nums[j] != nums[i] {
			// 将不重复元素移到i+1的位置
			i++
			nums[i] = nums[j]
		}
	}
	// 返回新数组的长度
	return i + 1
}

/*
------合并区间
以数组 intervals 表示若干个区间的集合，
其中单个区间为 intervals[i] = [starti, endi] 。
请你合并所有重叠的区间，并返回一个不重叠的区间数组，
该数组需恰好覆盖输入中的所有区间。可以先对区间数组按照区间的起始位置进行排序，
然后使用一个切片来存储合并后的区间，遍历排序后的区间数组，将当前区间与切片中最后一个区间进行比较
，如果有重叠，则合并区间；如果没有重叠，则将当前区间添加到切片中
*/
func merge(intervals [][]int) [][]int {
	if len(intervals) <= 1 {
		return intervals
	}
	// 按区间起始位置排序
	sort.Slice(intervals, func(i, j int) bool {
		return intervals[i][0] < intervals[j][0]
	})
	result := [][]int{intervals[0]}
	for i := 1; i < len(intervals); i++ {
		current := intervals[i]
		last := result[len(result)-1]
		// 如果当前区间的起始位置小于等于上一个区间的结束位置，合并区间
		if current[0] <= last[1] {
			// 更新结束位置为两个区间结束位置的较大值
			if current[1] > last[1] {
				last[1] = current[1]
			}
		} else {
			// 没有重叠，添加到结果中
			result = append(result, current)
		}
	}
	return result
}

/*
------两数之和
给定一个整数数组 nums 和一个目标值 target，请你在该数组中找出和为目标值的那两个整数
*/
func twoSum(nums []int, target int) []int {
	// 使用map存储数字和索引
	numMap := make(map[int]int)
	for i, num := range nums {
		// 计算当前数字需要的配对数字
		complement := target - num
		// 检查配对数字是否已经在map中
		if idx, ok := numMap[complement]; ok {
			return []int{idx, i}
		}
		// 将当前数字和索引存入map
		numMap[num] = i
	}
	// 如果没有找到，返回空切片（题目保证有唯一解，所以这里不会执行）
	return []int{}
}

func main() {

	//1 只出现一次的数字
	test1 := []int{4, 1, 2, 1, 2}
	fmt.Println("输入:", test1)
	fmt.Println("只出现一次的元素:", singleNumber(test1))

	//2 回文数
	test2 := []int{121, -121, 10, 12321, 12345, 0, 1}
	for _, tc := range test2 {
		fmt.Printf("输入: %d, 是否为回文数: %v\n", tc, isPalindrome(tc))
	}

	//3 有效的括号
	test3 := []string{
		"()",
		"()[]{}",
		"(]",
		"([)]",
		"{[]}",
		"",
	}
	for _, tc := range test3 {
		fmt.Printf("输入: \"%s\", 是否有效: %v\n", tc, isValid(tc))
	}

	//4 最长公共前缀
	test4 := []string{"flower", "flow", "flight"}
	fmt.Println("最长公共前缀:", longestCommonPrefix(test4))

	//5 加一
	test5 := []int{1, 2, 3}
	fmt.Println("整数加一:", plusOne(test5))

	//6 删除有序数组中的重复项
	test6 := []int{1, 1, 2}
	len1 := removeDuplicates(test6)
	fmt.Printf("删除重复项后的长度: %d, 数组前%d个元素: %v\n", len1, len1, test6[:len1])

	//7 合并区间
	test7 := [][]int{{1, 3}, {2, 6}, {8, 10}, {15, 18}}
	fmt.Println("合并区间:", merge(test7))

	//8 两数之和
	test8 := []int{2, 7, 11, 15}
	target1 := 9
	fmt.Printf("两数之和: 目标值=%d, 结果=%v\n", target1, twoSum(test8, target1))

}
