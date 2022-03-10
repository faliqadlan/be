package patient

type Profile struct {
	Patient_uid string `json:"patient_uid"`
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

type Record struct {
	Visit_uid        string `json:"visit_uid"`
	MainDiagnose     string `json:"mainDiagnose"`
	AdditionDiagnose string `json:"additionDiagnose"`
	Action           string `json:"action"`
	Recipe           string `json:"recipe"`
	BloodPressuse    string `json:"bloodPressuse"`
	HeartRate        string `json:"heartRate"`
	O2Saturatin      string `json:"o2Saturatin"`
	Weight           string `json:"weight"`
	Height           string `json:"height"`
	Bmi              string `json:"bmi"`
}

type Records struct {
	Records []Record `json:"records"`
}

type History struct {
	Date             string `json:"date"`
	Name             string `json:"name"`
	Address          string `json:"addres"`
	MainDiagnose     string `json:"mainDiagnose"`
	AdditionDiagnose string `json:"additionDiagnose"`
	Recipe           string `json:"recipe"`
}

type Histories struct {
	Histories []History `json:"histories"`
}

type Apppoinment struct {
	Day     string `json:"day"`
	Date    string `json:"date"`
	Name    string `json:"name"`
	Address string `json:"addres"`
}
