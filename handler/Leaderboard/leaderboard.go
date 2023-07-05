package leaderboard

import (
	"context"
	"strconv"

	"github.com/redis/go-redis/v9"

	pb "intern2023/pb"
	"intern2023/share"
)



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
        panic(err)
    }

    score, err := client.ZScore(ctx, share.LeaderBoardPattern(IdGame), IdUser).Result()
    if err != nil {
        panic(err)
    }
	return int32(rank + 1),  strconv.Itoa(int(score))
}

func GetLeaderboard(client *redis.Client, IdGame string, size int64, IdUser string) ([]*pb.LeaderBoardRank, error) {
	results, err := client.ZRevRangeWithScores(context.Background(), share.LeaderBoardPattern(IdGame), 0, size-1).Result()
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
	return scores, nil
}
