package user

import (
	"encoding/json"
	"strconv"
	"intern2023/database"
	pb "intern2023/pb"
	// game "intern2023/handler/Game"
	"github.com/redis/go-redis/v9"
	"intern2023/share"

)
type UserItem struct {
	ID         int32  `json:"_id"`
	Name       string `json:"name"`
	Password   string  `json:"password"`
}

func CreateUser(client *redis.Client, in *pb.CreateUserRequest) (int32, string){ //in *pb.CreateUserRequest not very okay

	min := 10000000
	max := 99999999
	XId := shareFunc.CreateRandomNumber(min, max)
	item := UserItem{ID: int32(XId), Name: in.Name, Password: in.Password}

	val, _ := json.Marshal(item)
	_, err1 := database.Set(client,"user:"+strconv.Itoa(XId), val, 0)

	if err1 != nil {
		panic(err1)
	}
	return item.ID, item.Name
}

func GetListUser(client *redis.Client) []*pb.User{
	var keys []string
	var userData []string
	keys, _ = database.Keys(client, "user:*")

	for _, key := range keys {
		val, _ := database.Get(client, key)
		userData = append(userData, val)
	}

	var Users []*pb.User
	for _, userData := range userData {
		var data *pb.User
		err := json.Unmarshal([]byte(userData), &data)
		if err != nil {
		}
		Users = append(Users, data)
	}
	return Users
}
