package api

import (
	"context"
	"jwtauthserver/auth"
	"jwtauthserver/pkg/errors"
	"jwtauthserver/pkg/rest"
	"net/http"

	"github.com/go-kit/kit/endpoint"
)

type Endpoint struct {
	RegisterEndpoint   endpoint.Endpoint
	AuthorizeEndpoint  endpoint.Endpoint
	GetProfileEndpoint endpoint.Endpoint
}

func MakeServerEndpoint(svc auth.Service) Endpoint {
	return Endpoint{
		RegisterEndpoint:  RegisterEndpoint(svc),
		AuthorizeEndpoint: AuthorizeEndpoint(svc),
	}
}

func RegisterEndpoint(svc auth.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(UserReqBody)

		if req.Username == "" {
			return nil, errors.ErrCreateEntity
		}

		if req.Password == "" {
			return nil, errors.ErrCreateEntity
		}

		err = svc.Register(ctx, req.Username, req.Password)
		if err != nil {
			return nil, err
		}

		return rest.HTTPResponse{
			Code:    http.StatusOK,
			Status:  "Success",
			Message: "Success Register User",
		}, nil
	}
}

func AuthorizeEndpoint(svc auth.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(UserReqBody)

		if req.Username == "" {
			return nil, errors.ErrCreateEntity
		}

		if req.Password == "" {
			return nil, errors.ErrCreateEntity
		}

		auth, err := svc.Authorize(ctx, req.Username, req.Password)
		if err != nil {
			return nil, err
		}

		return rest.HTTPResponse{
			Code:    http.StatusOK,
			Status:  "Success",
			Message: "Success Authorize User",
			Data: map[string]string{
				"token": auth.Token,
				"name":  auth.Username,
			},
		}, nil

	}
}

// func GetProfileEndpoint(svc auth.Service) endpoint.Endpoint {
// 	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
// 		req := request.(GetProfileReq)

// 		if req.Token == "" {
// 			return nil, errors.ErrUnauthorizedAccess
// 		}

// 		profile, err := svc.ViewAccount(ctx, req.Token)
// 		if err != nil {
// 			return nil, err
// 		}

// 		return rest.HTTPResponse{
// 			Code:   http.StatusOK,
// 			Status: "Success",
// 			Data:   profile,
// 		}, nil
// 	}
// }
