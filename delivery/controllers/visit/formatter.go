package visit

import (
	"be/entities"
	"time"

	"gorm.io/datatypes"
)

type Req struct {
	Doctor_uid       string `json:"doctor_uid" form:"doctor_uid" validate:"required"`
	Patient_uid      string `json:"patient_uid" form:"patient_uid"`
	Date             string `json:"date" form:"date" validate:"required"`
	Status           string `json:"status" form:"status"`
	Complaint        string `json:"complaint" form:"complaint"`
	MainDiagnose     string `json:"mainDiagnose" form:"mainDiagnose"`
	AdditionDiagnose string `json:"additionDiagnose" form:"additionDiagnose"`
	Action           string `json:"action" form:"action"`
	Recipe           string `json:"recipe" form:"recipe"`
	BloodPressure    string `json:"bloodPressuse" form:"bloodPressuse"`
	HeartRate        string `json:"heartRate" form:"heartRate"`
	O2Saturate       string `json:"o2Saturate" form:"o2Saturate"`
	Weight           int    `json:"weight" form:"weight"`
	Height           int    `json:"height" form:"height"`
	Bmi              int    `json:"bmi" form:"bmi"`
}

func (r *Req) ToVisit() *entities.Visit {
	var layout = "02-01-2006"

	var dateConv, _ = time.Parse(layout, r.Date)

	return &entities.Visit{
		Date:             datatypes.Date(dateConv),
		Status:           r.Status,
		Complaint:        r.Complaint,
		MainDiagnose:     r.MainDiagnose,
		AdditionDiagnose: r.AdditionDiagnose,
		Action:           r.Action,
		Recipe:           r.Recipe,
		BloodPressure:    r.BloodPressure,
		HeartRate:        r.HeartRate,
		O2Saturate:       r.O2Saturate,
		Weight:           r.Weight,
		Height:           r.Height,
		Bmi:              r.Bmi,
	}
}

type ResponseFormat struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}
