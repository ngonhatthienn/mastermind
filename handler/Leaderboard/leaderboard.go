package leaderboard

import (
	"context"
	"strconv"
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

	_,err := client.ZAdd(context.Background(), leaderBoardKey, redis.Z{
		Score:  float64(score),
		Member: userId,
	}).Result()
	return err
}

func GetUserRank(client *redis.Client, IdGame string, IdUser string) (int32, string) {
	ctx := context.Background()
	leaderboardKey := "leaderboard:" + IdGame
	rank, err := client.ZRevRank(ctx, leaderboardKey, IdUser).Result()
    if err != nil {
        // handle error
    }

    score, err := client.ZScore(ctx, leaderboardKey, IdUser).Result()
    if err != nil {
        // handle error
    }
	return int32(rank + 1),  strconv.Itoa(int(score))
}

func GetLeaderboard(client *redis.Client, IdGame string, size int64, IdUser string) ([]*pb.LeaderBoardRank, error) {
	leaderboardKey := "leaderboard:" + IdGame
	results, err := client.ZRevRangeWithScores(context.Background(), leaderboardKey, 0, size-1).Result()
	if err != nil {
		return nil, err
	}
	// var LeaderBoardData *pb.LeaderBoardReply
	var scores []*pb.LeaderBoardRank
	for _, result := range results {
		userId, _ := strconv.Atoi(result.Member.(string))
		score := &pb.LeaderBoardRank{
			UserId: int32(userId),
			Score:  strconv.Itoa(int(result.Score)),
		}
		scores = append(scores, score)
	}
	// LeaderBoardData.Ranks = scores

	// LeaderBoardData.UserRank, LeaderBoardData.UserScore = GetUserRank(client, leaderboardKey, IdUser)
	return scores, nil
}
