package user

type PostUserRequest struct {
	User User
}

type PutUserRequest struct {
	User User
}

type DeleteUserRequest struct {
	UserID int
}

type GetUserRequest struct {
	Email string
}

type GetAllUsersRequest struct {
}
