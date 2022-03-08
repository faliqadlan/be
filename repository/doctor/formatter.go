package doctor

type ProfileResp struct {
	UserName string `json:"userName"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
	Image    string `json:"image"`
	Address  string `json:"address"`
	Status   string `json:"status"`
	OpenDay  string `json:"openDay"`
	CloseDay string `json:"closeDay"`
	Capacity int    `json:"capacity"`
}

// type
