package types

import ()

type UserStoreInterface interface {
	GetUserByName(username string) (*User, error)
	GetUserByID(id int) (*User, error)
	CreateUser(RegisterUserPayload) error
}

type User struct {
	Userid      int    `json:"userid"`
	Username    string `json:"username"`
	Password    string `json:"password"`
	Createddate string `json:"createdDate"`
	Lastlogin   string `json:"lastLogin"`
}

type NewUser struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RegisterUserPayload struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginUserPayload struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
