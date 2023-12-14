package helper

import (
	"net/mail"
)

func Panic(err error) {
	if err != nil {
		panic(err)
	}
}

func EmailValidator(email string) error {
	_, err := mail.ParseAddress(email)
	return err
}
