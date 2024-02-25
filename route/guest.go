package route

import (
	"server-article/controller"

	"github.com/gofiber/fiber/v2"
)

func GuestRoute(app *fiber.App) {
	g := app.Group("/v1/guest")
	g.Post("/signup", controller.SignUp)
	g.Post("/signin", controller.SignIn)
	g.Get("verifyemail/:verification_code", controller.VerifyEmail)
}
