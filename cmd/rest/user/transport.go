package user

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	errUser "github.com/cfloress-gb-cl/final-project-bootcamp/cmd/rest/user/errors"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/transport"
	httptransport "github.com/go-kit/kit/transport/http"
)

var (
	ErrBadRouting = errors.New("bad request")
)

func MakeHTTPHandler(s UserProxy, logger log.Logger) http.Handler {
	r := mux.NewRouter()
	e := MakeServerEndpoints(s)
	options := []httptransport.ServerOption{
		httptransport.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
		httptransport.ServerErrorEncoder(encodeError),
	}

	r.Methods(http.MethodPost).Path(PostUser).Handler(httptransport.NewServer(
		e.PostUserEndpoint,
		decodePostProfileRequest,
		encodeResponse,
		options...,
	))
	r.Methods(http.MethodGet).Path(GetUser).Handler(httptransport.NewServer(
		e.GetUserEndpoint,
		decodeGetUserRequest,
		encodeResponse,
		options...,
	))
	r.Methods(http.MethodGet).Path(UsersBaseUri).Handler(httptransport.NewServer(
		e.GetAllUsersEndpoint,
		decodeGetAllUsersRequest,
		encodeResponse,
		options...,
	))
	r.Methods(http.MethodPut).Path(PutUser).Handler(httptransport.NewServer(
		e.PutUserEndpoint,
		decodePutProfileRequest,
		encodeResponse,
		options...,
	))
	r.Methods(http.MethodDelete).Path(DeleteUser).Handler(httptransport.NewServer(
		e.DeleteUserEndpoint,
		decodeDeleteProfileRequest,
		encodeResponse,
		options...,
	))

	return r
}

//decoders
func decodeGetUserRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	vars := mux.Vars(r)
	email, ok := vars["email"]
	if !ok {
		return nil, ErrBadRouting
	}
	return GetUserRequest{Email: email}, nil
}

func decodeGetAllUsersRequest(_ context.Context, r *http.Request) (request interface{}, err error) {

	return GetAllUsersRequest{}, nil
}

func decodePostProfileRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	var req PostUserRequest
	if e := json.NewDecoder(r.Body).Decode(&req.User); e != nil {
		return nil, e
	}
	return req, nil
}

func decodePutProfileRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	fmt.Println("decodePutProfileRequest rest user transport")
	vars := mux.Vars(r)
	email, ok := vars["email"]
	if !ok {
		return nil, ErrBadRouting
	}
	var usr User
	if err := json.NewDecoder(r.Body).Decode(&usr); err != nil {
		return nil, err
	}
	usr.Email = email
	return PutUserRequest{
		User: usr,
	}, nil
}

func decodeDeleteProfileRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	vars := mux.Vars(r)
	id, ok := strconv.Atoi(vars["id"])
	if ok != nil {
		return nil, ErrBadRouting
	}
	return DeleteUserRequest{UserID: id}, nil
}

type errorer interface {
	error() error
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {

	if e, ok := response.(errorer); ok && e.error() != nil {
		encodeError(ctx, e.error(), w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	if err == nil {
		panic("encodeError with nil error")
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(codeFrom(err))
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}

func codeFrom(err error) int {
	switch err {
	case errUser.ErrNotFound:
		return http.StatusNotFound
	case errUser.ErrUserAlreadyExists:
		return http.StatusConflict
	case errUser.ErrInvalidInput:
		return http.StatusUnprocessableEntity
	case ErrBadRouting:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}
