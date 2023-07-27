package repository

import (
	"context"
	"strconv"

	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"

	"intern2023/database"
	"intern2023/share"
)

type LeaderBoard struct {
	UserId         int    `json:"userId"`
	Score       string `json:"score"`
}


type LeaderBoardRepositoryImpl struct {
	redisClient *redis.Client
	mongoClient *mongo.Client
}

func NewLeaderBoardRepositoryImpl() *LeaderBoardRepositoryImpl {
	redisClient, _ := database.ConnectRedisDatabase()
	mongoClient := database.ConnectMongoDBConnection()

	return &LeaderBoardRepositoryImpl{
		redisClient: redisClient,
		mongoClient: mongoClient,
	}
}

func (r*LeaderBoardRepositoryImpl)AddScore( userId string, IdGame string, score int64) error {
	_,err := r.redisClient.ZAdd(context.Background(), share.LeaderBoardPattern(IdGame), redis.Z{
		Score:  float64(score),
		Member: userId,
	}).Result()
	return err
}

func (r*LeaderBoardRepositoryImpl)GetUserRank( IdGame string, IdUser string) (int32, string) {
	ctx := context.Background()
	rank, err := r.redisClient.ZRevRank(ctx, share.LeaderBoardPattern(IdGame), IdUser).Result()
    if err != nil {
        return 0, ""
    }

    score, err := r.redisClient.ZScore(ctx, share.LeaderBoardPattern(IdGame), IdUser).Result()
    if err != nil {
		return 0, ""
    }
	return int32(rank + 1),  strconv.Itoa(int(score))
}

func (r*LeaderBoardRepositoryImpl) GetLeaderboard(IdGame string, size int64, IdUser string) ([]LeaderBoard, error) {
	results, err := r.redisClient.ZRevRangeWithScores(context.Background(), share.LeaderBoardPattern(IdGame), 0, size-1).Result()
	if err != nil || len(results) == 0 {
		return nil, err
	}
	var leaderBoards []LeaderBoard
	for _, result := range results {
		userId, _ := strconv.Atoi(result.Member.(string))
		leaderBoard := LeaderBoard{
			UserId: userId,
			Score:  strconv.Itoa(int(result.Score)),
		}
		leaderBoards = append(leaderBoards, leaderBoard)
	}
	return leaderBoards, nil
}
