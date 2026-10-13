package main

func twoSum(nums []int, target int) []int {
	result := []int{0, 0}
	numIndex := make(map[int]int)
	for i, v := range nums {
		index, ok := numIndex[v]
		if ok {
			result[0] = index
			result[1] = i
			return result
		}
		numIndex[target-v] = i
	}
	return result
}
