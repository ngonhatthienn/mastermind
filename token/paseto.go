package token

import (
	"context"
	"fmt"
	"time"

	"intern2023/share"

	"aidanwoods.dev/go-paseto"
	"github.com/redis/go-redis/v9"
)
type PasetoMaker struct{
	SymmetricKey paseto.V4SymmetricKey
}

func NewPasetoMaker()(PasetoMaker, error) {
	maker := &PasetoMaker{
		SymmetricKey: paseto.NewV4SymmetricKey(),

	}
	return *maker, nil
}

func (maker *PasetoMaker)CreateToken(IdUser string) string {
	token := paseto.NewToken()

	token.SetIssuedAt(time.Now())
	token.SetNotBefore(time.Now())
	token.SetExpiration(time.Now().Add(2 * time.Hour))

	// Set the user id.
	token.SetString("id-user", IdUser)

	// Encrypt the token.
	encrypted := token.V4Encrypt(maker.SymmetricKey, nil)
	
	// Return the encrypted token.
	return encrypted
}

func getIdFromToken(token string, key paseto.V4SymmetricKey) (string, error) {
	parse := paseto.Parser{}
	decrypted, _ := parse.ParseV4Local(key, token, nil)
	return decrypted.GetString("id-user")
}

func (maker *PasetoMaker)VerifyUser(token string, client *redis.Client) (string, error) {
	IdUserString, err := getIdFromToken(token, maker.SymmetricKey)
	if err != nil {
		return "", err
	}
	_, err = client.Exists(context.Background(), share.UserPattern(IdUserString)).Result()
	if err != nil {
		return "", err
	}
	return IdUserString, nil
}

func Example() {
	publicKey, _ := paseto.NewV4AsymmetricPublicKeyFromHex("1eb9dbbbbc047c03fd70604e0071f0987e16b28b757225c11f00415d0e20b1a2")
	signed := "v4.public.eyJkYXRhIjoidGhpcyBpcyBhIHNpZ25lZCBtZXNzYWdlIiwiZXhwIjoiMjAyMi0wMS0wMVQwMDowMDowMCswMDowMCJ9v3Jt8mx_TdM2ceTGoqwrh4yDFn0XsHvvV_D0DtwQxVrJEBMl0F2caAdgnpKlt4p7xBnx1HcO-SPo8FPp214HDw.eyJraWQiOiJ6VmhNaVBCUDlmUmYyc25FY1Q3Z0ZUaW9lQTlDT2NOeTlEZmdMMVc2MGhhTiJ9"
	parser := paseto.NewParserWithoutExpiryCheck()
	token, _ := parser.ParseV4Public(publicKey, signed, nil) // this will fail if parsing failes, cryptographic checks fail, or validation rules fail
	thisToken := token
	fmt.Println(thisToken.GetString("id-user"))
}
