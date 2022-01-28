package user

import (
	"context"
	"errors"
	"fmt"
	errUser "github.com/cfloress-gb-cl/final-project-bootcamp/cmd/rest/user/errors"
	"github.com/go-kit/kit/endpoint"
)

type Endpoints struct {
	PostUserEndpoint     endpoint.Endpoint
	PostManyUserEndpoint endpoint.Endpoint
	GetUserEndpoint      endpoint.Endpoint
	GetAllUsersEndpoint  endpoint.Endpoint
	PutUserEndpoint      endpoint.Endpoint
	DeleteUserEndpoint   endpoint.Endpoint
}

func MakeServerEndpoints(s GrpcUsersProxy) Endpoints {
	return Endpoints{
		PostUserEndpoint:     MakePostUserEndpoint(s),
		PostManyUserEndpoint: MakePostUserEndpoint(s),
		GetUserEndpoint:      MakeGetUserEndpoint(s),
		GetAllUsersEndpoint:  MakeGetAllUsersEndpoint(s),
		PutUserEndpoint:      MakePutUserEndpoint(s),
		DeleteUserEndpoint:   MakeDeleteUserEndpoint(s),
	}
}

func MakePostUserEndpoint(s GrpcUsersProxy) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {

		// TODO change variables name and change comments
		req, validReq := request.(PostUserRequest)
		if !validReq {
			return nil, errors.New("invalid input data")
		}

		usr, e := s.Create(ctx, req.User)

		if e != nil {
			return nil, e
		}

		return postUserResponse{User: fmt.Sprintf("%s", usr.Email), Err: e}, nil
	}
}

func MakeGetUserEndpoint(s GrpcUsersProxy) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {

		req, validReq := request.(GetUserRequest)
		if !validReq {
			return nil, errors.New("invalid input data")
		}
		p, e := s.GetByEmail(ctx, req.Email) //pasar el context hasta el grpc

		if e != nil {
			return nil, e
		}

		return getUserResponse{User: p, Err: e}, nil
	}
}

func MakeGetAllUsersEndpoint(s GrpcUsersProxy) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {

		_, validReq := request.(GetAllUsersRequest)
		if !validReq {
			return nil, errors.New("invalid input data")
		}
		p, e := s.GetAll(ctx) //pasar el context hasta el grpc

		if e != nil {
			return errUser.WrapError(e), nil
		}

		return getAllUsersResponse{Users: p, Err: e}, nil
	}
}

func MakePutUserEndpoint(s GrpcUsersProxy) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		fmt.Println("makePutUserEndpoints rest user")
		req, validReq := request.(PutUserRequest)
		if !validReq {
			return nil, errors.New("invalid input data")
		}
		_, e := s.Update(ctx, req.User)

		if e != nil {
			return nil, e
		}

		return putUserResponse{User: fmt.Sprintf("%v", req.User.Email), Err: e}, nil
	}
}

func MakeDeleteUserEndpoint(s GrpcUsersProxy) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req, validReq := request.(DeleteUserRequest)
		if !validReq {
			return nil, errors.New("invalid input data")
		}
		_, e := s.Delete(ctx, req.UserID)

		if e != nil {
			return nil, e
		}

		return deleteUserResponse{User: fmt.Sprintf("%v", req.UserID), Err: e}, nil
	}
}
