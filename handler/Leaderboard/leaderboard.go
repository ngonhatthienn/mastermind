package leaderboard

import (
	"context"
	"strconv"

	"github.com/redis/go-redis/v9"

	"intern2023/share"
)

type LeaderBoard struct {
	UserId         int    `json:"userId"`
	Score       string `json:"score"`
}

func AddScore(client *redis.Client, userId string, IdGame string, score int64) error {
	_,err := client.ZAdd(context.Background(), share.LeaderBoardPattern(IdGame), redis.Z{
		Score:  float64(score),
		Member: userId,
	}).Result()
	return err
}

func GetUserRank(client *redis.Client, IdGame string, IdUser string) (int32, string) {
	ctx := context.Background()
	rank, err := client.ZRevRank(ctx, share.LeaderBoardPattern(IdGame), IdUser).Result()
    if err != nil {
        return 0, ""
    }

    score, err := client.ZScore(ctx, share.LeaderBoardPattern(IdGame), IdUser).Result()
    if err != nil {
		return 0, ""
    }
	return int32(rank + 1),  strconv.Itoa(int(score))
}

func GetLeaderboard(client *redis.Client, IdGame string, size int64, IdUser string) ([]LeaderBoard, error) {
	results, err := client.ZRevRangeWithScores(context.Background(), share.LeaderBoardPattern(IdGame), 0, size-1).Result()
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
