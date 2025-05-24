package RegisterController

import (
	"errors"
	"tama/backend/Lib/CodeGenerator"
	"tama/backend/Lib/Db"
	"tama/backend/Lib/Encryptor"
	"tama/backend/Lib/UserHandler"
	Validators "tama/backend/Lib/Validators"
	AuthStructs "tama/backend/Structs"
	"time"

	"github.com/gofiber/fiber/v3/log"
	"gorm.io/gorm"
)

func Register(user AuthStructs.User) (int, error) {
	valid := validateUsername(user.Username)

	if !valid {
		return 400, errors.New("invalid username")
	}

	user.ID = CodeGenerator.GetUUID()
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	user.EmailConfirmed = false
	user.EmailVerificationCode = CodeGenerator.RandStringRunes(6)
	user.Deleted = false
	user.CanChangePassword = false
	user.PasswordChangeCode = ""
	user.Password = Encryptor.ShaHashing(user.Password)
	user.Role = 0

	UserHandler.CreateUser(user)

	return 200, nil
}

func ConfirmEmail(user AuthStructs.User) (int, error) {
	foundUser, err := UserHandler.GetUserByEmail(user.Email)
	if err != nil {
		return 500, err
	}

	if foundUser.EmailConfirmed {
		return 400, errors.New("user is already confirmed")
	}

	if user.EmailVerificationCode != foundUser.EmailVerificationCode {
		return 400, errors.New("codes do not match")
	}

	foundUser.EmailConfirmed = true
	updates := map[string]interface{}{
		"EmailConfirmed": true,
	}
	_, err = UserHandler.UpdateUser(foundUser, updates)
	if err != nil {
		return 400, err
	}

	return 200, nil
}

func validateUsername(username string) bool {
	valid, err := Validators.ValidateString(username)
	if err != nil {
		log.Error("Invalid username: %v", err)
		return false
	}

	if !valid {
		log.Error("Username contains invalid characters")
		return false
	}

	db := Db.GetConnection()
	isInUse := db.First(&AuthStructs.User{}, "username = ?", username)
	if isInUse.Error != nil {
		if errors.Is(isInUse.Error, gorm.ErrRecordNotFound) {
			return true
		} else {
			log.Error("Error fetching user", isInUse.Error)
			return false
		}
	} else {
		return false
	}
}
