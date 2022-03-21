package visit

import (
	"be/entities"
	"errors"
	"time"

	"gorm.io/datatypes"
)

type Req struct {
	Event_uid        string
	Doctor_uid       string `json:"doctor_uid" form:"doctor_uid" validate:"required"`
	Patient_uid      string `json:"patient_uid" form:"patient_uid"`
	Date             string `json:"date" form:"date" validate:"required"`
	Status           string `json:"status" form:"status"`
	Complaint        string `json:"complaint" form:"complaint"  validate:"required"`
	MainDiagnose     string `json:"mainDiagnose" form:"mainDiagnose"`
	AdditionDiagnose string `json:"additionDiagnose" form:"additionDiagnose"`
	Action           string `json:"action" form:"action"`
	Recipe           string `json:"recipe" form:"recipe"`
	BloodPressure    string `json:"bloodPressuse" form:"bloodPressuse"`
	HeartRate        string `json:"heartRate" form:"heartRate"`
	RespiratoryRate  string `json:"respiratoryRate" form:"respiratoryRate"`
	O2Saturate       string `json:"o2Saturate" form:"o2Saturate"`
	Weight           string `json:"weight" form:"weight"`
	Height           string `json:"height" form:"height"`
	Bmi              string `json:"bmi" form:"bmi"`
}

func (r *Req) ToVisit() (*entities.Visit, error) {
	var layout = "02-01-2006"

	var dateConv, err = time.Parse(layout, r.Date)
	if err != nil && r.Date != "" {
		return &entities.Visit{}, errors.New("invalid date format")
	}
	// log.Info(time.Since(dateConv), r.Date)
	// if time.Since(dateConv) > 0 && r.Date != "" && time.Now() {
	// 	return &entities.Visit{}, errors.New("invalid date is in the past")
	// }
	return &entities.Visit{
		Event_uid:        r.Event_uid,
		Date:             datatypes.Date(dateConv),
		Status:           r.Status,
		Complaint:        r.Complaint,
		MainDiagnose:     r.MainDiagnose,
		AdditionDiagnose: r.AdditionDiagnose,
		Action:           r.Action,
		Recipe:           r.Recipe,
		BloodPressure:    r.BloodPressure,
		HeartRate:        r.HeartRate,
		RespiratoryRate:  r.RespiratoryRate,
		O2Saturate:       r.O2Saturate,
		Weight:           r.Weight,
		Height:           r.Height,
		Bmi:              r.Bmi,
	}, nil
}
