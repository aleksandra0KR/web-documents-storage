package model

type User struct {
	Login    string `json:"login,omitempty" gorm:"column:login"`
	Password string `json:"password,omitempty" gorm:"column:password"`
	Token    string `json:"token,omitempty" gorm:"column:token"`
}
