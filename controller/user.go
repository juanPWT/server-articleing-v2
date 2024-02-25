package controller

import (
	"fmt"
	"server-article/config"
	"server-article/model"
	s "server-article/service"
	"server-article/utils"

	goValidator "github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/thanhpk/randstr"
	"gorm.io/gorm"
)

var db *gorm.DB = config.ConnectDB()
var validate = goValidator.New()

var myValidation = &s.XValidator{
	Validator: validate,
}

func SignUp(c *fiber.Ctx) error {
	//body req
	u := new(model.UserSignUp)

	if err := c.BodyParser(u); err != nil {
		return utils.ResObject(c, fiber.StatusBadRequest, "cannot parse request body", nil)
	}

	// validate body parser request
	if errs := myValidation.Validate(u); len(errs) > 0 && errs[0].Error {
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

	// is user exits
	userExist := db.Where("email = ?", u.Email).First(&model.User{})
	if userExist.RowsAffected > 0 {
		return utils.ResObject(c, fiber.StatusBadRequest, "email already exist", nil)
	}

	// create user
	u.Password = utils.HashPass([]byte(u.Password))

	request := &model.User{
		Username: u.Username,
		Email:    u.Email,
		Password: u.Password,
	}

	errCreate := db.Create(&request)
	if errCreate.Error != nil {
		return utils.ResObject(c, fiber.StatusBadRequest, "cannot create user", nil)
	}

	//generate verification code
	code := randstr.String(20)

	verification_code := utils.Encode(code)

	// update user verification code
	request.Verification_code = verification_code
	db.Save(&request)

	// ? send email
	clientOrigin := utils.GetEnv("CLIENT_ORIGIN")
	var firstname = u.Username
	emailData := s.EmailData{
		URL:       clientOrigin + "/verifyemail/" + code,
		FirstName: firstname,
		Subject:   "your account verification code",
	}

	s.SendEmail(request, &emailData)

	message := "success sign up, please check your email to verify your account"
	return utils.ResObject(c, fiber.StatusCreated, message, u)
}

func SignIn(c *fiber.Ctx) error {
	//body req
	u := new(model.UserSignIn)

	if err := c.BodyParser(u); err != nil {
		return utils.ResObject(c, fiber.StatusBadRequest, "cannot parse request body", nil)
	}

	//validate req
	if errs := myValidation.Validate(u); len(errs) > 0 && errs[0].Error {
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

	// find user
	var user model.User
	err := db.Where("email = ?", u.Email).First(&user)
	if err.RowsAffected <= 0 {
		return utils.ResObject(c, fiber.StatusBadRequest, "email account not exist", nil)
	}

	// compare pass
	isValidPass := utils.VerifyPass([]byte(user.Password), []byte(u.Password))
	if !isValidPass {
		return utils.ResObject(c, fiber.StatusBadRequest, "invalid password", nil)
	}

	// geberate jwt token
	claims := jwt.MapClaims{
		"iss":        "user",
		"uid":        user.ID,
		"email":      user.Email,
		"image":      user.Image,
		"username":   user.Username,
		"isVeriveid": user.Verified_email,
		"created_at": user.Created_at,
	}
	token, error := s.CreateToken(claims)
	if error != nil {
		return utils.ResObject(c, fiber.StatusBadRequest, "sorry cannot create access token", nil)
	}

	return utils.ResObject(c, fiber.StatusOK, "success sign in", fiber.Map{
		"access_token": token,
	})
}
