package main

import (
	"fmt"
)

func main() {
	res := missingNumber([]int{1, 2, 3, 5, 5})
	fmt.Println(res)
	return
}

func missingNumber(nums []int) int {
	for i := 0; i < len(nums); i++ {
		if nums[i] > 0 {
			if nums[(nums[i]-1)] > 0 {
				nums[nums[i]-1] *= -1
			}
		}
	}
	for i := 0; i < len(nums); i++ {
		if nums[i] > 0 {
			return i + 1
		}
	}
	return 0
}
