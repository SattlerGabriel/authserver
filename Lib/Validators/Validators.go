package Validators

import (
	"errors"
	"regexp"
	"strings"
)

func ValidateString(text string) (bool, error) {
	text = strings.ReplaceAll(text, " ", "")
	text = strings.ReplaceAll(text, "	", "")
	if text == "" {
		return false, errors.New("error parsing username")
	}

	match, _ := regexp.MatchString("^[A-Za-z0-9][A-Za-z0-9_-]*$", text)
	if !match {
		return false, errors.New("error regex parsing username")
	}

	return true, nil
}
