package auth

import "github.com/alexedwards/argon2id"

func HashPassword(password string) (string, error) {
	hashPass, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		return "", err
	}
	return hashPass, nil
}

func CheckPasswordHash(password, hash string) (bool, error) {
	return argon2id.ComparePasswordAndHash(password, hash)
}
