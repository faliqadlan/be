package doctor

type ProfileResp struct {
	Doctor_uid     string `json:"doctor_uid"`
	Doctor_uid_ref string `json:"doctor_uid_ref"`
	UserName       string `json:"userName"`
	Email          string `json:"email"`
	Name           string `json:"name"`
	Image          string `json:"image"`
	Address        string `json:"address"`
	Status         string `json:"status"`
	OpenDay        string `json:"openDay"`
	CloseDay       string `json:"closeDay"`
	Capacity       int    `json:"capacity"`
}

type AllResp struct {
	Doctor_uid string `json:"doctor_uid"`
	Name       string `json:"name"`
	Image      string `json:"image"`
	Address    string `json:"address"`
	Status     string `json:"status"`
	Capacity   int    `json:"capacity"`
}

type All struct {
	Doctors []AllResp `JSON:"doctors"`
}
