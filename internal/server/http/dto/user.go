package dto

const bearer = "Bearer"

type UserRegistration struct {
	UserName string `json:"username"`
	Password string `json:"password"`
	TimeZone string `json:"time_zone"`
}

type UserLogin struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}

type UserUpdate struct {
	TimeZone string `json:"time_zone,omitempty"`
}

type IdObject struct {
	Id int64 `json:"id"`
}

type Token struct {
	Token string `json:"token"`
	Type string `json:"type"`
}

type JwtToken = string
func NewTokenBearer(tk JwtToken) *Token {
	return &Token{
		Token:  tk,
		Type: bearer,
	}
}