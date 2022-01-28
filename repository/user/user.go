package user

//User - represents a user
type User struct {
	ID       int    `json:"id"`
	Email    string `json:"email" validate:"required,email"`
	Name     string `json:"name" validate:"required"`
	Lastname string `json:"lastname" validate:"required"`
	Age      int32  `json:"age,string"`
	Status   int32  `json:"status,string"`
}
