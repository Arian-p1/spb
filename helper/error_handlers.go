package helper

import (
	"errors"
	"strings"
)



func Panic(err error) {
	if err != nil {
		panic(err)
	}
}

func EmailValidator(mail string) error {
	r := strings.Split(mail, "@")
	err := errors.New("email is wrong")
	if len(r) != 2 {
		return err
	}
	if strings.Count(r[1], ".") != 1 {
		return err
	}
	return nil
}
