package visit

type VisitResp struct {
	Visit_uid        string `json:"visit_uid"`
	Date             string `json:"date" form:"date" validate:"required"`
	Status           string `json:"status" form:"status"`
	Complaint        string `json:"complaint" form:"complaint"`
	MainDiagnose     string `json:"mainDiagnose" form:"mainDiagnose"`
	AdditionDiagnose string `json:"addiditonDiagnose" form:"addiditonDiagnose"`
	Action           string `json:"action" form:"action"`
	Recipe           string `json:"recipe" form:"recipe"`
	BloodPressure    string `json:"bloodPressuse" form:"bloodPressuse"`
	HeartRate        string `json:"heartRate" form:"heartRate"`
	O2Saturate       string `json:"o2Saturate" form:"o2Saturate"`
	Weight           string `json:"weight" form:"weight"`
	Height           string `json:"height" form:"height"`
	Bmi              string `json:"bmi" form:"bmi"`

	Doctor_uid    string `json:"doctor_uid"`
	DoctorName    string `json:"doctorName"`
	DoctorAddress string `json:"doctorAddress"`

	Patient_uid string `json:"patient_uid"`
	PatientName string `json:"patientName"`
	Gender      string `json:"gender"`
	Nik         string `json:"nik"`
}

type Visits struct {
	Visits []VisitResp `json:"visits"`
}

type VisitCalendar struct {
	Address      string `json:"address"`
	Complaint    string `json:"complaint"`
	Date         string `json:"date"`
	DoctorName   string `json:"doctorName"`
	PatientName  string `json:"patientName"`
	DoctorEmail  string `json:"doctorEmail"`
	PatientEmail string `json:"patientEmail"`
}
