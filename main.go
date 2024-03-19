package main

import (
	"server-article/config"
	"server-article/route"
	s "server-article/service"
	"server-article/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
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

	// file
	app.Static("/public", "./public", fiber.Static{
		Compress: true,
	})

	// cors
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:3000",
		AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH",
		AllowHeaders:     "Origin,Content-Type,Accept,Authorization,X-Requested-With,Access-Control-Allow-Origin,Access-Control-Allow-Headers,Access-Control-Allow-Methods,Access-Control-Allow-Credentials,Access-Control-Max-Age,Access-Control-Expose-Headers,Access-Control-Request-Headers,Content-Length,Accept-Language,Accept-Encoding,Connection",
		AllowCredentials: true,
	}))

	// migration
	config.Migrate()

	// route selection
	route.GuestRoute(app)
	route.AuthRoute(app)

	port := "127.0.0.1:" + utils.GetEnv("PORT")
	app.Listen(port)
}
