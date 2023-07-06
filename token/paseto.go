package token

import (
	"context"
	"os"
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

func GetIdFromToken(token string, key paseto.V4SymmetricKey) (string, bool) {
	parse := paseto.Parser{}
	decrypted, _ := parse.ParseV4Local(key, token, nil)
	if decrypted == nil {
		return "", false
	}
	IdUserString, err := decrypted.GetString("id-user")
	if err != nil {
		return "", false
	}
	return IdUserString, true
}

func (maker *PasetoMaker) VerifyUser(token string, client *redis.Client) (string, bool) {
	IdUserString, ok := GetIdFromToken(token, maker.SymmetricKey)

	if !ok {
		return "", false
	}
	_, err := client.Exists(context.Background(), share.UserPattern(IdUserString)).Result()
	if err != nil {
		return "", false
	}
	return IdUserString, true
}