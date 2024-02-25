package main

import (
	"server-article/route"
	s "server-article/service"
	"server-article/utils"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New(fiber.Config{
		Prefork:       true,
		CaseSensitive: true,
		StrictRouting: true,
		ServerHeader:  "articleing: 1.0.0",
		AppName:       "articleing",
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Status(fiber.StatusBadRequest).JSON(s.GlobalErrorHandlerResp{
				Success: false,
				Message: err.Error(),
			})
		},
	})

	// route selection
	route.GuestRoute(app)
	route.AuthRoute(app)

	port := "127.0.0.1:" + utils.GetEnv("PORT")
	app.Listen(port)
}
