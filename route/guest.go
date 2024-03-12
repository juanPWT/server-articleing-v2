package route

import (
	"server-article/controller"

	"github.com/gofiber/fiber/v2"
)

func GuestRoute(app *fiber.App) {
	g := app.Group("/v1/guest")

	// user
	g.Post("/signup", controller.SignUp)
	g.Post("/signin", controller.SignIn)
	g.Get("verifyemail/:verification_code", controller.VerifyEmail)
	g.Post("/forgotpassword", controller.ForgotPassword)
	g.Post("/forgotresetpassword", controller.ForgotResetPassword)

	// articles
	g.Get("/articles", controller.GetAllArticle)

	// category
	g.Get("/categories", controller.GetAllCategory)
}
