package PasswordReset

import (
	"errors"
	"fmt"
	"tama/backend/Lib/CodeGenerator"
	"tama/backend/Lib/Encryptor"
	"tama/backend/Lib/Mailer"
	"tama/backend/Lib/UserHandler"
	AuthStructs "tama/backend/Structs"
)

func SendPasswordCode(user AuthStructs.User) (int, error) {
	foundUser, err := UserHandler.GetUserByEmail(user.Email)
	if err != nil {
		return 500, err
	}

	code := CodeGenerator.RandStringRunes(6)
	subject := "Password change request"
	body := fmt.Sprintf("Password change code: %v", code)

	_, err = UserHandler.UpdateUser(foundUser, map[string]interface{}{"PasswordChangeCode": code})
	if err != nil {
		return 500, err
	}

	err = Mailer.SendEmail(foundUser.Email, subject, body)
	if err != nil {
		return 500, err
	}

	return 200, nil
}

func ValidatePasswordChangeCode(user AuthStructs.User) (int, error) {
	foundUser, err := UserHandler.GetUserByEmail(user.Email)
	if err != nil {
		return 500, err
	}

	if user.PasswordChangeCode != foundUser.PasswordChangeCode {
		return 401, errors.New("codes do not match")
	}

	_, err = UserHandler.UpdateUser(foundUser, map[string]interface{}{"CanChangePassword": true})

	return 200, nil
}

func UpdatePassword(user AuthStructs.User) (int, error) {
	foundUser, err := UserHandler.GetUserByEmail(user.Email)
	if err != nil {
		return 500, err
	}

	if !foundUser.CanChangePassword {
		return 401, errors.New("cant change password")
	}

	updates := map[string]interface{}{
		"Password": Encryptor.ShaHashing(user.Password),
	}
	_, err = UserHandler.UpdateUser(foundUser, updates)
	if err != nil {
		return 500, err
	}

	return 200, nil
}
