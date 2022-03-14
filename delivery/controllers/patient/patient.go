package patient

import (
	utils "be/api/aws"
	"be/delivery/controllers/templates"
	"be/delivery/middlewares"
	"be/repository/patient"
	"errors"
	"net/http"
	"strings"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

type Controller struct {
	r    patient.Patient
	conf *session.Session
}

func New(r patient.Patient, awsS3 *session.Session) *Controller {
	return &Controller{
		r:    r,
		conf: awsS3,
	}
}

func (cont *Controller) Create() echo.HandlerFunc {
	return func(c echo.Context) error {
		var req Req

		// request

		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, templates.BadRequest(nil, "error binding for add patient "+err.Error(), nil))
		}

		var v = validator.New()
		if err := v.Struct(req); err != nil {
			return c.JSON(http.StatusBadRequest, templates.BadRequest(nil, "error validator for add patient "+err.(validator.ValidationErrors).Error(), nil))
		}

		// aws s3

		file, err := c.FormFile("file")
		if err != nil {
			log.Warn(err)
		}
		if err == nil {
			link, err := utils.UploadFileToS3(cont.conf, *file)
			if err != nil {
				log.Warn(err)
				return c.JSON(http.StatusInternalServerError, templates.InternalServerError(nil, errors.New("there's some problem is server"), nil))
			}

			req.Image = link
		}

		// database

		res, err := cont.r.Create(*req.ToPatient())

		if err != nil {
			// log.Info(err)
			return c.JSON(http.StatusInternalServerError, templates.InternalServerError(nil, "error internal server for add patient "+err.Error(), nil))
		}

		return c.JSON(http.StatusCreated, templates.Success(http.StatusCreated, "success add patient", res.Name))
	}
}

func (cont *Controller) Update() echo.HandlerFunc {
	return func(c echo.Context) error {
		var uid = middlewares.ExtractTokenUid(c)
		var req Req

		// request

		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, templates.BadRequest(nil, "error binding for update patient "+err.Error(), nil))
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

				var res = utils.UpdateFileS3(cont.conf, nameFile, *file)
				log.Info(res)
				if res != "success" {
					log.Warn(res)
					return c.JSON(http.StatusInternalServerError, templates.InternalServerError(nil, errors.New("there's some problem is server"), nil))
				}
			} else {
				var link, err = utils.UploadFileToS3(cont.conf, *file)
				if err != nil {
					log.Warn(err)
					return c.JSON(http.StatusInternalServerError, templates.InternalServerError(nil, errors.New("there's some problem is server"), nil))
				}

				req.Image = link
			}
		}

		// database

		res, err := cont.r.Update(uid, *req.ToPatient())

		if err != nil {
			// log.Info(err)
			return c.JSON(http.StatusInternalServerError, templates.InternalServerError(nil, "error internal server for update patient "+err.Error(), nil))
		}

		return c.JSON(http.StatusAccepted, templates.Success(http.StatusAccepted, "success update patient", res.Name))
	}
}

func (cont *Controller) Delete() echo.HandlerFunc {
	return func(c echo.Context) error {
		var uid = middlewares.ExtractTokenUid(c)

		// database

		var res, err = cont.r.Delete(uid)

		if err != nil {
			// log.Info(err)
			return c.JSON(http.StatusInternalServerError, templates.InternalServerError(nil, "error internal server for delete patient "+err.Error(), nil))
		}

		// aws s3

		res1, err := cont.r.GetProfile(uid)
		if err != nil {
			log.Error(err)
		}

		if res1.Image != "https://www.teralogistics.com/wp-content/uploads/2020/12/default.png" {

			var nameFile = res1.Image

			nameFile = strings.Replace(nameFile, "https://karen-givi-bucket.s3.ap-southeast-1.amazonaws.com/", "", -1)
			log.Info(cont.conf)
			res := utils.DeleteFileS3(cont.conf, nameFile)
			log.Info(res)
			if res != "success" {
				log.Warn(res)
			}
		}

		return c.JSON(http.StatusAccepted, templates.Success(http.StatusAccepted, "success delete patient", res.DeletedAt))
	}
}

func (cont *Controller) GetProfile() echo.HandlerFunc {
	return func(c echo.Context) error {
		var uid string

		if uid = c.QueryParam("patient_uid"); uid == "" {
			uid = middlewares.ExtractTokenUid(c)
		}

		// database
		// log.Info(uid)
		var res, err = cont.r.GetProfile(uid)

		if err != nil {
			// log.Info(err)
			return c.JSON(http.StatusInternalServerError, templates.InternalServerError(nil, "error internal server for get profile patient "+err.Error(), nil))
		}

		return c.JSON(http.StatusOK, templates.Success(http.StatusOK, "success get profile patient", res))
	}
}
