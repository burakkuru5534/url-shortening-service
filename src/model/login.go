package model

type LoginData struct {
	UpassFromDb *string
	UserID      int64
	UserEmail   string
}

type LoginReqData struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
