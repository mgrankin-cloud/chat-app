package utils

import (
    "crypto/sha256"
    "encoding/hex"
)

func HashPassword(password string) string {
    hash := sha256.Sum256([]byte(password))
    return hex.EncodeToString(hash[:])
}

func IsValidEmail(email string) bool {
    // тут написать логику проверки email
    return true
}