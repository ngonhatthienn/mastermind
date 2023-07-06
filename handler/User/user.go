package user

import (
	"context"
	"encoding/json"
	"strconv"
	"fmt"
	"github.com/redis/go-redis/v9"

	password "intern2023/handler/Password"
	pb "intern2023/pb"
	"intern2023/share"
)

type UserItem struct {
	ID       int32  `json:"_id"`
	Name     string `json:"name"`
	Password string `json:"password"`
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

func CreateUser(client *redis.Client, Name string, Password string) int32 { // in *pb.CreateUserRequest not very okay

	min := 10000000
	max := 99999999
	XId := share.CreateRandomNumber(min, max)
	hashedPassword := password.HashPassword(Password)
	item := UserItem{ID: int32(XId), Name: Name, Password: hashedPassword}

	val, _ := json.Marshal(item)
	_, _ = client.Set(context.Background(), share.UserPattern(strconv.Itoa(XId)), val, 0).Result()

	return item.ID
}

func GetListUser(client *redis.Client) (int, []*pb.User) {
	keys, _ := client.Keys(context.Background(), share.AllUserPattern()).Result()

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

func LogIn(client *redis.Client, username string, Password string) bool {
	keys, _ := client.Keys(context.Background(), share.AllUserPattern()).Result()
	fmt.Println("Login: ", keys)
	cmdS, _ := client.Pipelined(context.Background(), func(pipe redis.Pipeliner) error {
		for _, key := range keys {
			pipe.Get(context.Background(), key).Result()
		}
		return nil
	})
	for _, cmd := range cmdS {
		val := cmd.(*redis.StringCmd).Val()
		var data *pb.User
		_ = json.Unmarshal([]byte(val), &data)
		fmt.Println("Data: ", data.Name)
		fmt.Println("username: ",username)

		if data.Name == username {

			return password.CheckPassword(data.Password, Password)
		}
	}
	return false
}
