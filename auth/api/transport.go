package api

import (
	"context"
	"encoding/json"
	"errors"
	"jwtauthserver/auth"
	"jwtauthserver/pkg/rest"
	"net/http"

	"github.com/go-kit/log"

	"github.com/go-kit/kit/transport"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
)

var (
	// ErrBadRouting is returned when an expected path variable is missing.
	ErrBadRouting = errors.New("inconsistent mapping between route and handler ")
)

func MakeHTTPHandler(svc auth.Service, logger log.Logger) http.Handler {
	r := mux.NewRouter()
	r.Use(corsMiddleware)

	e := MakeServerEndpoint(svc)
	opt := []httptransport.ServerOption{
		httptransport.ServerErrorEncoder(rest.EncodeError),
		httptransport.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
	}

	r.Methods("POST").Path("/register").Handler(httptransport.NewServer(
		e.RegisterEndpoint,
		decodeUserRequest,
		rest.EncodeResponse,
		opt...,
	))

	r.Methods("POST").Path("/authorize").Handler(httptransport.NewServer(
		e.AuthorizeEndpoint,
		decodeUserRequest,
		rest.EncodeResponse,
		opt...,
	))

	return r
}

func decodeUserRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	var req UserReqBody
	if e := json.NewDecoder(r.Body).Decode(&req); e != nil {
		return nil, e
	}
	return req, nil
}

func corsMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type")

		if r.Method == "OPTIONS" {
			return
		}

		h.ServeHTTP(w, r)
	})
}

type UserReqBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type GetProfileReq struct {
	Token string
}
