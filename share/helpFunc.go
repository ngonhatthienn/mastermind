package share

import (
	"context"
	"math/rand"
	"strings"
	"time"
	"github.com/redis/go-redis/v9"
)

func CreateRandomNumber(min int, max int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min+1) + min
}

// Create array number
func CreateArrRand(randoms *[5]int) {
	check := [10]int{}

	for i := 0; i < 5; i++ {
		create := CreateRandomNumber(1, 9)
		if check[create] == 0 {
			randoms[i] = create
			check[create]++
		} else {
			i--
		}
	}
}

// Get One Element in key
func GetKeyElement(key string, index int) string {
	parts := strings.Split(key, ":")
	return parts[index]
}

func GetTokenElement(Token string, index int) string {
	parts := strings.Split(Token, " ")
	return parts[index]
}

func CalcScore(redisClient *redis.Client,keySession string, guessLeft int) int{
	timeStart, _ := redisClient.HGet(context.Background(), keySession, "timeStart").Int64()
	savedTime := time.Unix(timeStart, 0)

	diffInSeconds := 5000 - time.Now().Sub(savedTime).Seconds()

	// Get right and pos
	right, _ := redisClient.HGet(context.Background(), keySession, "sumRight").Int()
	pos, _ := redisClient.HGet(context.Background(), keySession, "sumPos").Int()

	score := int(diffInSeconds) + guessLeft*100 + (right+pos)*2
	return score
}
