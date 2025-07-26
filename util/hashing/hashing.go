package hashing

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"slices"
)

func GeneratePassword(password string) (string, error) {
	HashedPass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash pass %v", err)
	}
	return string(HashedPass), nil
}

// Хеширует id юзеров комнаты чтоб потом легче было искать
func HashUsersForRoomUnique(ids []int32) string {
	slices.Sort(ids)
	var buf bytes.Buffer
	for _, id := range ids {
		fmt.Fprintf(&buf, "%d,", id)
	}
	hash := sha256.Sum256(buf.Bytes())
	return hex.EncodeToString(hash[:])
}
