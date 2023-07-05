package share

// GAME
func AllGamePattern() string {
	return "game:*"
}

func GamePattern(IdGame string) string {
	return "game:" + IdGame
}
// USER
func AllUserPattern() string {
	return "user:*"
}

func UserPattern(IdUser string) string {
	return "user:" + IdUser
}

// LEADERBOARD
func AllLeaderBoardPatterns(IdUser string) string {
	return "leaderboard:*"
}

func LeaderBoardPattern(IdGame string) string {
	return "leaderboard:" + IdGame
}

// SESSION
func AllSessionPatterns(IdUser string) string {
	return "session:" + IdUser + "*"
}

func SessionPattern(IdGame string, IdUser string) string {
	return "session:" + IdUser + ":" + IdGame
}