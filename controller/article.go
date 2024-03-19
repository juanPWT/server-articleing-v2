package controller

import (
	"fmt"
	"server-article/model"
	"server-article/utils"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func GetAllArticle(c *fiber.Ctx) error {
	article := new([]model.Article)
	errGet := db.Where("is_post = ?", true).Joins("User").Joins("Category").Find(&article)

	if errGet.Error != nil {
		return utils.ResObject(c, fiber.StatusBadRequest, "cannot get article", nil)
	}

	return utils.ResObject(c, fiber.StatusOK, "success get article", article)
}

func GetArticleByCategory(c *fiber.Ctx) error {
	category_id := c.Params("category_id")
	if category_id == "" {
		return utils.ResObject(c, fiber.StatusBadRequest, "category cannot be empty", nil)
	}

	var article []model.Article
	errGet := db.Where("is_post = ?", true).Where("category_id = ?", category_id).Joins("User").Joins("Category").Find(&article)
	if errGet.RowsAffected == 0 {
		return utils.ResObject(c, fiber.StatusBadRequest, "zero result search by category", nil)
	}

	return utils.ResObject(c, fiber.StatusOK, "success get article by category", article)
}

func SearchArticle(c *fiber.Ctx) error {
	search := c.Query("search")

	if search == "" {
		return utils.ResObject(c, fiber.StatusBadRequest, "search cannot be empty", nil)
	}

	var article []model.Article
	errGet := db.Where("is_post = ?", true).Where("title LIKE ?", "%"+search+"%").Joins("User").Joins("Category").Find(&article)
	if errGet.RowsAffected == 0 {
		return utils.ResObject(c, fiber.StatusBadRequest, "zero result search", nil)
	}

	return utils.ResObject(c, fiber.StatusOK, "Success search article", article)
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
	var user model.User
	userExist := db.Where("id = ?", r.User_id).First(&user)
	if userExist.RowsAffected == 0 {
		return utils.ResObject(c, fiber.StatusBadGateway, "user not found", nil)
	}

	// if category is exist
	var category model.Category
	categoryExist := db.Where("id = ?", r.Category_id).First(&category)
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
	var r []model.CreateContent
	article_id := c.Params("article_id")

	if err := c.BodyParser(&r); err != nil {
		return utils.ResObject(c, fiber.StatusBadRequest, "cannot parse input", nil)
	}

	// valdiate body parser request
	// for _, req := range r {
	// 	if errs := myValidation.Validate(&req); len(errs) > 0 && errs[0].Error {
	// 		errsMsgs := make([]string, 0)

	// 		for _, err := range errs {
	// 			errsMsgs = append(errsMsgs, fmt.Sprintf(
	// 				"[%s]: '%v' | Need tobe '%s'",
	// 				err.FailedField,
	// 				err.Value,
	// 				err.Tag,
	// 			))
	// 		}

	// 		return utils.ResObject(c, fiber.StatusBadRequest, "validation error", errsMsgs)
	// 	}
	// }

	// if article is exist
	var article model.Article
	articleExist := db.Where("id = ?", article_id).First(&article)
	if articleExist.RowsAffected == 0 {
		return utils.ResObject(c, fiber.StatusBadGateway, "article not found", nil)
	}

	// insert batch content article
	var newBodys []model.Body
	for _, req := range r {
		newBody := model.Body{
			Article_id: article.ID,
			Content:    req.Content,
		}
		newBodys = append(newBodys, newBody)
	}

	// save in db
	if err := db.Create(&newBodys); err.Error != nil {
		return utils.ResObject(c, fiber.StatusInternalServerError, "cannot create content", nil)
	}

	return utils.ResObject(c, fiber.StatusCreated, "success save content", newBodys)
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
	if article.IsPost {
		article.IsPost = false
		errUpdate := db.Save(&article)
		if errUpdate.Error != nil {
			return utils.ResObject(c, fiber.StatusBadRequest, "cannot update article", nil)
		}

		return utils.ResObject(c, fiber.StatusOK, "success unpost article", article)
	} else {
		article.IsPost = true
		errUpdate := db.Save(&article)
		if errUpdate.Error != nil {
			return utils.ResObject(c, fiber.StatusBadRequest, "cannot update article", nil)
		}

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

func GetContentForEdit(c *fiber.Ctx) error {
	article_id := c.Params("article_id")
	user_id := c.Query("user_id")

	if article_id == "" || user_id == "" {
		return utils.ResObject(c, fiber.StatusBadRequest, "article id cannot be empty", nil)
	}

	// get  article
	var articleDetail model.Article
	errGet := db.Where("user_id = ?", user_id).Joins("User").Joins("Category").First(&articleDetail, article_id)
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

	return utils.ResObject(c, fiber.StatusOK, "success get article detail by user", response)
}

func DeleteContent(c *fiber.Ctx) error {
	article_id := c.Params("article_id")
	body_id := c.Query("body_id")

	if article_id == "" && body_id == "" {
		return utils.ResObject(c, fiber.StatusBadRequest, "article id and body id cannot be empty", nil)
	}

	// if article is exist
	var article model.Article
	articleExist := db.Where("id = ?", article_id).First(&article)
	if articleExist.RowsAffected == 0 {
		return utils.ResObject(c, fiber.StatusBadGateway, "article not found", nil)
	}

	// if body is exist
	var body model.Body
	bodyExist := db.Where("id = ?", body_id).First(&body)
	if bodyExist.RowsAffected == 0 {
		return utils.ResObject(c, fiber.StatusBadGateway, "body not found", nil)
	}

	// delete body
	errDelete := db.Delete(&body)
	if errDelete.Error != nil {
		return utils.ResObject(c, fiber.StatusBadRequest, "cannot delete body", nil)
	}

	return utils.ResObject(c, fiber.StatusOK, "success delete body", nil)
}

func DeleteFullArticle(c *fiber.Ctx) error {
	article_id := c.Params("article_id")

	if article_id == "" {
		return utils.ResObject(c, fiber.StatusBadRequest, "article id cannot be empty", nil)
	}

	// if article is exist
	var article model.Article
	articleExist := db.Where("id = ?", article_id).First(&article)
	if articleExist.RowsAffected == 0 {
		return utils.ResObject(c, fiber.StatusBadGateway, "article not found", nil)
	}

	// delete all body
	var body []model.Body
	errGetContent := db.Where("article_id = ?", article_id).Find(&body)
	if errGetContent.Error != nil {
		return utils.ResObject(c, fiber.StatusBadRequest, "cannot get article detail", nil)
	}

	errDeleteBody := db.Delete(&body)
	if errDeleteBody.Error != nil {
		return utils.ResObject(c, fiber.StatusBadRequest, "cannot delete body", nil)
	}

	// delete article
	errDelete := db.Delete(&article)
	if errDelete.Error != nil {
		return utils.ResObject(c, fiber.StatusBadRequest, "cannot delete article", nil)
	}

	return utils.ResObject(c, fiber.StatusOK, "success delete article", nil)
}

func EditProject(c *fiber.Ctx) error {
	article_id := c.Params("article_id")
	r := new(model.EditProject)

	if err := c.BodyParser(r); err != nil {
		return utils.ResObject(c, fiber.StatusBadRequest, "cannot parse input", nil)
	}

	// article is exist
	var article model.Article
	articleExist := db.Where("id = ?", article_id).First(&article)
	if articleExist.RowsAffected == 0 {
		return utils.ResObject(c, fiber.StatusBadGateway, "article not found", nil)
	}

	// update article
	article.Title = r.Title
	article.Introduction = r.Introduction

	errUpdate := db.Save(&article)
	if errUpdate.Error != nil {
		return utils.ResObject(c, fiber.StatusBadRequest, "cannot update article", nil)
	}

	return utils.ResObject(c, fiber.StatusOK, "success update article", article)
}

func EditThumbnail(c *fiber.Ctx) error {
	article_id := c.Params("article_id")
	file, err := c.FormFile("thumbnail")
	if err != nil {
		return utils.ResObject(c, fiber.StatusBadRequest, "cannot get file", nil)
	}

	// article exist
	var article model.Article
	articleExist := db.Where("id = ?", article_id).First(&article)
	if articleExist.RowsAffected == 0 {
		return utils.ResObject(c, fiber.StatusBadGateway, "article not found", nil)
	}

	// setting file name
	uniqueId := uuid.New()
	fileName := strings.Replace(uniqueId.String(), "-", "", -1)
	fileExt := strings.Split(file.Filename, ".")[1]
	image := fmt.Sprintf("%s.%s", fileName, fileExt)

	// save to db
	path := utils.GetEnv("FILE_PATH")
	thumbnail := path + "thumbnail/" + image
	article.Thumbnail = thumbnail
	errSave := db.Save(&article)
	if errSave.Error != nil {
		return utils.ResObject(c, fiber.StatusBadRequest, "cannot save thumbnail", nil)
	}

	// save to local
	if errSaveLoc := c.SaveFile(file, fmt.Sprintf("./public/thumbnail/%s", image)); errSaveLoc != nil {
		return utils.ResObject(c, fiber.StatusBadRequest, "cannot save thumbnail", nil)
	}

	return utils.ResObject(c, fiber.StatusOK, "success save thumbnail", article)
}
