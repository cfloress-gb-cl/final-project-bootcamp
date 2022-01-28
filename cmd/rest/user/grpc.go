package user

import (
	"context"
	"fmt"
	glog "google.golang.org/grpc/grpclog"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/caarlos0/env/v6"

	proto "github.com/cfloress-gb-cl/final-project-bootcamp/cmd/grpc/user/pb"
	errUser "github.com/cfloress-gb-cl/final-project-bootcamp/cmd/rest/user/errors"
	"google.golang.org/grpc"
)

type ServerConnection struct {
	client  proto.UsersClient
	context context.Context
	dispose func()
}

type GrpcUsersProxy interface {
	GetAll(context.Context) ([]User, error)
	Create(context.Context, User) (User, error)
	Update(context.Context, User) (User, error)
	Delete(context.Context, int) (bool, error)
	GetByEmail(context.Context, string) (User, error)
}

type UserProxy struct {
	grpcLog glog.LoggerV2
}

type config struct {
	Port int    `env:"proto_PORT" envDefault:"9000"`
	Host string `env:"proto_HOST" envDefault:"127.0.0.1"`
}

func NewUserProxy() *UserProxy {
	return &UserProxy{
		grpcLog: glog.NewLoggerV2(os.Stdout, os.Stdout, os.Stdout),
	}
}

func (up UserProxy) GetAll(ctx context.Context) ([]User, error) {

	serverCon, err := OpenServerConection(ctx)

	if err != nil {
		log.Fatalf(err.Error())
	}

	defer serverCon.dispose()
	c := serverCon.client
	result, errorFromCall := c.GetAllUsers(serverCon.context, &proto.GetAllUsersRequest{})

	if errorFromCall != nil {
		log.Fatalf(errorFromCall.Error())
	}

	response := []User{}

	for _, o := range result.Users {
		response = append(response, User{
			Id:       int(o.Id),
			Email:    o.Email,
			Name:     o.Name,
			Lastname: o.LastName,
			Age:      o.Age,
			Status:   o.Status,
		})
	}

	return response, errorFromCall
}

func (up UserProxy) Create(ctx context.Context, u User) (User, error) {

	serverCon, err := OpenServerConection(ctx)

	if err != nil {
		log.Fatalf("did not connect to server: %s", err)
	}

	defer serverCon.dispose()
	c := serverCon.client
	externalUser := &proto.User{
		Id:       0,
		Email:    u.Email,
		Name:     u.Name,
		LastName: u.Lastname,
		Age:      u.Age,
		Status:   u.Status,
	}

	result, errorFromCall := c.Create(serverCon.context, &proto.CreateUserRequest{User: externalUser})

	if result.Code == proto.CodeResult_FAILED {
		return User{}, errUser.ErrUserAlreadyExists
	}

	if result.Code == proto.CodeResult_INVALIDINPUT {
		return User{}, errUser.ErrInvalidInput
	}

	if errorFromCall != nil {
		return User{}, errorFromCall
	}

	u.Id = int(result.UserId)
	return u, errorFromCall
}

func (up UserProxy) Update(ctx context.Context, u User) (User, error) {

	serverCon, err := OpenServerConection(ctx)

	if err != nil {
		log.Fatalf("did not connect to server: %s", err)
	}

	defer serverCon.dispose()
	c := serverCon.client
	externalUser := proto.User{
		Id:       int32(u.Id),
		Email:    u.Email,
		Name:     u.Name,
		LastName: u.Lastname,
		Age:      u.Age,
		Status:   u.Status,
	}

	result, errorFromCall := c.Update(serverCon.context, &proto.UpdateUserRequest{User: &externalUser})

	switch result.Code {
	case proto.CodeResult_FAILED:
		return User{}, errUser.ErrUserAlreadyExists
	case proto.CodeResult_INVALIDINPUT:
		return User{}, errUser.ErrInvalidInput
	case proto.CodeResult_NOTFOUND:
		return User{}, errUser.ErrNotFound
	case proto.CodeResult_MISSINGFIELD:
		return User{}, errUser.ErrNoMissingField
	case proto.CodeResult_NOCHANGES:
		return User{}, errUser.ErrNoChangesUpdate
	}

	if errorFromCall != nil {
		return User{}, errorFromCall
	}

	return u, nil
}

func (up UserProxy) Delete(ctx context.Context, id int) (bool, error) {

	serverCon, err := OpenServerConection(ctx)

	if err != nil {
		log.Fatalf("did not connect to server: %s", err)
	}

	defer serverCon.dispose()
	c := serverCon.client
	externalUserId := &proto.Id{
		Value: int32(id),
	}
	result, errorFromCall := c.Delete(serverCon.context, externalUserId)

	switch result.Code {
	case proto.CodeResult_FAILED:
		return false, errUser.ErrUserAlreadyExists
	case proto.CodeResult_INVALIDINPUT:
		return false, errUser.ErrInvalidInput
	case proto.CodeResult_NOTFOUND:
		return false, errUser.ErrNotFound
	case proto.CodeResult_MISSINGFIELD:
		return false, errUser.ErrNoMissingField
	case proto.CodeResult_NOCHANGES:
		return false, errUser.ErrNoChangesUpdate
	}

	if errorFromCall != nil {
		return true, errorFromCall
	}

	return false, errorFromCall
}

func (up UserProxy) GetByEmail(ctx context.Context, email string) (User, error) {

	serverCon, err := OpenServerConection(ctx)

	if err != nil {
		log.Fatalf("did not connect to server: %s", err)
	}

	defer serverCon.dispose()
	c := serverCon.client
	result, errorFromCall := c.GetUser(serverCon.context, &proto.EmailAddress{Value: email})

	if result.User.Id == 0 {
		return User{}, errUser.ErrNotFound
	}

	if errorFromCall != nil {
		fmt.Println("server call fail:", errorFromCall)
		return User{}, errorFromCall
	}

	userFromGrpc := result.User

	response := User{
		Id:       int(userFromGrpc.Id),
		Email:    userFromGrpc.Email,
		Name:     userFromGrpc.Name,
		Lastname: userFromGrpc.LastName,
		Age:      userFromGrpc.Age,
		Status:   userFromGrpc.Status,
	}

	return response, nil
}

func OpenServerConection(ctx context.Context) (*ServerConnection, error) {

	cfg := config{}

	if err := env.Parse(&cfg); err != nil {
		fmt.Printf("%+v\n", err)
	}

	conn, err := grpc.Dial(fmt.Sprintf("%s:%s", cfg.Host, strconv.Itoa(cfg.Port)), grpc.WithInsecure())

	if err != nil {
		log.Fatalf("did not connect to server: %s", err)
		return nil, err
	}

	ctxTO, cancel := context.WithTimeout(ctx, 10*time.Second)

	c := proto.NewUsersClient(conn)

	return &ServerConnection{
		client:  c,
		context: ctxTO,
		dispose: func() {
			cancel()
			conn.Close()
		},
	}, nil
}
