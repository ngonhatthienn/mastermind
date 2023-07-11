package user

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/redis/go-redis/v9"

	password "intern2023/handler/Password"
	// pb "intern2023/pb"
	"intern2023/share"
)

type User struct {
	ID       int32  `json:"_id"`
	Username     string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

func CheckExistUser(client *redis.Client, IdUser int) bool {
	IdUserString := strconv.Itoa(IdUser)
	UserKey := share.UserPattern(IdUserString)
	valUser, _ := client.Get(context.Background(), UserKey).Result()
	if valUser == "" {
		return false
	}
	return true
}

func CreateUser(client *redis.Client, Username string, Password string, Email string, Role string) int32 { // in *pb.CreateUserRequest not very okay

	min := 10000000
	max := 99999999
	XId := share.CreateRandomNumber(min, max)
	hashedPassword := password.HashPassword(Password)
	item := User{ID: int32(XId), Username: Username, Email: Email, Password: hashedPassword, Role: Role}

	val, _ := json.Marshal(item)
	_, _ = client.Set(context.Background(), share.UserPattern(strconv.Itoa(XId)), val, 0).Result()

	return item.ID
}

// func GetListUser(client *redis.Client) (int, []*pb.User) {
// 	keys, _ := client.Keys(context.Background(), share.AllUserPattern()).Result()

// 	cmdS, _ := client.Pipelined(context.Background(), func(pipe redis.Pipeliner) error {
// 		for _, key := range keys {
// 			pipe.Get(context.Background(), key).Result()
// 		}
// 		return nil
// 	})

// 	var Users []*pb.User
// 	for _, cmd := range cmdS {
// 		val := cmd.(*redis.StringCmd).Val()
// 		var data *pb.User
// 		_ = json.Unmarshal([]byte(val), &data)
// 		Users = append(Users, data)
// 	}

// 	return len(Users), Users
// }

func LogIn(client *redis.Client, username string, Password string) (int, bool) {
	keys, _ := client.Keys(context.Background(), share.AllUserPattern()).Result()
	cmdS, _ := client.Pipelined(context.Background(), func(pipe redis.Pipeliner) error {
		for _, key := range keys {
			pipe.Get(context.Background(), key).Result()
		}
		return nil
	})
	for _, cmd := range cmdS {
		val := cmd.(*redis.StringCmd).Val()
		// var data *pb.User
		var data User
		_ = json.Unmarshal([]byte(val), &data)

		if data.Username == username {
			return int(data.ID), password.CheckPassword(data.Password, Password)
		}
	}
	return 0, false
}
