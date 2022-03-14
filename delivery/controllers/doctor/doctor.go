package doctor

import (
	"be/api/aws/s3"
	"be/delivery/controllers/templates"
	"be/delivery/middlewares"
	"be/repository/doctor"
	"errors"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

type Controller struct {
	r      doctor.Doctor
	taskS3 s3.TaskS3M
}

func New(r doctor.Doctor, taskS3 s3.TaskS3M) *Controller {
	return &Controller{
		r:      r,
		taskS3: taskS3,
	}
}

func (cont *Controller) Create() echo.HandlerFunc {
	return func(c echo.Context) error {
		var req Req

		// request

		if err := c.Bind(&req); err != nil {
			log.Warn(err)
			return c.JSON(http.StatusBadRequest, templates.BadRequest(nil, "invalid input ", nil))
		}

		var v = validator.New()
		if err := v.Struct(req); err != nil {
			log.Warn(err)
			return c.JSON(http.StatusBadRequest, templates.BadRequest(nil, "invalid input ", nil))
		}

		// aws s3

		file, err := c.FormFile("file")
		if err != nil {
			log.Warn(err)
		}
		if err == nil {
			link, err := cont.taskS3.UploadFileToS3(*file)
			if err != nil {
				log.Warn(err)
				return c.JSON(http.StatusInternalServerError, templates.InternalServerError(nil, errors.New("there's some problem is server"), nil))
			}

			req.Image = link
		}

		// database

		res, err := cont.r.Create(*req.ToDoctor())

		if err != nil {
			log.Warn(err)
			switch err.Error() {
			case errors.New("can't assign capacity below zero").Error():
				err = errors.New("can't assign capacity below zero")
			case errors.New("user name is already exist").Error():
				err = errors.New("user name is already exist")
			default:
				err = errors.New("there's some problem is server")
			}
			return c.JSON(http.StatusInternalServerError, templates.InternalServerError(nil, err.Error(), nil))
		}

		return c.JSON(http.StatusCreated, templates.Success(http.StatusCreated, "success add Doctor", map[string]interface{}{
			"name": res.Name,
		}))
	}
}

func (cont *Controller) Update() echo.HandlerFunc {
	return func(c echo.Context) error {
		var uid = middlewares.ExtractTokenUid(c)
		var req Req

		// request

		if err := c.Bind(&req); err != nil {
			log.Warn(err)
			return c.JSON(http.StatusBadRequest, templates.BadRequest(nil, "invalid input", nil))
		}

		// aws s3

		res1, err := cont.r.GetProfile(uid)
		if err != nil {
			log.Warn(err)
		}

		file, err := c.FormFile("file")
		if err != nil {
			log.Warn(err)
		}
		if err == nil {
			if res1.Image != "https://www.teralogistics.com/wp-content/uploads/2020/12/default.png" {
				var nameFile = res1.Image

				nameFile = strings.Replace(nameFile, "https://karen-givi-bucket.s3.ap-southeast-1.amazonaws.com/", "", -1)

				var res = cont.taskS3.UpdateFileS3(nameFile, *file)
				log.Info(res)
				if res != "success" {
					log.Warn(res)
					return c.JSON(http.StatusInternalServerError, templates.InternalServerError(nil, errors.New("there's some problem is server"), nil))
				}
			} else {
				var link, err = cont.taskS3.UploadFileToS3(*file)
				if err != nil {
					log.Warn(err)
					return c.JSON(http.StatusInternalServerError, templates.InternalServerError(nil, errors.New("there's some problem is server"), nil))
				}

				req.Image = link
			}
		}

		// database

		res, err := cont.r.Update(uid, *req.ToDoctor())

		if err != nil {
			log.Warn(err)
			switch err.Error() {
			case errors.New("user name is already exist").Error():
				err = errors.New("user name is already exist")
			case errors.New("can't update capacity below total pending patients").Error():
				err = errors.New("can't update capacity below total pending patients")
			default:
				err = errors.New("there's some problem is server")
			}
			return c.JSON(http.StatusInternalServerError, templates.InternalServerError(nil, err.Error(), nil))
		}

		return c.JSON(http.StatusAccepted, templates.Success(http.StatusAccepted, "success update Doctor", res.Name))
	}
}

func (cont *Controller) Delete() echo.HandlerFunc {
	return func(c echo.Context) error {
		var uid = middlewares.ExtractTokenUid(c)

		// aws s3

		res1, err := cont.r.GetProfile(uid)
		if err != nil {
			log.Error(err)
		}

		if res1.Image != "https://www.teralogistics.com/wp-content/uploads/2020/12/default.png" {

			var nameFile = res1.Image

			nameFile = strings.Replace(nameFile, "https://karen-givi-bucket.s3.ap-southeast-1.amazonaws.com/", "", -1)
			res := cont.taskS3.DeleteFileS3(nameFile)
			log.Info(res)
			if res != "success" {
				log.Warn(res)
			}
		}

		// database

		res, err := cont.r.Delete(uid)

		if err != nil {
			log.Warn(err)
			return c.JSON(http.StatusInternalServerError, templates.InternalServerError(nil, "there's problem in server", nil))
		}

		return c.JSON(http.StatusAccepted, templates.Success(http.StatusAccepted, "success delete Doctor", res.DeletedAt))
	}
}

func (cont *Controller) GetProfile() echo.HandlerFunc {
	return func(c echo.Context) error {
		var uid = middlewares.ExtractTokenUid(c)

		// database

		var res, err = cont.r.GetProfile(uid)

		if err != nil {
			log.Warn(err)
			return c.JSON(http.StatusInternalServerError, templates.InternalServerError(nil, "there's problem in server", nil))
		}

		return c.JSON(http.StatusOK, templates.Success(nil, "success get profile Doctor", res))
	}
}

func (cont *Controller) GetAll() echo.HandlerFunc {
	return func(c echo.Context) error {

		// database

		var res, err = cont.r.GetAll()

		if err != nil {
			log.Warn(err)
			return c.JSON(http.StatusInternalServerError, templates.InternalServerError(nil, "there's problem in server", nil))
		}

		return c.JSON(http.StatusOK, templates.Success(nil, "success get all doctor's patient", res))
	}
}
