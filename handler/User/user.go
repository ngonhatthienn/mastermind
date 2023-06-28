package user

import (
	"context"
	"encoding/json"
	"strconv"

	"intern2023/database"
	pb "intern2023/pb"
	shareFunc "intern2023/share"

	"github.com/redis/go-redis/v9"
)

type UserItem struct {
	ID       int32  `json:"_id"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

func AllUserPatterns(IdUser string) string {
	return "user:*"
}

func UserPattern(IdUser string) string {
	return "user:" + IdUser
}

func CheckExistUser(client *redis.Client, IdUser int) bool {
	IdUserString := strconv.Itoa(IdUser)
	UserKey := UserPattern(IdUserString)
	valUser, _ := client.Get(context.Background(), UserKey).Result()
	if valUser == "" {
		return false
	}

	return true
}

func CreateUser(client *redis.Client, in *pb.CreateUserRequest) (int32, string) { // in *pb.CreateUserRequest not very okay

	min := 10000000
	max := 99999999
	XId := shareFunc.CreateRandomNumber(min, max)
	item := UserItem{ID: int32(XId), Name: in.Name, Password: in.Password}

	val, _ := json.Marshal(item)
	_, _ = database.Set(client, "user:"+strconv.Itoa(XId), val, 0)

	return item.ID, item.Name
}

func GetListUser(client *redis.Client) (int, []*pb.User)  {
	// var keys []string
	// var userData []string
	// keys, _ = database.Keys(client, "user:*")

	// for _, key := range keys {
	// 	val, _ := database.Get(client, key)
	// 	userData = append(userData, val)
	// }

	// var Users []*pb.User
	// for _, userData := range userData {
	// 	var data *pb.User
	// 	err := json.Unmarshal([]byte(userData), &data)
	// 	if err != nil {
	// 	}
	// 	Users = append(Users, data)
	// }
	// return Users
	keys, _ := client.Keys(context.Background(), "user:*").Result() 

	cmdS, _ := client.Pipelined(context.Background(), func(pipe redis.Pipeliner) error {
		for _, key := range keys {
			pipe.Get(context.Background(), key).Result()
		}
		return nil
	})

	var Users []*pb.User
    for _, cmd := range cmdS {
        val := cmd.(*redis.StringCmd).Val()
		var data *pb.User
		_ = json.Unmarshal([]byte(val), &data)
		Users = append(Users, data)
    }

	return len(Users), Users
}
