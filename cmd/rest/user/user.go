package user

type User struct {
	Id       int    `json:"id,omitempty"`
	Email    string `json:"email,omitempty"`
	Name     string `json:"name,omitempty"`
	Lastname string `json:"lastname,omitempty"`
	Age      int32  `json:"age,string"`
	Status   int32  `json:"status,string"`
}
