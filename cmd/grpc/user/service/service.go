package grpc

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/go-kit/kit/log"

	"github.com/go-kit/kit/transport"
	grpctransport "github.com/go-kit/kit/transport/grpc"

	user "github.com/cfloress-gb-cl/final-project-bootcamp/cmd/grpc/user"
	endpoints "github.com/cfloress-gb-cl/final-project-bootcamp/cmd/grpc/user/endpoints"
	proto "github.com/cfloress-gb-cl/final-project-bootcamp/cmd/grpc/user/pb"
)

type grpcUserServer struct {
	proto.UsersServer
	getUser     grpctransport.Handler
	create      grpctransport.Handler
	getAllUsers grpctransport.Handler
	update      grpctransport.Handler
	delete      grpctransport.Handler
}

func NewGrpcUserServer(endpoints endpoints.GrpcUserServerEndpoints, logger log.Logger) proto.UsersServer {

	options := []grpctransport.ServerOption{
		grpctransport.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
	}

	server := &grpcUserServer{

		create:      grpctransport.NewServer(endpoints.CreateUserEndpoint, decodeCreateUserRequest, encodeCreateUserResponse, options...),
		getUser:     grpctransport.NewServer(endpoints.GetUserByEmailEndpoint, decodeGetUserRequest, encodeGetUserResponse, options...),
		getAllUsers: grpctransport.NewServer(endpoints.GetAllUsersEndpoint, decodeGetAllUsersRequest, encodeGetAllUsersResponse, options...),
		update:      grpctransport.NewServer(endpoints.UpdateUserEndpoint, decodeUpdateUserRequest, encodeUpdateUserResponse, options...),
		delete:      grpctransport.NewServer(endpoints.DeleteUserEndpoint, decodeDeleteUserRequest, encodeDeleteUserResponse, options...),
	}
	fmt.Println("Call to NewGrpcUserServer service.go grp pkg...")

	return server
}

func (u grpcUserServer) GetUser(ctx context.Context, uid *proto.EmailAddress) (*proto.GetUserResponse, error) {

	_, grpcResponse, err := u.getUser.ServeGRPC(ctx, uid)

	return grpcResponse.(*proto.GetUserResponse), err
}

func (u grpcUserServer) Create(ctx context.Context, user *proto.CreateUserRequest) (*proto.CreateUserResponse, error) {

	_, grpcResponse, err := u.create.ServeGRPC(ctx, user)

	fmt.Println("grpcResponse-->", grpcResponse)

	return grpcResponse.(*proto.CreateUserResponse), err

}

func (u grpcUserServer) GetAllUsers(ctx context.Context, request *proto.GetAllUsersRequest) (*proto.GetAllUsersResponse, error) {
	fmt.Println("getAllUsers method", request)

	_, grpcResponse, err := u.getAllUsers.ServeGRPC(ctx, request)

	return grpcResponse.(*proto.GetAllUsersResponse), err
}

func (u grpcUserServer) Update(ctx context.Context, userInfo *proto.UpdateUserRequest) (*proto.UpdateUserResponse, error) {

	fmt.Println("Update method")

	_, grpcResponse, err := u.update.ServeGRPC(ctx, userInfo)

	return grpcResponse.(*proto.UpdateUserResponse), err
}

func (u grpcUserServer) Delete(ctx context.Context, userId *proto.Id) (*proto.DeleteUserResponse, error) {

	ctx, grpcResponse, err := u.delete.ServeGRPC(ctx, userId)

	return grpcResponse.(*proto.DeleteUserResponse), err
}

func decodeCreateUserRequest(ctx context.Context, grpcReq interface{}) (interface{}, error) {
	req, validReq := grpcReq.(*proto.CreateUserRequest)
	if !validReq {
		return nil, errors.New("invalid input request data")
	}
	usr := user.User{Email: req.User.Email, Name: req.User.Name, Lastname: req.User.LastName, Age: req.User.Age, Status: req.User.Status}

	return user.PostUserRequest{User: usr}, nil
}

func encodeCreateUserResponse(ctx context.Context, resp interface{}) (interface{}, error) {

	response, validResp := resp.(user.PostUserResponse)

	if !validResp {
		return nil, errors.New("invalid input data")
	}

	if response.Error != nil {
		if response.Error.Error() == "user already exists" {
			return &proto.CreateUserResponse{Code: proto.CodeResult_FAILED}, nil
		} else {
			return &proto.CreateUserResponse{Code: proto.CodeResult_INVALIDINPUT}, nil
		}
	}

	return &proto.CreateUserResponse{UserId: int32(response.Id), Code: proto.CodeResult_OK}, nil
}

func decodeGetUserRequest(ctx context.Context, grpcReq interface{}) (interface{}, error) {
	req, validReq := grpcReq.(*proto.EmailAddress)
	if !validReq {
		return nil, errors.New("invalid input data")
	}
	return user.GetUserRequest{Email: req.Value}, nil
}

