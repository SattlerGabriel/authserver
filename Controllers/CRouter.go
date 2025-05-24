package CRouter

import (
	"fmt"
	"os"
	LoginController "tama/backend/Controllers/Authentication/Login"
	"tama/backend/Controllers/Authentication/PasswordReset"
	RegisterController "tama/backend/Controllers/Authentication/Register"
	AuthStructs "tama/backend/Structs"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/log"
	"github.com/golang-jwt/jwt"
)

func Register(ctx fiber.Ctx) error {
	user := new(AuthStructs.User)

	if err := ctx.Bind().Body(user); err != nil {
		return err
	}

	code, err := RegisterController.Register(*user)
	if err != nil {
		return err
	}

	ctx.SendStatus(code)

	return nil
}

func Login(ctx fiber.Ctx) error {
	user := new(AuthStructs.User)

	if err := ctx.Bind().Body(user); err != nil {
		return err
	}

	existingUser, _, err := LoginController.Login(*user)
	if err != nil {
		return err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": existingUser.Username,
		"role":     existingUser.Role,
		"exp":      time.Now().Add(time.Hour * 72).Unix(),
	})

	var jwtSecret = []byte(os.Getenv("TOKEN_SECRET"))
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		log.Error(err)
		return ctx.Status(fiber.StatusInternalServerError).SendString("Could not login")
	}

	response := fmt.Sprintf(`{
		"user": "%v",
		"token": "%v"
	}`,
		AuthStructs.User{ID: existingUser.ID, Username: existingUser.Username, Email: existingUser.Email},
		tokenString)

	ctx.Set("Content-Type", "application/json")
	return ctx.SendString(response)
}

func ConfirmEmail(ctx fiber.Ctx) error {
	user := new(AuthStructs.User)
	if err := ctx.Bind().Body(user); err != nil {
		return err
	}

	code, err := RegisterController.ConfirmEmail(*user)
	if err != nil {
		return err
	}

	return ctx.SendStatus(code)
}

func SendPasswordCode(ctx fiber.Ctx) error {
	user := new(AuthStructs.User)
	if err := ctx.Bind().Body(user); err != nil {
		return err
	}

	code, err := PasswordReset.SendPasswordCode(*user)
	if err != nil {
		return err
	}

	return ctx.SendStatus(code)
}

func RequestPasswordChange(ctx fiber.Ctx) error {
	user := new(AuthStructs.User)
	if err := ctx.Bind().Body(user); err != nil {
		return err
	}

	code, err := PasswordReset.ValidatePasswordChangeCode(*user)
	if err != nil {
		return err
	}

	return ctx.SendStatus(code)
}

func UpdatePassword(ctx fiber.Ctx) error {
	user := new(AuthStructs.User)
	if err := ctx.Bind().Body(user); err != nil {
		return err
	}

	code, err := PasswordReset.UpdatePassword(*user)
	if err != nil {
		return err
	}

	return ctx.SendStatus(code)
}

func JwtTest(ctx fiber.Ctx) error {
	username := ctx.Locals("username")
	role := ctx.Locals("role")
	return ctx.SendString(fmt.Sprintf("User: %v, Role: %v", username, role))
}
