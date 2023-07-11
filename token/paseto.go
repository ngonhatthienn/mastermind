package token

import (
	"context"
	"os"
	"time"
	"encoding/json"
	"intern2023/share"
	pb "intern2023/pb"
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

func IsTokenExpired(decrypted paseto.Token) bool {
	// Get the token's expiration time
	expirationTime, err := decrypted.GetExpiration()
	if err != nil {
		return true
	}
	// Check if the expiration time has passed
	return time.Now().After(expirationTime)
}

func (maker *PasetoMaker) CreateToken(IdUserString string) string {
	token := paseto.NewToken()

	token.SetIssuedAt(time.Now())
	token.SetNotBefore(time.Now())
	token.SetExpiration(time.Now().Add(2 * time.Hour))

	// Set the user id.
	token.SetString("id-user", IdUserString)

	// Encrypt the token.
	encrypted := token.V4Encrypt(maker.SymmetricKey, nil)

	// Return the encrypted token.
	return encrypted
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

func (maker *PasetoMaker) DecryptedToken(token string) (*paseto.Token, bool) {
	parse := paseto.Parser{}
	decrypted, err := parse.ParseV4Local(maker.SymmetricKey, token, nil)
	if err != nil || decrypted == nil {
		return nil, false
	}
	return decrypted, true
}

// We should verify user in token
func Authentication(role string, IdUserString string, client *redis.Client) bool {
	val, err := client.Get(context.Background(), share.UserPattern(IdUserString)).Result()
	if err != nil || val == "" {
		return false
	}
	var userData *pb.User
		_ = json.Unmarshal([]byte(val), &userData)
	if(role != userData.Role) {
		return false
	}
	return true
}

func (maker *PasetoMaker) VerifyUser(decrypted *paseto.Token, client *redis.Client) (string, bool) {
	IdUserString, ok := GetUserIdFromToken(decrypted)
	if !ok {
		return "", false
	}
	val, err := client.Get(context.Background(), share.UserPattern(IdUserString)).Result()
	if err != nil || val == "" {
		return "", false
	}
	return IdUserString, true
}
