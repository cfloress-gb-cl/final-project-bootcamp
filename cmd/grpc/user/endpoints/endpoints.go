package grpc

import (
	"context"
	"errors"
	"fmt"
	domain "github.com/cfloress-gb-cl/final-project-bootcamp/cmd/grpc/user"
	"github.com/cfloress-gb-cl/final-project-bootcamp/repository/user"
	"github.com/go-kit/kit/endpoint"
)

type GrpcUserServerEndpoints struct {
	CreateUserEndpoint     endpoint.Endpoint
	GetUserByEmailEndpoint endpoint.Endpoint
	UpdateUserEndpoint     endpoint.Endpoint
	DeleteUserEndpoint     endpoint.Endpoint
	GetAllUsersEndpoint    endpoint.Endpoint
}

func NewGrpcUsersServer(s user.Service) *GrpcUserServerEndpoints {
	fmt.Println("call to NewGrpcUsersServer endpoints")
	return &GrpcUserServerEndpoints{
		CreateUserEndpoint:     MakePostUserEndpoint(s),
		GetUserByEmailEndpoint: MakeGetUserEndpoint(s),
		UpdateUserEndpoint:     MakeUpdateUserEndpoint(s),
		DeleteUserEndpoint:     MakeDeleteUserEndpoint(s),
		GetAllUsersEndpoint:    MakeGetAllUsersEndpoint(s),
	}
}

func MakePostUserEndpoint(s user.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req, validReq := request.(domain.PostUserRequest)
		if !validReq {
			return nil, errors.New("invalid input data")
		}

		usr := user.User{Email: req.Email, Name: req.Name, Lastname: req.Lastname, Age: req.Age, Status: req.Status}

		usrID, err := s.Create(ctx, usr)

		fmt.Println("Error-->", domain.PostUserResponse{Id: usrID, Error: err})

		return domain.PostUserResponse{Id: usrID, Error: err}, nil
	}
}

func MakeGetUserEndpoint(s user.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {

		req, validReq := request.(domain.GetUserRequest)
		if !validReq {
			return nil, errors.New("invalid input data")
		}

		usr, err := s.GetByEmail(ctx, req.Email)

		return domain.GetUserResponse{User: domain.User{Id: int32(usr.ID), Email: usr.Email, Name: usr.Name, Lastname: usr.Lastname, Age: usr.Age, Status: usr.Status}}, nil
	}
}

func MakeGetAllUsersEndpoint(s user.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {

		usrs, err := s.GetAll(ctx)

		responseData := domain.GetAllUsersResponse{Users: []domain.User{}}

		for _, usr := range usrs {
			responseData.Users = append(responseData.Users, domain.User{Id: int32(usr.ID), Email: usr.Email, Name: usr.Name, Lastname: usr.Lastname, Age: usr.Age, Status: usr.Status})
		}

		fmt.Println("users", responseData)

		return responseData, nil
	}
}

func MakeUpdateUserEndpoint(s user.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {

		fmt.Println("Call to update user  Endpoint.go grp pkg...")
		req, validReq := request.(domain.UpdateUserRequest)

		fmt.Println(req)

		if !validReq {
			return nil, errors.New("invalid request type")
		}

		usr := user.User{Email: req.Email, Name: req.Name, Lastname: req.Lastname, Age: req.Age, Status: req.Status}

		if err != nil {
			return nil, errors.New("invalid object cast")
		}

		err = s.Update(ctx, usr)

		return domain.UpdateUserResponse{Error: err}, nil
	}
}

func MakeDeleteUserEndpoint(s user.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req, validReq := request.(domain.DeleteUserRequest)

		if !validReq {
			return nil, errors.New("invalid request type")
		}

		err = s.Delete(ctx, int(req.Id))

		return domain.DeleteUserResponse{Error: err}, nil
	}
}
