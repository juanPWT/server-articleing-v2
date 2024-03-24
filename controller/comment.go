package controller

import (
	"fmt"
	"server-article/model"
	"server-article/utils"

	"github.com/gofiber/fiber/v2"
)

func GetCommentByArticle(c *fiber.Ctx) error {
	article_id := c.Params("article_id")

	if article_id == "" {
		return utils.ResObject(c, fiber.StatusBadRequest, "Invalid request", nil)
	}

	var comments []model.Comment
	errComment := db.Where("article_id = ? ", article_id).Joins("User").Find(&comments)
	if errComment.RowsAffected == 0 {
		return utils.ResObject(c, fiber.StatusNotFound, "No comment found", nil)
	}

	return utils.ResObject(c, fiber.StatusOK, "Success get comment", comments)
}

func GetReplyByComment(c *fiber.Ctx) error {
	comment_id := c.Params("comment_id")

	if comment_id == "" {
		return utils.ResObject(c, fiber.StatusBadRequest, "Invalid request", nil)
	}

	var replies []model.Reply
	errReplies := db.Where("comment_id = ?", comment_id).Joins("User").Find(&replies)
	if errReplies.RowsAffected == 0 {
		return utils.ResObject(c, fiber.StatusNotFound, "No reply found", nil)
	}

	return utils.ResObject(c, fiber.StatusOK, "Success get reply", replies)
}

func CommentArticle(c *fiber.Ctx) error {
	r := new(model.CommentRequest)

	if err := c.BodyParser(r); err != nil {
		return utils.ResObject(c, fiber.StatusBadRequest, "Invalid request", nil)
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

	// save comment
	var comment model.Comment
	comment.Article_id = r.Article_id
	comment.User_id = r.User_id
	comment.Content = r.Content

	if errSave := db.Save(&comment); errSave.Error != nil {
		return utils.ResObject(c, fiber.StatusInternalServerError, "Failed save comment", nil)
	}

	return utils.ResObject(c, fiber.StatusOK, "Success save comment", comment)
}

func ReplyComment(c *fiber.Ctx) error {
	r := new(model.ReplyRequest)

	if err := c.BodyParser(r); err != nil {
		return utils.ResObject(c, fiber.StatusBadRequest, "Invalid request", nil)
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

	// is comment exist
	var comment model.Comment
	errComment := db.First(&comment, r.Comment_id)
	if errComment.RowsAffected == 0 {
		return utils.ResObject(c, fiber.StatusNotFound, "Comment not found", nil)
	}

	// save reply
	var reply model.Reply
	reply.Comment_id = r.Comment_id
	reply.User_id = r.User_id
	reply.Content = r.Content

	errSave := db.Save(&reply)
	if errSave.Error != nil {
		return utils.ResObject(c, fiber.StatusInternalServerError, "Failed save reply", nil)
	}

	// update replies comment
	comment.Replies += 1
	db.Save(&comment)

	return utils.ResObject(c, fiber.StatusOK, "Success save reply", reply)
}
