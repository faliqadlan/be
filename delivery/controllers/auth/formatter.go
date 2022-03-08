package auth

type Userlogin struct {
	UserName string `json:"userName"  form:"userName"`
	Password string `json:"password" form:"password"`
}

type LoginRespFormat struct {
	Code    int                    `json:"code"`
	Message string                 `json:"message"`
	Data    map[string]interface{} `json:"data"`
}
