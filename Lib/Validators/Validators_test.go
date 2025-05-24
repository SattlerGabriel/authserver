package Validators_test

import (
	Validators "tama/backend/Lib/Validators"
	"testing"
)

func TestValidateString(t *testing.T) {
	username := " I hate spaces & s*^pecial c h	aracters      "

	res, _ := Validators.ValidateString(username)
	if !res {
		t.Errorf("Expected res to be false was true")
	}
}
