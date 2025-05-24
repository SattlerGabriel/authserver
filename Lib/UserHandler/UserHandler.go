package UserHandler

import (
	"errors"
	"tama/backend/Lib/Db"
	AuthStructs "tama/backend/Structs"

	"github.com/gofiber/fiber/v3/log"
)

func GetUser(username string) (AuthStructs.User, error) {
	db := Db.GetConnection()

	user := AuthStructs.User{}

	result := db.First(&user, "username = ?", username)
	if result.Error != nil {
		log.Error("Error fetching user", result.Error)
		return AuthStructs.User{}, errors.New("no user found")
	} else {
		return user, nil
	}
}

func GetUserByEmail(email string) (AuthStructs.User, error) {
	db := Db.GetConnection()

	user := AuthStructs.User{}

	result := db.First(&user, "email = ?", email)
	if result.Error != nil {
		log.Error("Error fetching user", result.Error)
		return AuthStructs.User{}, errors.New("no user found")
	} else {
		return user, nil
	}
}

func CreateUser(user AuthStructs.User) (AuthStructs.User, error) {
	db := Db.GetConnection()

	r := db.Create(&user)
	if r.Error != nil {
		log.Error(r.Error)
		return AuthStructs.User{}, r.Error
	}

	user, err := GetUser(user.Username)
	if err != nil {
		log.Error(err)
		return AuthStructs.User{}, err
	}

	return user, nil
}

func UpdateUser(user AuthStructs.User, updates interface{}) (AuthStructs.User, error) {
	db := Db.GetConnection()

	r := db.Model(AuthStructs.User{}).Where("username = ?", user.Username).Updates(updates)
	if r.Error != nil {
		log.Error(r.Error)
		return AuthStructs.User{}, r.Error
	}

	user, err := GetUser(user.Username)
	if err != nil {
		log.Error(err)
		return AuthStructs.User{}, err
	}

	return user, nil
}

func GetConfirmationCode(username string) (string, error) {
	db := Db.GetConnection()

	var code string

	result := db.First(&AuthStructs.User{}).Select("email_verification_code").Where("username = ?", username).Scan(&code)
	if result.Error != nil {
		log.Error(result.Error)
		return "", result.Error
	}

	return code, nil
}

func GetPasswordResetCode(username string) (string, error) {
	db := Db.GetConnection()

	var code string

	result := db.First(&AuthStructs.User{}).Select("password_change_code").Where("username = ?", username).Scan(&code)
	if result.Error != nil {
		log.Error(result.Error)
		return "", result.Error
	}

	return code, nil
}
