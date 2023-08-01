package handler

import (
	"strconv"
	"time"

	"google.golang.org/grpc/metadata"

	gameApp "intern2023/app"
	"intern2023/internal/gameLogic/repository"
	"intern2023/internal/gameLogic/service"
	"intern2023/token"
	pb "intern2023/pb/game"
	share "intern2023/share"
)

type GameLogicHandler struct {
	gameRepository        *repository.GameRepositoryImpl
	leaderBoardRepository *repository.LeaderBoardRepositoryImpl
	userRepository        *repository.UserRepositoryImpl
	sessionRepository     *repository.SessionRepositoryImpl

	pasetoMaker token.PasetoMaker
}

func NewService() *GameLogicHandler {
	gameRepository := repository.NewGameRepositoryImpl()
	leaderBoardRepository := repository.NewLeaderBoardRepositoryImpl()
	userRepository := repository.NewUserRepositoryImpl()
	sessionRepository := repository.NewSessionRepositoryImpl()
	

	pasetoMaker, _ := token.NewPasetoMaker()
	return &GameLogicHandler{gameRepository: gameRepository, leaderBoardRepository: leaderBoardRepository, userRepository: userRepository, sessionRepository: sessionRepository, pasetoMaker: pasetoMaker}


}

// GAME
func (h *GameLogicHandler) CreateGame(sizeGame int, GuessLimit int) {
	h.gameRepository.UpdateListGame(GuessLimit)
}

func (h *GameLogicHandler) ListGame() (int, []repository.Game) {
	// Check Any Games, if not, generate it
	h.gameRepository.CacheGameFromDB()
	// Get list game
	length, Games := h.gameRepository.GetListGame()
	return length, Games
}

func (h *GameLogicHandler) GetCurrent(IdUser int) (share.Status, *pb.GameReply) {
	// Check Any Games, if not, generate it
	h.gameRepository.CacheGameFromDB()
	_, IdGame := h.sessionRepository.CreateNewSession(int32(IdUser))
	Game, isExist := h.gameRepository.GetGame(IdGame)
	if !isExist {
		status := share.GenerateStatus(404, "Game")
		return status, &pb.GameReply{}
	}
	GameReply := pb.GameReply{XId: int32(Game.ID), GuessLimit: int32(Game.GuessLimit)}

	status := share.GenerateStatus(200, "Get current")
	return status, &GameReply
}

// Pick one game
func (h *GameLogicHandler) PickGame(IdUser int, IdGame int) (share.Status, *pb.GameReply) {
	// Check Any Games, if not, generate it
	h.gameRepository.CacheGameFromDB()
	// Check exist game
	if h.gameRepository.CheckExistGame(IdGame) == false {
		status := share.GenerateStatus(404, "Game")
		return status, &pb.GameReply{}
	}
	// Handle pick game
	h.sessionRepository.CreateSessionWithId(int32(IdUser), int32(IdGame))
	Game, isExist := h.gameRepository.GetGame(IdGame)
	if !isExist {
		status := share.GenerateStatus(404, "Game")
		return status, &pb.GameReply{}
	}
	GameReply := pb.GameReply{XId: int32(Game.ID), GuessLimit: int32(Game.GuessLimit)}
	status := share.GenerateStatus(200, "Pick game")
	return status, &GameReply
}

// Update Game
func (h *GameLogicHandler) UpdateGame(GuessLimit int) share.Status {
	h.gameRepository.UpdateListGame(GuessLimit)
	status := share.GenerateStatus(200, "Update Game")
	return status
}

// Play Game
func (h *GameLogicHandler) PlayGame(IdUser int, UserGuess string) (share.Status, int, []*pb.ListHistory) {
	// Check Any Games, if not, generate it
	h.gameRepository.CacheGameFromDB()

	// Check Exist session or not
	keySession, isExistKeySession := h.sessionRepository.GetKeySessionByUserID(IdUser)
	if !isExistKeySession {
		_, IdGame := h.sessionRepository.CreateNewSession(int32(IdUser))
		keySession = share.SessionPattern(strconv.Itoa(IdGame), strconv.Itoa(IdUser))
	}

	// Check if user already win or not
	isWin, _ := h.sessionRepository.GetSessionValue(keySession, "isWin").Bool()
	if isWin {
		status := share.GenerateStatus(200, "")
		status.Message = "You'd already won, please get another game"
		return status, 0, nil
	}
	guessLeft, _ := h.sessionRepository.GetSessionValue(keySession, "guessLeft").Int()
	if guessLeft == 0 {
		status := share.GenerateStatus(200, "")
		status.Message = "You lose!!!"
		return status, 0, nil
	} else {
		guessLeft--
		h.sessionRepository.SetSessionValue(keySession, "guessLeft", guessLeft)
	}

	IdGameString := share.GetKeyElement(keySession, 2)
	IdGame, _ := strconv.Atoi(IdGameString)
	getGameValue, isExist := h.gameRepository.GetGame(IdGame)
	if !isExist {
		status := share.GenerateStatus(404, "Game")
		return status, 0, nil
	}
	IdUserString := strconv.Itoa(IdUser)
	rightNumber, rightPosition := gameApp.OutputGame(UserGuess, getGameValue.Game)

	// If user win this game
	if rightNumber == rightPosition && rightNumber == 5 {
		h.sessionRepository.SetSessionValue(keySession, "isWin", true)

		// CalcScore
		// score := share.CalcScore(s.redisClient, keySession, guessLeft)
		timeStart, _ := h.sessionRepository.GetSessionValue(keySession, "timeStart").Int64()
		savedTime := time.Unix(timeStart, 0)
		diffInSeconds := 5000 - time.Now().Sub(savedTime).Seconds()
		// Get right and pos
		right, _ := h.sessionRepository.GetSessionValue(keySession, "sumRight").Int()
		pos, _ := h.sessionRepository.GetSessionValue(keySession, "sumPos").Int()
		score := int(diffInSeconds) + guessLeft*100 + (right+pos)*2

		_ = h.leaderBoardRepository.AddScore(IdUserString, IdGameString, int64(score))
		status := share.GenerateStatus(200, "")
		status.Message = "You win!!!"
		return status, guessLeft, nil
	}


	var listHistory []*pb.ListHistory
	listHistory, _ = h.sessionRepository.PushAndGetHistory(keySession, UserGuess, int32(rightNumber), int32(rightPosition))	
	
	status := share.GenerateStatus(200, "")
	status.Message = "Try your best !!!"
	return status, guessLeft, listHistory
}

