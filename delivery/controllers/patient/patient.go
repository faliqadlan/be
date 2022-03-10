package patient

import (
	"be/delivery/controllers/templates"
	"be/delivery/middlewares"
	"be/repository/patient"
	"be/utils"
	"net/http"

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

		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, templates.BadRequest(nil, "error binding for add patient "+err.Error(), nil))
		}

		var v = validator.New()
		if err := v.Struct(req); err != nil {
			return c.JSON(http.StatusBadRequest, templates.BadRequest(nil, "error validator for add patient "+err.(validator.ValidationErrors).Error(), nil))
		}
		var file, err1 = c.FormFile("file")
		if err1 != nil {
			log.Info(err1)
		}
		if err1 == nil {
			var link = utils.UploadFileToS3(cont.conf, *file)

			req.Image = link
		}
		var res, err = cont.r.Create(*req.ToPatient())

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

		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, templates.BadRequest(nil, "error binding for update patient "+err.Error(), nil))
		}

		var res, err = cont.r.Update(uid, *req.ToPatient())

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

		var res, err = cont.r.Delete(uid)

		if err != nil {
			// log.Info(err)
			return c.JSON(http.StatusInternalServerError, templates.InternalServerError(nil, "error internal server for delete patient "+err.Error(), nil))
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
		// log.Info(uid)
		var res, err = cont.r.GetProfile(uid)

		if err != nil {
			// log.Info(err)
			return c.JSON(http.StatusInternalServerError, templates.InternalServerError(nil, "error internal server for get profile patient "+err.Error(), nil))
		}

		return c.JSON(http.StatusOK, templates.Success(http.StatusOK, "success get profile patient", res))
	}
}

func (cont *Controller) GetRecords() echo.HandlerFunc {
	return func(c echo.Context) error {
		var uid string

		if uid = c.QueryParam("patient_uid"); uid == "" {
			uid = middlewares.ExtractTokenUid(c)
		}
		log.Info(uid)
		var res, err = cont.r.GetRecords(uid)

		if err != nil {
			// log.Info(err)
			return c.JSON(http.StatusInternalServerError, templates.InternalServerError(nil, "error internal server for get patient record "+err.Error(), nil))
		}

		return c.JSON(http.StatusOK, templates.Success(http.StatusOK, "success get patient record", res))
	}
}

func (cont *Controller) GetHistories() echo.HandlerFunc {
	return func(c echo.Context) error {

		var uid = middlewares.ExtractTokenUid(c)

		log.Info(uid)
		var res, err = cont.r.GetHistories(uid)

		if err != nil {
			// log.Info(err)
			return c.JSON(http.StatusInternalServerError, templates.InternalServerError(nil, "error internal server for get patient history "+err.Error(), nil))
		}

		return c.JSON(http.StatusOK, templates.Success(http.StatusOK, "success get patient history", res))
	}
}

func (cont *Controller) GetAppointMent() echo.HandlerFunc {
	return func(c echo.Context) error {

		var uid = middlewares.ExtractTokenUid(c)

		log.Info(uid)
		var res, err = cont.r.GetAppointMent(uid)

		if err != nil {
			// log.Info(err)
			return c.JSON(http.StatusInternalServerError, templates.InternalServerError(nil, "error internal server for get patient appoinment "+err.Error(), nil))
		}

		return c.JSON(http.StatusOK, templates.Success(http.StatusOK, "success get patient appoinment", res))
	}
}
