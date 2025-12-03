package format

import (
	"github.com/Drengin6306/ZeroBank/pkg/errorx"
	"github.com/nyaruka/phonenumbers"
)

func Format(raw string) (string, error) {
	defaultCountry := "CN"
	num, err := phonenumbers.Parse(raw, defaultCountry)
	if err != nil {
		return "", err
	}
	if !phonenumbers.IsValidNumber(num) {
		return "", errorx.NewErrorWithMsg(errorx.ErrInvalidParams, "invalid phone number")
	}
	return phonenumbers.Format(num, phonenumbers.E164), nil
}
