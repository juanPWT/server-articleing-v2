package route

import (
	"server-article/controller"
	"server-article/utils"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
)

func AuthRoute(app *fiber.App) {
	auth := app.Group("/v1/auth", jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: []byte(utils.GetEnv("SECRET_TOKEN_JWT"))},
	}))

	// user
	auth.Get("/user", controller.GetUser)
	auth.Get("/user/logout", controller.Logout)

	// article
	auth.Get("/articles", controller.GetAllArticle)

	// category
	auth.Post("/category", controller.CreateCategory)
	auth.Get("/categories", controller.GetAllCategory)

}
