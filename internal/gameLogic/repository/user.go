package repository

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"

	"intern2023/database"
	password "intern2023/handler/Password"
	"intern2023/share"
)

type User struct {
	ID       int32  `json:"_id"`
	Fullname string `json:"fullname"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type UserRepositoryImpl struct {
	redisClient *redis.Client
	mongoClient *mongo.Client
}

func NewUserRepositoryImpl() *UserRepositoryImpl {
	redisClient, _ := database.ConnectRedisDatabase()
	mongoClient := database.ConnectMongoDBConnection()

	return &UserRepositoryImpl{
		redisClient: redisClient,
		mongoClient: mongoClient,
	}
}

func (r *UserRepositoryImpl) GetListUser() (int, []User) {
	keys, _ := r.redisClient.Keys(context.Background(), share.AllUserPattern()).Result()

	cmdS, _ := r.redisClient.Pipelined(context.Background(), func(pipe redis.Pipeliner) error {
		for _, key := range keys {
			pipe.Get(context.Background(), key).Result()
		}
		return nil
	})

	var Users []User
	for _, cmd := range cmdS {
		val := cmd.(*redis.StringCmd).Val()
		var user User
		_ = json.Unmarshal([]byte(val), &user)
		Users = append(Users, user)
	}
	return len(Users), Users

}

func (r *UserRepositoryImpl) CheckExistUser(IdUser int) bool {
	IdUserString := strconv.Itoa(IdUser)
	UserKey := share.UserPatternValue(IdUserString)
	valUser, _ :=  r.redisClient.Get(context.Background(), UserKey).Result()
	if valUser == "" {
		return false
	}
	return true
}

func (r *UserRepositoryImpl) LogIn(username string, Password string) (int, string, bool) {
	keys, _ := r.redisClient.Keys(context.Background(), share.AllUserPattern()).Result()
	cmdS, _ :=  r.redisClient.Pipelined(context.Background(), func(pipe redis.Pipeliner) error {
		for _, key := range keys {
			pipe.Get(context.Background(), key).Result()
		}
		return nil
	})
	for _, cmd := range cmdS {
		val := cmd.(*redis.StringCmd).Val()
		var data User
		_ = json.Unmarshal([]byte(val), &data)

		if data.Username == username {
			return int(data.ID), data.Role, password.CheckPassword(data.Password, Password)
		}
	}
	return 0, "", false
}

func (r *UserRepositoryImpl) CreateUser(Fullname string, Username string, Password string, Email string, Role string) int32 {
	min := 10000000
	max := 99999999
	XId := share.CreateRandomNumber(min, max)
	hashedPassword := password.HashPassword(Password)
	item := User{ID: int32(XId), Fullname: Fullname, Username: Username, Email: Email, Password: hashedPassword, Role: Role}

	val, _ := json.Marshal(item)
	_, _ =  r.redisClient.Set(context.Background(), share.UserPatternValue(strconv.Itoa(XId)), val, 0).Result()

	return item.ID
}
