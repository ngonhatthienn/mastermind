package ToProto

import (
	game "intern2023/handler/Game"
	leaderboard "intern2023/handler/Leaderboard"
	user "intern2023/handler/User"
	pb "intern2023/pb"
)

func ToListUserProto(users []user.User) []*pb.User {
	var userProtos []*pb.User
	for _, user := range users {
		userProto := &pb.User{
			XId:      user.ID,
			Username: user.Username,
			Email:    user.Email,
			Password: user.Password,
			Role:     user.Role,
		}
		userProtos = append(userProtos, userProto)
	}
	return userProtos
}

func ToListGameProto(games []game.Game, isAdmin bool) []*pb.Game {
	var gameProtos []*pb.Game
	for _, game := range games {
		gameProto := &pb.Game{}
		if isAdmin {
			gameProto = &pb.Game{
				XId:        int32(game.ID),
				Game:       game.Game,
				GuessLimit: int32(game.GuessLimit),
			}
		} else {
			gameProto = &pb.Game{
				XId:        int32(game.ID),
				GuessLimit: int32(game.GuessLimit),
			}
		}
		gameProtos = append(gameProtos, gameProto)
	}
	return gameProtos
}

func ToLeaderBoardProto(leaderboards []leaderboard.LeaderBoard) []*pb.LeaderBoard {
	var leaderboardProtos []*pb.LeaderBoard
	for _, leaderboard := range leaderboards {
		leaderboardProto := &pb.LeaderBoard{
			UserId: int32(leaderboard.UserId),
			Score:  leaderboard.Score,
		}
		leaderboardProtos = append(leaderboardProtos, leaderboardProto)
	}
	return leaderboardProtos
}
