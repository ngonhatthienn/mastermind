package gameApp

import (
	"strconv"
)


func FindString(arr []string, target string) bool {
	for _, str := range arr {
		if str == target {
			return true
		}
	}
	return false
}



// convert from int to Array
func convertIntArr(num int, arr *[5]int) {
	index := 4
	for index >= 0 {
		arr[index] = num % 10
		num /= 10
		index--
	}
}

func ConvertArrString(arr *[5]int) string {
	res := ""
	for i := 0; i < 5; i++ {
		res += strconv.Itoa(arr[i])
	}
	return res
}

func convertStringArr(s string) []int {
	intArr := make([]int, len(s))
	for i, c := range s {
		num, _ := strconv.Atoi(string(c))
		intArr[i] = num
	}
	return intArr
}


// Check number return false if repeat
func checkRepeat(input *[5]int, check *[10]int) bool {
	for i := 0; i < 5; i++ {
		if check[input[i]] > 1 {
			return false
		} else {
			check[input[i]]++
		}
	}
	return true
}

// Compare and return output when user play
func OutputGame(randomString string, inputString string) (int, int) {
	rightNumber := 0
	rightPosition := 0
	check := [10]int{}
	randomArr := convertStringArr(randomString)
	inputArr := convertStringArr(inputString)

	for i := 0; i < 5; i++ {
		check[randomArr[i]] = i + 1
	}

	for i := 0; i < 5; i++ {
		if check[inputArr[i]] != 0 {
			rightNumber++
			if check[inputArr[i]] == i+1 {
				rightPosition++
			}
		}
	}

	return rightNumber, rightPosition
}

func outputGameUser(randoms *[5]int, input *[5]int, rightNumber *int, rightPosition *int) bool {
	*rightNumber = 0
	*rightPosition = 0
	check := [10]int{}
	for i := 0; i < 5; i++ {
		check[randoms[i]] = i + 1
	}

	for i := 0; i < 5; i++ {
		if check[input[i]] != 0 {
			*rightNumber++
			if check[input[i]] == i+1 {
				*rightPosition++
			}
		}
	}

	return *rightPosition == 5
}

// Compare and return output when tool play
func outputGameTool(randoms *[5]int, input *[5]int) bool {
	for index := 0; index < 5; index++ {
		if randoms[index] != input[index] {
			return false
		}
		if index == 4 {
			if randoms[index] != input[index] {
				break
			} else {
				return true
			}
		}
	}
	return false
}

func generateDown(randoms *[5]int, input *[5]int) bool {
	for i := 98765; i >= 12345; i-- {
		digits := make(map[int]bool)
		n := i
		for j := 0; j < 5; j++ {
			digit := n % 10
			if digits[digit] {
				// Not distinct digits => skip
				break
			}
			digits[digit] = true
			n /= 10
			if j == 4 {
				// Check Ouput
				convertIntArr(i, input)
				if !outputGameTool(randoms, input) {
					break
				} else {
					return true
				}
			}
		}
	}
	return false
}

func generateUp(randoms *[5]int, input *[5]int) bool {
	for i := 12345; i <= 98765; i++ {
		digits := make(map[int]bool)
		n := i
		for j := 0; j < 5; j++ {
			digit := n % 10
			if digits[digit] {
				// Not distinct digits => skip
				break
			}
			digits[digit] = true
			n /= 10
			if j == 4 {
				// Check
				convertIntArr(i, input)
				if !outputGameTool(randoms, input) {
					break
				} else {
					return true
				}
			}
		}
	}
	return false
}

func toolPlay(randoms *[5]int, input *[5]int) bool {
	input1 := [5]int{}
	input2 := [5]int{}

	for {
		ch1 := make(chan bool)
		ch2 := make(chan bool)

		done := false
		go func() {
			ch1 <- generateUp(randoms, &input1)
		}()
		go func() {
			ch2 <- generateDown(randoms, &input2)
		}()
		if <-ch2 {
			done = true
			*input = input2
			break
		} else if <-ch1 {
			done = true
			*input = input1
			break
		}
		if done {
			break
		}
	}
	return true
}


