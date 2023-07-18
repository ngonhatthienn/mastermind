package token

import (
	"context"
	"os"
	"strconv"
	"time"
	"intern2023/share"

	"aidanwoods.dev/go-paseto"
	"github.com/golobby/dotenv"
	"github.com/redis/go-redis/v9"
)

type PasetoMaker struct {
	SymmetricKey paseto.V4SymmetricKey
}
type pasetoConfig struct {
	SymmetricKeyHex string `env:"SYMMETRIC_KEY_HEX"`
}

func NewPasetoMaker() (PasetoMaker, error) {
	config := pasetoConfig{}
	file, err := os.Open("app.env")
	err = dotenv.NewDecoder(file).Decode(&config)
	if err != nil {
		panic(err)
	}
	key, _ := paseto.V4SymmetricKeyFromHex(config.SymmetricKeyHex)
	maker := &PasetoMaker{
		SymmetricKey: key,
	}
	return *maker, nil
}

func GetUserIdFromToken(decrypted *paseto.Token) (string, bool) {
	if decrypted == nil {
		return "", false
	}
	IdUserString, err := decrypted.GetString("id-user")
	if err != nil {
		return "", false
	}
	return IdUserString, true
}
func GetSessionIdFromToken(decrypted *paseto.Token) (string, bool) {
	if decrypted == nil {
		return "", false
	}
	IdSessionString, err := decrypted.GetString("id-session")
	if err != nil {
		return "", false
	}
	return IdSessionString, true
}


func (maker *PasetoMaker) CreateToken(IdUserString string, userRole string) (string, string) {
	token := paseto.NewToken()

	token.SetIssuedAt(time.Now())
	token.SetNotBefore(time.Now())
	token.SetExpiration(time.Now().Add(3 * time.Hour))
	token.SetSubject(userRole)
	// Set the user id.
	token.SetString("id-user", IdUserString)
	// Set the session id.
	IdSession := share.CreateRandomNumber(10000, 99999)
	IdSessionString := strconv.Itoa(IdSession)
	token.SetString("id-session", IdSessionString)

	// Encrypt the token.
	encrypted := token.V4Encrypt(maker.SymmetricKey, nil)

	// Return the encrypted token.
	return encrypted, IdSessionString
}

func (maker *PasetoMaker) DecryptedToken(token string) (*paseto.Token, bool) {
	parse := paseto.NewParser()
	decrypted, err := parse.ParseV4Local(maker.SymmetricKey, token, nil)
	if err != nil || decrypted == nil {
		return nil, false
	}
	return decrypted, true
}

// We should verify user in token

func (maker *PasetoMaker) CheckExistUser(decrypted *paseto.Token, client *redis.Client) (string, bool) {
	IdUserString, ok := GetUserIdFromToken(decrypted)
	if !ok {
		return "", false
	}
	val, err := client.Get(context.Background(), share.UserPatternValue(IdUserString)).Result()
	if err != nil || val == "" {
		return "", false
	}
	return IdUserString, true
}
func (maker *PasetoMaker) CheckExactSession(decrypted *paseto.Token, client *redis.Client) (string, bool) {
	IdUserString, _ := GetUserIdFromToken(decrypted)
	IdSessionString, ok := GetSessionIdFromToken(decrypted)
	if !ok {
		return "", false
	}
	IdSession, err := client.Get(context.Background(), share.UserPatternSession(IdUserString)).Result()
	if err != nil || IdSession == "" {
		return "", false
	}
	if(IdSession != IdSessionString ) {
		return "", false
	}
	return IdSessionString, true
}

func (maker *PasetoMaker) Authentication(decrypted *paseto.Token, permission string) bool {
	if decrypted == nil {
		return false
	}
	UserRole, err := decrypted.GetSubject()
	if err != nil {
		return false
	}
	if permission == UserRole {
		return true
	} else if permission == "none" {
		return true
	}
	return false
}
