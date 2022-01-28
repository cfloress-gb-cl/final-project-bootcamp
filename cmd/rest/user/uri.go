package user

import "fmt"

var (
	UsersBaseUri = "/user/"
	PostUser     = fmt.Sprintf("%s", UsersBaseUri)
	GetUser      = fmt.Sprintf("%s{%s}", UsersBaseUri, "email")
	PutUser      = GetUser
	DeleteUser   = fmt.Sprintf("%s{%s}", UsersBaseUri, "id")
)
