package token

import (
	"fmt"
	"time"

	"aidanwoods.dev/go-paseto"
)

func GenToken(IdUser string) {
	token := paseto.NewToken()

	token.SetIssuedAt(time.Now())
	token.SetNotBefore(time.Now())
	token.SetExpiration(time.Now().Add(2 * time.Hour))

	// Set the user id.
	token.SetString("id-user", IdUser)

	// Encrypt the token.
	key := paseto.NewV4SymmetricKey()
	encrypted := token.V4Encrypt(key, nil)

	// Print the encrypted token.
	fmt.Println(encrypted)
	parse := paseto.NewParser()
	decrypted, _ := paseto.Parser.ParseV4Local(parse, key, "", nil)
	fmt.Println(decrypted)
}

func Example() {
	publicKey, _ := paseto.NewV4AsymmetricPublicKeyFromHex("1eb9dbbbbc047c03fd70604e0071f0987e16b28b757225c11f00415d0e20b1a2")
	signed := "v4.public.eyJkYXRhIjoidGhpcyBpcyBhIHNpZ25lZCBtZXNzYWdlIiwiZXhwIjoiMjAyMi0wMS0wMVQwMDowMDowMCswMDowMCJ9v3Jt8mx_TdM2ceTGoqwrh4yDFn0XsHvvV_D0DtwQxVrJEBMl0F2caAdgnpKlt4p7xBnx1HcO-SPo8FPp214HDw.eyJraWQiOiJ6VmhNaVBCUDlmUmYyc25FY1Q3Z0ZUaW9lQTlDT2NOeTlEZmdMMVc2MGhhTiJ9"
	parser := paseto.NewParserWithoutExpiryCheck()
	_, _ = parser.ParseV4Public(publicKey, signed, nil) // this will fail if parsing failes, cryptographic checks fail, or validation rules fail
}
