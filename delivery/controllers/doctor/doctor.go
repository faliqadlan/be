package doctor

import (
	"be/api/aws/s3"
	"be/delivery/controllers/templates"
	"be/delivery/middlewares"
	"be/repository/doctor"
	"be/utils"
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
			switch {
			case strings.Contains(err.Error(), "UserName"):
				err = errors.New("invalid userName")
			case strings.Contains(err.Error(), "Email"):
				err = errors.New("invalid email")
			case strings.Contains(err.Error(), "Password"):
				err = errors.New("invalid password")
			case strings.Contains(err.Error(), "Name"):
				err = errors.New("invalid name")
			case strings.Contains(err.Error(), "Address"):
				err = errors.New("invalid address")
			case strings.Contains(err.Error(), "Status"):
				err = errors.New("invalid status")
			case strings.Contains(err.Error(), "OpenDay"):
				err = errors.New("invalid open day")
			case strings.Contains(err.Error(), "CloseDay"):
				err = errors.New("invalid close day")
			case strings.Contains(err.Error(), "Capacity"):
				err = errors.New("invalid capacity ")
			default:
				err = errors.New("invalid input")
			}
			return c.JSON(http.StatusBadRequest, templates.BadRequest(nil, err.Error(), nil))
		}

		// check capacity

		if req.Capacity < 0 {
			return c.JSON(http.StatusBadRequest, templates.BadRequest(nil, "can't assign capacity below zero", nil))
		}

		// check format user name

		if err := utils.UserNameValid(req.UserName); err != nil {
			log.Warn(err)
			return c.JSON(http.StatusBadRequest, templates.BadRequest(nil, err.Error(), nil))
		}

		// check format name

		if err := utils.NameValid(req.Name); err != nil {
			log.Warn(err)
			return c.JSON(http.StatusBadRequest, templates.BadRequest(nil, err.Error(), nil))
		}

		// check format address

		if err := utils.AddressValid(req.Address); err != nil {
			log.Warn(err)
			return c.JSON(http.StatusBadRequest, templates.BadRequest(nil, err.Error(), nil))
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
			switch {
			case err.Error() == errors.New("user name is already exist").Error():
				err = errors.New("user name is already exist")
			case err.Error() == errors.New("email is already exist").Error():
				err = errors.New("email is already exist")
			case strings.Contains(err.Error(), "status"):
				err = errors.New("invalid status")
			case strings.Contains(err.Error(), "open_day"):
				err = errors.New("invalid open day")
			case strings.Contains(err.Error(), "close_day"):
				err = errors.New("invalid close day")
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

		// check capacity

		if req.Capacity < 0 && req.Capacity != 0 {
			return c.JSON(http.StatusBadRequest, templates.BadRequest(nil, "can't assign capacity below zero", nil))
		}

		// check format user name

		if err := utils.UserNameValid(req.UserName); err != nil && req.UserName != "" {
			log.Warn(err)
			return c.JSON(http.StatusBadRequest, templates.BadRequest(nil, err.Error(), nil))
		}

		// check format name

		if err := utils.NameValid(req.Name); err != nil && req.Name != "" {
			log.Warn(err)
			return c.JSON(http.StatusBadRequest, templates.BadRequest(nil, err.Error(), nil))
		}

		// check format address

		if err := utils.AddressValid(req.Address); err != nil && req.Address != "" {
			log.Warn(err)
			return c.JSON(http.StatusBadRequest, templates.BadRequest(nil, err.Error(), nil))
		}

		// aws s3

		file, err := c.FormFile("file")
		if err != nil {
			log.Warn(err)
		}
		if err == nil {
			res1, err := cont.r.GetProfile(uid)
			if err != nil {
				log.Warn(err)
				return c.JSON(http.StatusInternalServerError, templates.InternalServerError(nil, errors.New("there's some problem is server"), nil))
			}
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
			switch {
			case err.Error() == errors.New("user name is already exist").Error():
				err = errors.New("user name is already exist")
			case err.Error() == errors.New("can't update capacity below total pending patients").Error():
				err = errors.New("can't update capacity below total pending patients")
			case strings.Contains(err.Error(), "status"):
				err = errors.New("invalid status")
			case strings.Contains(err.Error(), "open_day"):
				err = errors.New("invalid open day")
			case strings.Contains(err.Error(), "close_day"):
				err = errors.New("invalid close day")
			case err.Error() == "record not found":
				err = errors.New("account is not found")
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
			switch err.Error() {
			case "record not found":
				err = errors.New("account is not found")
			default:
				err = errors.New("there's some problem is server")
			}
			return c.JSON(http.StatusInternalServerError, templates.InternalServerError(nil, err.Error(), nil))
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
			switch err.Error() {
			case "record not found":
				err = errors.New("account is not found")
			default:
				err = errors.New("there's some problem is server")
			}
			return c.JSON(http.StatusInternalServerError, templates.InternalServerError(nil, err.Error(), nil))
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
