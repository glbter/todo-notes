package model

type User struct {
	Id           Id
	Name         string
	PasswordHash []byte
	TimeZone     string
}

func NewUser(id Id, name string, hash []byte, timeZone string) *User {
	return &User{
		Id: id,
		Name: name,
		PasswordHash: hash,
		TimeZone: timeZone,
	}
}

type UserInReq struct {
	Id int64
	//TimeZone TimeZone
}

type UserNew struct {
	Name         string
	Password     string
	TimeZone     string
}

type UserUpdate struct {
	Id Id
	TimeZone string
}