package patient

import (
	"be/api/aws/s3"
	"be/delivery/controllers/templates"
	logic "be/delivery/logic/patient"
	"be/delivery/middlewares"
	"be/repository/patient"
	"errors"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

type Controller struct {
	r      patient.Patient
	taskS3 s3.TaskS3M
	l      logic.Patient
}

func New(r patient.Patient, taskS3 s3.TaskS3M, l logic.Patient) *Controller {
	return &Controller{
		r:      r,
		taskS3: taskS3,
		l:      l,
	}
}

func (cont *Controller) Create() echo.HandlerFunc {
	return func(c echo.Context) error {
		var req logic.Req

		// request

		if err := c.Bind(&req); err != nil {
			log.Warn(err)
			return c.JSON(http.StatusBadRequest, templates.BadRequest(nil, "invalid input", nil))
		}

		// validation struct

		switch {
		case req.UserName == "":
			return c.JSON(http.StatusBadRequest, templates.BadRequest(nil, "invalid user name ", nil))
		case req.Email == "":
			return c.JSON(http.StatusBadRequest, templates.BadRequest(nil, "invalid email ", nil))
		case req.Password == "":
			return c.JSON(http.StatusBadRequest, templates.BadRequest(nil, "invalid password ", nil))
		}

		// validaion request

		if err := cont.l.ValidationRequest(req); err != nil {
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
				return c.JSON(http.StatusInternalServerError, templates.InternalServerError(nil, "there's some problem is server", nil))
			}

			req.Image = link
		}

		// database

		entity, err := req.ToPatient()
		if err != nil {
			log.Warn(err)
			return c.JSON(http.StatusBadRequest, templates.BadRequest(nil, err.Error(), nil))
		}

		res, err := cont.r.Create(*entity)

		if err != nil {
			log.Warn(err)
			switch err.Error() {
			case errors.New("user name is already exist").Error():
				err = errors.New("user name is already exist")
			case errors.New("email is already exist").Error():
				err = errors.New("email is already exist")
			default:
				err = errors.New("there's some problem is server")
			}

			return c.JSON(http.StatusInternalServerError, templates.InternalServerError(nil, err.Error(), nil))
		}

		// generate token

		token, err := middlewares.GenerateToken(res.Patient_uid, "patient")

		if err != nil {
			log.Warn(err)
			err = errors.New("there's some problem is server")
			return c.JSON(http.StatusNotAcceptable, templates.BadRequest(http.StatusNotAcceptable, err.Error(), nil))
		}

		return c.JSON(http.StatusCreated, templates.Success(http.StatusCreated, "success add patient", map[string]interface{}{
			"token":    token,
			"userName": res.UserName,
		}))
	}
}

func (cont *Controller) Update() echo.HandlerFunc {
	return func(c echo.Context) error {
		var uid, _ = middlewares.ExtractTokenUid(c)
		var req logic.Req

		// request

		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, templates.BadRequest(nil, err.Error(), nil))
		}

		var v = validator.New()
		if err := v.Struct(req); err != nil {
			log.Warn(err)
			switch {
			case strings.Contains(err.Error(), "Nik"):
				err = errors.New("invalid nik")
			case strings.Contains(err.Error(), "Name"):
				err = errors.New("invalid name")
			case strings.Contains(err.Error(), "Gender"):
				err = errors.New("invalid gender")
			case strings.Contains(err.Error(), "Address"):
				err = errors.New("invalid address")
			case strings.Contains(err.Error(), "PlaceBirth"):
				err = errors.New("invalid place birth")
			case strings.Contains(err.Error(), "Dob"):
				err = errors.New("invalid date of birth")
			case strings.Contains(err.Error(), "Job"):
				err = errors.New("invalid job")
			case strings.Contains(err.Error(), "Status"):
				err = errors.New("invalid status")
			case strings.Contains(err.Error(), "Religion"):
				err = errors.New("invalid religion")
			default:
				err = errors.New("invalid input")
			}
			return c.JSON(http.StatusBadRequest, templates.BadRequest(nil, err.Error(), nil))
		}

		// validaion request

		if err := cont.l.ValidationRequest(req); err != nil {
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

		entity, err := req.ToPatient()
		if err != nil {
			log.Warn(err)
			return c.JSON(http.StatusBadRequest, templates.BadRequest(nil, err.Error(), nil))
		}

		res, err := cont.r.Update(uid, *entity)

		if err != nil {
			log.Warn(err)
			switch err.Error() {
			case errors.New("user name is already exist").Error():
				err = errors.New("user name is already exist")
			case errors.New("email is already exist").Error():
				err = errors.New("email is already exist")
			case "record not found":
				err = errors.New("account is not found")
			default:
				err = errors.New("there's some problem is server")
			}
			return c.JSON(http.StatusInternalServerError, templates.InternalServerError(nil, err.Error(), nil))
		}

		return c.JSON(http.StatusAccepted, templates.Success(http.StatusAccepted, "success update patient", res.Name))
	}
}

func (cont *Controller) Delete() echo.HandlerFunc {
	return func(c echo.Context) error {
		var uid, _ = middlewares.ExtractTokenUid(c)

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

		return c.JSON(http.StatusAccepted, templates.Success(http.StatusAccepted, "success delete patient", res.DeletedAt))
	}
}

func (cont *Controller) GetProfile() echo.HandlerFunc {
	return func(c echo.Context) error {
		var uid string

		if uid = c.QueryParam("patient_uid"); uid == "" {
			uid, _ = middlewares.ExtractTokenUid(c)
		}

		// database
		// log.Info(uid)
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

		return c.JSON(http.StatusOK, templates.Success(http.StatusOK, "success get profile patient", res))
	}
}
