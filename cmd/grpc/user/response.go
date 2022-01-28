package grpc

type PostUserResponse struct {
	Error error
	Id    int
}

type GetUserResponse struct {
	User
}

type GetAllUsersResponse struct {
	Users []User
}

type UpdateUserResponse struct {
	Error error
}

type DeleteUserResponse struct {
	Error error
}
