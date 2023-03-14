package authentication

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

func RegisterAuth(mongoClient *mongo.Client, ctx context.Context, app *fiber.App) error {
	authRepo := NewMongoRepository(mongoClient)
	jwtHelper := NewJWTHelper()
	authService := NewService(authRepo, jwtHelper)
	authHandler := NewHandler(authService)

	app.Post("/login", func(c *fiber.Ctx) error {
		return authHandler.Login(ctx, c)
	})

	return nil
}

type AuthHandler struct {
	authService *Service
}

func NewHandler(authService *Service) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

func (a *AuthHandler) Login(ctx context.Context, c *fiber.Ctx) error {
	var userCredential UserCredential
	err := c.BodyParser(&userCredential)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
			"error":   err.Error(),
		})
	}

	jwt, err := a.authService.AuthenticateUser(ctx, &userCredential)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Invalid credentials",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Login successful",
		"jwt":     jwt,
	})

}
