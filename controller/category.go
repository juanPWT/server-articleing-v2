package controller

import (
	"fmt"
	"server-article/model"
	"server-article/utils"

	"github.com/gofiber/fiber/v2"
)

func CreateCategory(c *fiber.Ctx) error {
	r := new(model.CategoryRequest)

	if err := c.BodyParser(r); err != nil {
		return utils.ResObject(c, fiber.StatusBadRequest, "cannot parse input", nil)
	}

	// validate body parser request
	if errs := myValidation.Validate(r); len(errs) > 0 && errs[0].Error {
		errMsgs := make([]string, 0)

		for _, err := range errs {
			errMsgs = append(errMsgs, fmt.Sprintf(
				"[%s]: '%v' | Needs to be '%s'",
				err.FailedField,
				err.Value,
				err.Tag,
			))
		}

		return utils.ResObject(c, fiber.StatusBadRequest, "validation error", errMsgs)
	}

	// if category is exist
	categoryExist := db.Where("name = ?", r.Name).First(&model.Category{})
	if categoryExist.RowsAffected > 0 {
		return utils.ResObject(c, fiber.StatusBadRequest, "category already exist", nil)
	}

	// create category
	category := &model.Category{
		Name: r.Name,
	}
	errCreate := db.Create(&category)
	if errCreate.Error != nil {
		return utils.ResObject(c, fiber.StatusBadRequest, "cannot create category", nil)
	}

	return utils.ResObject(c, fiber.StatusCreated, "success", category)
}

func GetAllCategory(c *fiber.Ctx) error {
	var categories []model.Category

	err := db.Find(&categories)
	if err.Error != nil {
		return utils.ResObject(c, fiber.StatusBadRequest, "cannot get categories", nil)
	}

	return utils.ResObject(c, fiber.StatusOK, "success", categories)
}
