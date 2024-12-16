package otp

import (
	"crypto/rand"
	"encoding/base32"
	"fmt"
	"log"
	"net/smtp"
	"time"
)

type Backend struct {
}

func GenerateOTP() string {
	otp := make([]byte, 4)
	_, err := rand.Read(otp)
	if err != nil {
		log.Fatalf("Failed to generate OTP: %v", err)
	}
	otpInt := int(base32.StdEncoding.EncodeToString(otp)[0] & 0xF)
	return fmt.Sprintf("%04d", otpInt)
}

func SendEmail(to, subject, body string) error {
	from := "greengranchat@gmail.com"
	password := "androidchat11082024"
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	message := []byte("To: " + to + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"\r\n" +
		body + "\r\n")

	auth := smtp.PlainAuth("", from, password, smtpHost)
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, message)

	if err != nil {
		return err
	}
	return nil
}

var otpStore = make(map[string]string)

func StoreOTP(email, otp string, expiration time.Duration) {
	otpStore[email] = otp
	time.AfterFunc(expiration, func() {
		delete(otpStore, email)
	})
}

func VerifyOTP(email, otp string) bool {
	storedOTP, exists := otpStore[email]
	if !exists {
		return false
	}
	return storedOTP == otp
}