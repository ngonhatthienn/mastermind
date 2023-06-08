package main

import (
	"fmt"
	"math/rand"
	"time"
)

// convert from int to Array
func convertIntArr(num int, arr *[5]int) {
	index := 4
	for index >= 0 {
		arr[index] = num % 10
		num /= 10
		index--
	}
}

// Create randoms number
func createNumRand(randoms *[5]int) {
	check := [10]int{}

	for i := 0; i < 5; i++ {
		create := rand.Intn(9) + 1
		if check[create] == 0 {
			randoms[i] = create
			check[create]++
		} else {
			i--
		}
	}
	// for i := 0; i < 5; i++ {
	// 	fmt.Println(randoms[i])
	// }
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

// Get input user
func inputUser(input *[5]int) {
	check := [10]int{}
	for {
		// Input
		for i := 0; i < 10; i++ {
			check[i] = 0
		}
		for i := 0; i < 5; i++ {
			fmt.Scan(&input[i])
			check[input[i]]++
		}
		// Check
		if checkRepeat(input, &check) {
			break
		} else {
			fmt.Println("Your number is invalid! Please enter 5 numbers again: ")
		}
	}
}

// Compare and return output when user play
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

func main() {
	randoms := [5]int{}
	input := [5]int{}
	
	win := false
	rightNumber := 0
	rightPosition := 0
	var player int

	createNumRand(&randoms)
	fmt.Print("Nguoi choi: 1, May choi: 2: ")
	fmt.Scanln(&player)

	if player == 1 {
		for i := 1; i <= 100; i++ {
			fmt.Println("Doan 5 so lan", i, ": ")
			inputUser(&input)
			if outputGameUser(&randoms, &input, &rightNumber, &rightPosition) {
				win = true
				break
			}
			fmt.Println(rightNumber, " ", rightPosition)

		}
		if win {
			fmt.Println("You are win!!!")
		} else {
			fmt.Println("You are lose!!!")
		}
	} else if player == 2 {
		start := time.Now()
		if toolPlay(&randoms, &input) {
			fmt.Print("The result is:")
			for i := 0; i < 5; i++ {
				if i == 4 {
					fmt.Println(input[i])
				} else {
					fmt.Print(input[i], " ")
				}
			}
		}
		t := time.Now()
		elapsed := t.Sub(start).Microseconds()
		fmt.Println("It takes:", elapsed, "microseconds")
	}
}
