package helper

import "github.com/google/uuid"

func GenerateToken(userId string) (string, error) {
	// Generate a basic token
	token, err := uuid.NewUUID()
	if err != nil {
		return "", err
	}
	return token.String(), nil
}
