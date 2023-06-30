package share
import (
	"math/rand"
	"time"
	"strings"
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
