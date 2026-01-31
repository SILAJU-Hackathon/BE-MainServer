package utils

import (
	"fmt"
	"math/rand"
	"net/smtp"
	"os"
	"time"
)

func GenerateOTP() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return fmt.Sprintf("%06d", r.Intn(1000000))
}

func SendOTP(email, otp string) error {
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")
	smtpEmail := os.Getenv("SMTP_EMAIL")
	smtpPassword := os.Getenv("SMTP_PASSWORD")

	if smtpHost == "" || smtpPort == "" || smtpEmail == "" || smtpPassword == "" {
		fmt.Printf("----------------------------------------------------------------\n")
		fmt.Printf("SENDING OTP TO: %s\n", email)
		fmt.Printf("OTP CODE: %s\n", otp)
		fmt.Printf("----------------------------------------------------------------\n")
		return nil
	}

	auth := smtp.PlainAuth("", smtpEmail, smtpPassword, smtpHost)

	subject := "Your OTP Verification Code"
	body := fmt.Sprintf(`
Hello,

Your OTP verification code is: %s

This code will expire in 5 minutes. Do not share this code with anyone.

Best regards,
Dinacom Team
`, otp)

	msg := []byte(fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: text/plain; charset=UTF-8\r\n\r\n%s",
		smtpEmail, email, subject, body))

	addr := fmt.Sprintf("%s:%s", smtpHost, smtpPort)
	err := smtp.SendMail(addr, auth, smtpEmail, []string{email}, msg)
	if err != nil {
		fmt.Printf("Failed to send OTP email: %v\n", err)
		return err
	}

	fmt.Printf("OTP sent successfully to: %s\n", email)
	return nil
}
