package leaderboard

import (
	"fmt"
	"strconv"

	"intern2023/database"
	pb "intern2023/pb"

	"github.com/redis/go-redis/v9"
)

func AllLeaderBoardPatterns(IdUser string) string {
	return "leaderboard:*"
}

func LeaderBoardPattern(IdGame string) string {
	return "leaderboard:" + IdGame
}

func AddScore(client *redis.Client, userId string, IdGame string, score int64) error {
	leaderBoardKey := "leaderboard:" + IdGame

	err := database.ZAdd(client, leaderBoardKey, redis.Z{
		Score:  float64(score),
		Member: userId,
	})
	return err
}

func GetLeaderboard(client *redis.Client, IdGame string, size int64) ([]*pb.LeaderBoardData, error) {
	leaderboardKey := "leaderboard:" + IdGame
	results, err := database.ZRevRangeWithScores(client, leaderboardKey, 0, size-1)
	if err != nil {
		return nil, err
	}

	var scores []*pb.LeaderBoardData
	for _, result := range results {
		userId, _ := strconv.Atoi(result.Member.(string))
		score := &pb.LeaderBoardData{
			UserId: int32(userId),
			Score:  strconv.Itoa(int(result.Score)),
		}
		scores = append(scores, score)
	}
	fmt.Println(scores)

	return scores, nil
}
