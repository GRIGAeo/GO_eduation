package main

func main(){
	nums := []int{2, 7, 11, 15}
	target := 9
	twoSum(nums, target)
}

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
