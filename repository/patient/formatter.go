package patient

type Profile struct {
	Patient_uid string `json:"patient_uid"`
	UserName    string `json:"userName"`
	Email       string `json:"email"`
	Nik         string `json:"nik"`
	Name        string `json:"name"`
	Image       string `json:"image"`
	Gender      string `json:"gender"`
	Address     string `json:"address"`
	PlaceBirth  string `json:"placeBirth"`
	Dob         string `json:"dob"`
	Religion    string `json:"religion"`
	Status      string `json:"status"`
	Job         string `json:"job"`
}

type Apppoinment struct {
	Day     string `json:"day"`
	Date    string `json:"date"`
	Name    string `json:"name"`
	Address string `json:"addres"`
}

type PatientAll struct {
	Patient_uid string `json:"patient_uid"`
	Nik         string `json:"nik"`
	Name        string `json:"name"`
	Gender      string `json:"gender"`
}

type All struct {
	Patients []PatientAll `json:"patients"`
}
