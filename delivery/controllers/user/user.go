package user

import (
	"be/delivery/controllers/templates"
	"be/entities"
	"be/repository/user"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type UserController struct {
	repo user.User
}

func New(repo user.User) *UserController {
	return &UserController{
		repo: repo,
	}
}

func (uc *UserController) Create() echo.HandlerFunc {
	return func(c echo.Context) error {
		user := UserCreateRequest{}

		if err := c.Bind(&user); err != nil {
			return c.JSON(http.StatusBadRequest, templates.BadRequest(nil, "There is some problem from input", err))
		}
		v := validator.New()
		if err := v.Struct(user); err != nil {
			return c.JSON(http.StatusBadRequest, templates.BadRequest(nil, "There is some problem from input", nil))
		}

		res, err := uc.repo.Create(entities.User{Name: user.Name, Email: user.Email, Password: user.Password})

		if err != nil {
			return c.JSON(http.StatusInternalServerError, templates.InternalServerError(nil, "error internal server error fo create new user "+err.Error(), err))
		}

		return c.JSON(http.StatusCreated, templates.Success(http.StatusCreated, "Success create new user", res.Name))
	}
}
