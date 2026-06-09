package password

import "golang.org/x/crypto/bcrypt"

const DefaultCost = bcrypt.DefaultCost

func Hash(plain string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(plain), DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func Compare(hash, plain string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(plain))
}
