package validator

import (
	"errors"
	"net/url"
	"regexp"
)

type Validator interface {
	IsValid(key interface{}) bool
}

var (
	ErrBadFormat = errors.New("invalid format")
	emailRegexp  = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	phoneRegexp  = regexp.MustCompile(`^(?:(?:\(?(?:00|\+)([1-4]\d\d|[1-9]\d?)\)?)?[\-\.\ \\\/]?)?((?:\(?\d{1,}\)?[\-\.\ \\\/]?){0,})(?:[\-\.\ \\\/]?(?:#|ext\.?|extension|x)[\-\.\ \\\/]?(\d+))?$`)
)

type EmailValidator struct{}

func (v EmailValidator) IsValid(key interface{}) bool {
	emailStr, ok := key.(string)
	if !ok {
		return false
	}
	return emailRegexp.MatchString(emailStr)
}

type PhoneNumberValidator struct{}

func (v PhoneNumberValidator) IsValid(phoneNumber interface{}) bool {
	phoneStr, ok := phoneNumber.(string)
	if !ok {
		return false
	}
	return phoneRegexp.MatchString(phoneStr)
}

type URLValidator struct{}

func (v URLValidator) IsValid(URL interface{}) bool {
	URLStr, ok := URL.(string)
	if !ok {
		return false
	}

	_, err := url.ParseRequestURI(URLStr)
	if err != nil {
		return false
	}
	return true
}

const (
	minPasswordLength = 8
)

type PasswordValidator struct{}

func (v PasswordValidator) IsValid(password interface{}) bool {
	return false
}
