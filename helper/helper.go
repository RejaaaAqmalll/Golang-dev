package helper

import (
	"fmt"
	"math/rand"
	"time"

	"gopkg.in/gomail.v2"
)

// Const auth email
const CONFIG_SMTP_HOST = "smtp.gmail.com"
const CONFIG_SMTP_PORT = 587
const CONFIG_SENDER_NAME = "PT. Ada Ide Langsung Jalan <absensmart@gmail.com>"
const CONFIG_AUTH_EMAIL = "absensmart@gmail.com"
const CONFIG_AUTH_PASSWORD = "kwihycajcmnuuevg"

// Generate a random token
func GenerateCode() string {
	const lengthcode = 4
	const codechar = "0123456789"

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	r.Seed(time.Now().UnixNano())
	code := make([]byte, lengthcode)
	for i := range code {
		code[i] = codechar[rand.Intn(len(codechar))]
	}
	return string(code)
}

// Send code to email
func Sendemail(email string, code string) error {
	//  menangkap email dan kode authy
	//Set send message
	m := gomail.NewMessage()
	m.SetHeader("From", CONFIG_SENDER_NAME)
	m.SetHeader("To", email)
	m.SetHeader("Subject", "Password reset code")
	m.SetBody("text/plain", fmt.Sprintf("Your password reset code is: %s", code))

	d := gomail.NewDialer(
		CONFIG_SMTP_HOST,
		CONFIG_SMTP_PORT,
		CONFIG_AUTH_EMAIL,
		CONFIG_AUTH_PASSWORD,
	)
	if err := d.DialAndSend(m); err != nil {
		fmt.Println(err)
	}

	// Respon
	fmt.Println("Your code has been send")
	return nil
}
