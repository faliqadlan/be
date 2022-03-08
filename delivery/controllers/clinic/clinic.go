package clinic

import (
	"be/delivery/controllers/templates"
	"be/delivery/middlewares"
	"be/repository/clinic"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type Controller struct {
	r clinic.Clinic
}

func New(r clinic.Clinic) *Controller {
	return &Controller{
		r: r,
	}
}

func (cont *Controller) Create() echo.HandlerFunc {
	return func(c echo.Context) error {
		var req Req

		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, templates.BadRequest(nil, "error bad request for add clinic "+err.Error(), nil))
		}

		var v = validator.New()
		if err := v.Struct(req); err != nil {
			return c.JSON(http.StatusBadRequest, templates.BadRequest(nil, "error bad request for add clinic "+err.(validator.ValidationErrors).Error(), nil))
		}

		var res, err = cont.r.Create(*req.ToClinic())

		if err != nil {
			c.JSON(http.StatusInternalServerError, templates.InternalServerError(nil, "error internal server for add clinic "+err.Error(), nil))
		}

		token, err := middlewares.GenerateToken(res.Clinic_uid)

		if err != nil {
			return c.JSON(http.StatusNotAcceptable, templates.BadRequest(http.StatusNotAcceptable, "error in process token "+err.Error(), err))
		}

		return c.JSON(http.StatusCreated, templates.Success(http.StatusCreated, "success add clinic", token))
	}
}