// Hint Game
func (h *GameLogicHandler) HintGame(IdUser int, Type string) (share.Status, string) {
	// Check Any Games, if not, generate it
	h.gameRepository.CacheGameFromDB()
	keySession, isExistKeySession := h.sessionRepository.GetKeySessionByUserID(IdUser)
	if !isExistKeySession {
		status := share.GenerateStatus(404, "Session")
		return status, ""
	}

	IdGame := share.GetKeyElement(keySession, 2)
	IdGameInt,_ := strconv.Atoi(IdGame)
	Result, isExist := h.gameRepository.GetGame(IdGameInt)
	if !isExist {
		status := share.GenerateStatus(404, "Game")
		return status, ""
	}

	res, isSuccess := h.gameRepository.GenerateHint(Result.Game, Type)
	if !isSuccess {
		status := share.GenerateStatus(400, "Get hint")
		return status, ""
	}
	status := share.GenerateStatus(200, "")
	return status, res
}

// USER
func (h *GameLogicHandler) LogIn(Name string, Password string) (share.Status, int, string, bool) {
	IdUser, userRole, ok := h.userRepository.LogIn(Name, Password)
	if ok {
		status := share.GenerateStatus(200, "LogIn")
		return status, IdUser, userRole, ok
	}
	status := share.GenerateStatus(404, "User")
	return status, IdUser, userRole, ok
}

func (h *GameLogicHandler) CreateToken(IdUser int, userRole string) string {
	IdUserString := strconv.Itoa(IdUser)
	token, IdSession := h.pasetoMaker.CreateToken(IdUserString, userRole)
	h.sessionRepository.SetSessionID(share.UserPatternSession(IdUserString), IdSession)
	return token
}

func (h *GameLogicHandler) CreateUser(Fullname string, Username string, Password string, Email string, Role string) (int32, error) {
	Id := h.userRepository.CreateUser(Fullname, Username, Password, Email, Role) // Not in best practices
	return Id, nil
}

func (h *GameLogicHandler) ListUsers() ([]repository.User, error) {
	_,Users := h.userRepository.GetListUser()
	return Users, nil
}

// LEADERBOARD
func (h *GameLogicHandler) GetLeaderBoard(IdGame int, IdUser int, Size int, isAdmin bool) (share.Status, []repository.LeaderBoard, int32, string) {
	// Check Any Games, if not, generate it
	h.gameRepository.CacheGameFromDB()
	// Check exist game
	if h.gameRepository.CheckExistGame(IdGame) == false {
		status := share.GenerateStatus(404, "Id Game")
		return status, nil, 0, ""
	}

	IdUserString := strconv.Itoa(IdUser)
	var UserRank int32
	var UserScore string
	leaderboardData, err := h.leaderBoardRepository.GetLeaderboard(strconv.Itoa(IdGame), int64(Size), IdUserString)
	if err != nil || leaderboardData == nil {
		status := share.GenerateStatus(200, "")
		status.Message = "No user has won this game yet"
		return status, leaderboardData, UserRank, UserScore
	}
	status := share.GenerateStatus(200, "Get LeaderBoard")
	if !isAdmin {
		UserRank, UserScore = h.leaderBoardRepository.GetUserRank( strconv.Itoa(IdGame), strconv.Itoa(IdUser))
	}

	return status, leaderboardData, UserRank, UserScore
}

// AUTHORIZATION
func (h *GameLogicHandler) AuthorAndAuthn(md metadata.MD, permission string) (share.Status, int) {
	bearerToken := md.Get("authorization")
	if len(bearerToken) <= 0 {
		status := share.GenerateStatus(401, "")
		return status, 0
	}
	reqToken := share.GetTokenElement(bearerToken[0], 1)
	decryptedToken, decryptedOk := h.pasetoMaker.DecryptedToken(reqToken)
	if !decryptedOk {
		status := share.GenerateStatus(401, "Token")
		status.Message = "Invalid or Expired token "
		return status, 0
	}

	IdUserString, ok := token.GetUserIdFromToken(decryptedToken)
	if !ok {
		status := share.GenerateStatus(404, "User")
		return status, 0
	}
	IdUser,_ := strconv.Atoi(IdUserString)
	ok = h.userRepository.CheckExistUser(IdUser)
	if !ok {
		status := share.GenerateStatus(404, "User")
		return status, 0
	}
	// GRPC
	IdSessionString, ok := token.GetSessionIdFromToken(decryptedToken)
	userService := service.NewUserService()
	isExist := userService.CheckSessionId(IdUserString, IdSessionString)
	if !isExist {
		status := share.GenerateStatus(401, "Token")
		status.Message = "Expired token "
		return status, 0
	}

	isAuthn := h.pasetoMaker.Authorization(decryptedToken, permission)
	if !isAuthn {
		status := share.GenerateStatus(403, "")
		return status, 0
	}

	status := share.GenerateStatus(200, "")
	return status, IdUser
}
