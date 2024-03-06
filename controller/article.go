package controller

import (
	"fmt"
	"server-article/model"
	"server-article/utils"

	"github.com/gofiber/fiber/v2"
)

func GetAllArticle(c *fiber.Ctx) error {
	article := new([]model.Article)
	errGet := db.Joins("User").Joins("Category").Find(&article)

	if errGet.Error != nil {
		return utils.ResObject(c, fiber.StatusBadRequest, "cannot get article", nil)
	}

	return utils.ResObject(c, fiber.StatusOK, "success get article", article)
}

func CreateProject(c *fiber.Ctx) error {
	r := new(model.CreteProject)
	if err := c.BodyParser(r); err != nil {
		return utils.ResObject(c, fiber.StatusBadRequest, "cannot parse input", nil)
	}

	// validate body parser request
	if errs := myValidation.Validate(r); len(errs) > 0 && errs[0].Error {
		errsMsgs := make([]string, 0)

		for _, err := range errs {
			errsMsgs = append(errsMsgs, fmt.Sprintf(
				"[%s]: '%v' | Need tobe '%s'",
				err.FailedField,
				err.Value,
				err.Tag,
			))
		}

		return utils.ResObject(c, fiber.StatusBadRequest, "validation error", errsMsgs)
	}

	// if user is exist
	userExist := db.Where("id = ?", r.User_id).First(&model.User{})
	if userExist.RowsAffected == 0 {
		return utils.ResObject(c, fiber.StatusBadGateway, "user not found", nil)
	}

	// if category is exist
	var user model.User
	categoryExist := db.Where("id = ?", r.Category_id).First(&user)
	if categoryExist.RowsAffected == 0 {
		return utils.ResObject(c, fiber.StatusBadGateway, "category not found", nil)
	}

	// create project article
	newProject := &model.Article{
		User_id:      r.User_id,
		Category_id:  r.Category_id,
		Title:        r.Title,
		Introduction: r.Introduction,
		Thumbnail:    r.Thumbnail,
	}

	errCreate := db.Create(&newProject)
	if errCreate.Error != nil {
		return utils.ResObject(c, fiber.StatusBadRequest, "cannot create project", nil)
	}

	return utils.ResObject(c, fiber.StatusCreated, "success create project "+user.Username, newProject)
}

func GetProjectByUser(c *fiber.Ctx) error {

	user_id := c.Params("user_id")

	if user_id == "" {
		return utils.ResObject(c, fiber.StatusBadRequest, "user id cannot be empty", nil)
	}

	// get project article by user
	project := new([]model.Article)
	errGet := db.Where("user_id = ?", user_id).Joins("User").Joins("Category").Find(&project)
	if errGet.Error != nil {
		return utils.ResObject(c, fiber.StatusBadRequest, "cannot get project", nil)
	}

	return utils.ResObject(c, fiber.StatusOK, "success get project", project)

}

func CreateContent(c *fiber.Ctx) error {
	r := new(model.CreateContent)
	article_id := c.Params("article_id")

	if err := c.BodyParser(r); err != nil {
		return utils.ResObject(c, fiber.StatusBadRequest, "cannot parse input", nil)
	}

	// valdiate body parser request
	if errs := myValidation.Validate(r); len(errs) > 0 && errs[0].Error {
		errsMsgs := make([]string, 0)

		for _, err := range errs {
			errsMsgs = append(errsMsgs, fmt.Sprintf(
				"[%s]: '%v' | Need tobe '%s'",
				err.FailedField,
				err.Value,
				err.Tag,
			))
		}

		return utils.ResObject(c, fiber.StatusBadRequest, "validation error", errsMsgs)
	}

	// if article is exist
	var article model.Article
	articleExist := db.Where("id = ?", article_id).First(&article)
	if articleExist.RowsAffected == 0 {
		return utils.ResObject(c, fiber.StatusBadGateway, "article not found", nil)
	}

	// create content article
	newContent := &model.Body{
		Article_id: article.ID,
		Content:    r.Content,
	}

	errCreate := db.Create(&newContent)
	if errCreate.Error != nil {
		return utils.ResObject(c, fiber.StatusBadRequest, "cannot create content", nil)
	}

	return utils.ResObject(c, fiber.StatusCreated, "success create content", newContent)
}

func PostArticle(c *fiber.Ctx) error {
	article_id := c.Query("article_id")
	user_id := c.Query("user_id")

	if article_id == "" {
		return utils.ResObject(c, fiber.StatusBadRequest, "article id cannot be empty", nil)
	}

	if user_id == "" {
		return utils.ResObject(c, fiber.StatusBadRequest, "user id cannot be empty", nil)
	}

	// if article is exist
	var article model.Article
	articleExist := db.Where("id = ?", article_id).Where("user_id = ?", user_id).First(&article)
	if articleExist.RowsAffected == 0 {
		return utils.ResObject(c, fiber.StatusBadGateway, "article not found", nil)
	}

	// update article
	article.IsPost = true
	errUpdate := db.Save(&article)
	if errUpdate.Error != nil {
		return utils.ResObject(c, fiber.StatusBadRequest, "cannot update article", nil)
	}

	return utils.ResObject(c, fiber.StatusOK, "success post article", article)
}

func GetFullContentDetail(c *fiber.Ctx) error {
	article_id := c.Params("article_id")

	if article_id == "" {
		return utils.ResObject(c, fiber.StatusBadRequest, "article id cannot be empty", nil)
	}

	// get  article
	var articleDetail model.Article
	errGet := db.Where("is_post = ?", true).Joins("User").Joins("Category").First(&articleDetail, article_id)
	if errGet.RowsAffected == 0 {
		return utils.ResObject(c, fiber.StatusBadRequest, "cannot get article detail", nil)
	}

	// get content article
	var body []model.Body
	errGetContent := db.Where("article_id = ?", article_id).Find(&body)
	if errGetContent.Error != nil {
		return utils.ResObject(c, fiber.StatusBadRequest, "cannot get article detail", nil)
	}

	response := new(model.ArticleDetail)
	response.Article = articleDetail
	response.Body = body

	return utils.ResObject(c, fiber.StatusOK, "success get article detail", response)

}
