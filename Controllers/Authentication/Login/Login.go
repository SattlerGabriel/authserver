package LoginController

import (
	"errors"
	"tama/backend/Lib/Encryptor"
	"tama/backend/Lib/UserHandler"
	AuthStructs "tama/backend/Structs"
)

func Login(user AuthStructs.User) (AuthStructs.User, int, error) {
	if user.Password == "" || (user.Username == "" || user.Email == "") {
		return AuthStructs.User{}, 401, errors.New("fields can't be empty")
	}

	foundUser, err := UserHandler.GetUser(user.Username)
	if err != nil {
		return AuthStructs.User{}, 500, err
	}

	if foundUser.Username == "" {
		return AuthStructs.User{}, 500, errors.New("password or username don't match")
	}

	hashedPassword := Encryptor.ShaHashing(user.Password)
	if hashedPassword != foundUser.Password {
		return AuthStructs.User{}, 402, errors.New("password or username don't match")
	}

	if !foundUser.EmailConfirmed {
		return AuthStructs.User{}, 403, errors.New("user has not confirmed it's email")
	}

	return foundUser, 200, nil
}
