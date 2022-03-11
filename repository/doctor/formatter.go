package doctor

type ProfileResp struct {
	Doctor_uid   string `json:"doctor_uid"`
	UserName     string `json:"userName"`
	Email        string `json:"email"`
	Name         string `json:"name"`
	Image        string `json:"image"`
	Address      string `json:"address"`
	Status       string `json:"status"`
	OpenDay      string `json:"openDay"`
	CloseDay     string `json:"closeDay"`
	Capacity     int    `json:"capacity"`
	LeftCapacity int    `json:"leftCapacity"`
}

type PatientResp struct {
	Patient_uid string `json:"patient_uid"`
	Name        string `json:"name"`
	Gender      string `json:"gender"`
	Nik         string `json:"nik"`
	TotalVisit  int    `json:"totalVisit"`
}

type PatientsResp struct {
	Patients []PatientResp `json:"patientResp"`
}

type Visit struct {
	Patient_uid string `json:"patient_uid"`
	Name        string `json:"name"`
	Gender      string `json:"gender"`
	Nik         string `json:"nik"`
	Status      string `json:"status"`
}

type Dashboard struct {
	TotalPatient     int `json:"totalPatient"`
	TotalVisitDay    int `json:"totalVisitDay"`
	TotalAppointment int `json:"totalAppointment"`
	// Visits           []Visit `json:"visits"`
}

type AllResp struct {
	Doctor_uid string `json:"doctor_uid"`
	Name       string `json:"name"`
	Image      string `json:"image"`
	Address    string `json:"address"`
	Status     string `json:"status"`
}

type All struct {
	Doctors []AllResp `JSON:"doctors"`
}
