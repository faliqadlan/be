package visit

type VisitResp struct {
	Visit_uid string `json:"visit_uid"`
	Name      string `json:"name"`
	Nik       string `json:"nik"`
	Gender    string `json:"gender"`
	Date      string `json:"date"`
	Recipe    string `json:"recipe"`
	Diagnose  string `json:"diagnose"`
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
