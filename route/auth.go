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
	auth.Put("/user/:user_id", controller.EditUser)
	auth.Put("/profile-picture/:user_id", controller.EditProfilePicture)

	// article
	auth.Post("/project", controller.CreateProject)
	auth.Get("/project/:user_id", controller.GetProjectByUser)
	auth.Post("/content/:article_id", controller.CreateContent)
	auth.Post("/post", controller.PostArticle)
	auth.Get("/content/:article_id", controller.GetFullContentDetail)
	auth.Get("/edit/:article_id", controller.GetContentForEdit)
	auth.Delete("/content/:article_id", controller.DeleteContent)
	auth.Delete("/project/:article_id", controller.DeleteFullArticle)
	auth.Put("/project/:article_id", controller.EditProject)
	auth.Put("/project/thumbnail/:article_id", controller.EditThumbnail)

	// comment
	auth.Post("/comment", controller.CommentArticle)
	auth.Post("/reply", controller.ReplyComment)
	auth.Get("/comment/:article_id", controller.GetCommentByArticle)
	auth.Get("/reply/:comment_id", controller.GetReplyByComment)

	// category
	auth.Post("/category", controller.CreateCategory)
}
