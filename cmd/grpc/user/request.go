package grpc

type PostUserRequest struct {
	User `json:"user,omitempty"`
}

type UpdateUserRequest struct {
	User `json:"user,omitempty"`
}

type DeleteUserRequest struct {
	Id int32 `json:"id,omitempty"`
}

type GetUserRequest struct {
	Email string `json:"email,omitempty"`
}

type GetAllUsersRequest struct {
}

type User struct {
	//The user id to update
	Id int32 `json:"id,omitempty"`
	//The user email
	Email string `json:"email,omitempty"`
	//The user name
	Name string `json:"name,omitempty"`
	//The user last name
	Lastname string `json:"last_name,omitempty"`
	// user's age
	Age int32 `json:"age"`
	//user's status
	Status int32 `json:"status"`
}