func encodeGetUserResponse(ctx context.Context, resp interface{}) (interface{}, error) {
	response, validResp := resp.(user.GetUserResponse)
	if !validResp {
		return nil, errors.New("invalid input data to encode")
	}
	usr := proto.User{Id: response.Id, Name: response.Name, Email: response.Email, LastName: response.Lastname, Age: response.Age, Status: response.Status}

	return &proto.GetUserResponse{User: &usr}, nil
}

func decodeGetAllUsersRequest(ctx context.Context, grpcReq interface{}) (interface{}, error) {
	_, validReq := grpcReq.(*proto.GetAllUsersRequest)
	if !validReq {
		return nil, errors.New("invalid input data decode")
	}
	return user.GetUserRequest{}, nil
}

func encodeGetAllUsersResponse(ctx context.Context, resp interface{}) (interface{}, error) {
	respData, validResp := resp.(user.GetAllUsersResponse)
	if !validResp {
		return nil, errors.New("invalid input data to encode")
	}

	response := &proto.GetAllUsersResponse{Users: []*proto.User{}}

	for _, usr := range respData.Users {
		pbUser := proto.User{Id: usr.Id, Name: usr.Name, Email: usr.Email, LastName: usr.Lastname, Age: usr.Age, Status: usr.Status}

		response.Users = append(response.Users, &pbUser)
	}

	return response, nil
}

func decodeUpdateUserRequest(ctx context.Context, grpcReq interface{}) (interface{}, error) {

	req, validReq := grpcReq.(*proto.UpdateUserRequest)

	fmt.Println("decode", req)

	if !validReq {
		return nil, errors.New("invalid input data to decode")
	}

	usr := user.User{Id: req.User.Id, Email: req.User.Email, Name: req.User.Name, Lastname: req.User.LastName, Age: req.User.Age, Status: req.User.Status}

	fmt.Println("user", usr)

	return user.UpdateUserRequest{User: usr}, nil
}

func encodeUpdateUserResponse(ctx context.Context, resp interface{}) (interface{}, error) {

	response, validResp := resp.(user.UpdateUserResponse)

	fmt.Println("response-->", response)

	if !validResp {
		fmt.Println("!validResp")
		return nil, errors.New("invalid input data to encode")
	}

	if response.Error != nil {
		switch {
		case strings.Contains(response.Error.Error(), "user not found"):
			return &proto.UpdateUserResponse{Code: proto.CodeResult_NOTFOUND}, nil
		case strings.Contains(response.Error.Error(), "cannot update the user information"):
			return &proto.UpdateUserResponse{Code: proto.CodeResult_FAILED}, nil
		case strings.Contains(response.Error.Error(), "or missing"):
			return &proto.UpdateUserResponse{Code: proto.CodeResult_MISSINGFIELD}, nil
		case strings.Contains(response.Error.Error(), "no records were updated"):
			return &proto.UpdateUserResponse{Code: proto.CodeResult_NOCHANGES}, nil
		default:
			return &proto.UpdateUserResponse{Code: proto.CodeResult_INVALIDINPUT}, nil
		}
	}

	return &proto.UpdateUserResponse{Code: proto.CodeResult_OK}, nil
}

func decodeDeleteUserRequest(ctx context.Context, req interface{}) (interface{}, error) {
	fmt.Println(req)
	request, validReq := req.(*proto.Id)

	if !validReq {
		return nil, errors.New("invalid input data to decode")
	}

	return user.DeleteUserRequest{Id: request.Value}, nil
}

func encodeDeleteUserResponse(ctx context.Context, resp interface{}) (interface{}, error) {
	fmt.Println(resp)
	response, validReq := resp.(user.DeleteUserResponse)

	if !validReq {
		return nil, errors.New("invalid input data to encode")
	}

	if response.Error != nil {
		errorMessage := response.Error.Error()

		switch {
		case strings.Contains(errorMessage, "user not found"):
			return &proto.DeleteUserResponse{Code: proto.CodeResult_NOTFOUND}, nil
		case strings.Contains(errorMessage, "cannot update the user information"):
			return &proto.DeleteUserResponse{Code: proto.CodeResult_FAILED}, nil
		case strings.Contains(errorMessage, "or missing"):
			return &proto.DeleteUserResponse{Code: proto.CodeResult_MISSINGFIELD}, nil
		case strings.Contains(errorMessage, "no records were"):
			return &proto.DeleteUserResponse{Code: proto.CodeResult_NOCHANGES}, nil
		default:
			return &proto.DeleteUserResponse{Code: proto.CodeResult_INVALIDINPUT}, nil
		}
	}

	return &proto.DeleteUserResponse{Code: proto.CodeResult_OK}, nil
}
