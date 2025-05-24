package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	CRouter "tama/backend/Controllers"
	"tama/backend/Lib/Db"
	AuthStructs "tama/backend/Structs"

	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	err = runMigrations()
	if err != nil {
		log.Fatal(err)
	}

	app := fiber.New()

	auth := app.Group("/auth")
	auth.Post("/register", CRouter.Register)
	auth.Put("/confirmemail", CRouter.ConfirmEmail)
	auth.Get("/sendpasswordchangecode", CRouter.SendPasswordCode)
	auth.Post("/requestpasswordchange", CRouter.RequestPasswordChange)
	auth.Put("/updatepassword", CRouter.UpdatePassword)
	auth.Post("/login", CRouter.Login)

	app.Get("/jwtTest", CRouter.JwtTest, jwtMiddleware())

	log.Fatal(app.Listen("localhost:3000"))
}

func runMigrations() error {
	db := Db.GetConnection()

	err := db.AutoMigrate(&AuthStructs.User{})
	if err != nil {
		return errors.New("failed to run automatic migration for User")
	}

	return nil
}

func jwtMiddleware() fiber.Handler {
	return func(c fiber.Ctx) error {
		auth := c.Get("Authorization")
		if auth == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Missing token"})
		}

		tokenStr := auth[len("Bearer "):]

		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			// Ensure signing method is HMAC
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fiber.ErrUnauthorized
			}
			return []byte(os.Getenv("TOKEN_SECRET")), nil
		})

		fmt.Print(token)

		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token"})
		}

		// Store user info in context
		claims := token.Claims.(jwt.MapClaims)
		c.Locals("username", claims["username"])
		c.Locals("role", claims["role"])

		return c.Next()
	}
}
