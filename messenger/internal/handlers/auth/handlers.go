package auth

import (
	"context"
	"errors"
	ssov1 "github.com/mgrankin-cloud/messenger/contract/gen/go/sso"
	"github.com/mgrankin-cloud/messenger/internal/logger/validator"
	"github.com/mgrankin-cloud/messenger/internal/services/auth"
	errors2 "github.com/mgrankin-cloud/messenger/pkg/storage"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type serverAPI struct {
	ssov1.UnimplementedAuthServer
	auth Auth
}

// Auth интерфейс для передачи в gRPC app
type Auth interface {
	Authenticate(
		ctx context.Context,
		identifier string,
		password string,
		appID int,
	) (success bool, token string, err error)
	RegisterNewUser(
		ctx context.Context,
		email string,
		username string,
		password string,
		phone string,
	) (success bool, userID int64, err error)
}

func Register(gRPCServer *grpc.Server, auth Auth) {
	ssov1.RegisterAuthServer(gRPCServer, &serverAPI{auth: auth})
}

func (s *serverAPI) Authenticate(
	ctx context.Context,
	in *ssov1.AuthenticateRequest,
) (*ssov1.AuthenticateResponse, error) {
	if err := validator.ValidateField(in.Identifier, "identifier", validator.IsValidUsername); err != nil {
		return nil, err
	}

	if err := validator.ValidateField(in.Password, "password", validator.IsValidPassword); err != nil {
		return nil, err
	}

	if in.GetAppId() == 0 {
		return nil, status.Error(codes.InvalidArgument, "missing app id")
	}

	success, token, err := s.auth.Authenticate(ctx, in.GetIdentifier(), in.GetPassword(), int(in.GetAppId()))
	if err != nil {
		if errors.Is(err, auth.ErrInvalidCredentials) {
			return nil, status.Error(codes.InvalidArgument, "Invalid credentials")
		}

		return nil, status.Error(codes.Internal, "failed to login")
	}

	return &ssov1.AuthenticateResponse{Success: success, Token: token}, nil
}

func (s *serverAPI) Register(
	ctx context.Context,
	in *ssov1.RegisterRequest,
) (*ssov1.RegisterResponse, error) {

	if err := validator.ValidateField(in.Email, "email", validator.IsValidEmail); err != nil {
		return nil, err
	}

	if err := validator.ValidateField(in.Username, "username", validator.IsValidUsername); err != nil {
		return nil, err
	}

	if err := validator.ValidateField(in.Password, "password", validator.IsValidPassword); err != nil {
		return nil, err
	}

	if err := validator.ValidateField(in.Phone, "phone", validator.IsValidPhoneNumber); err != nil {
		return nil, err
	}

	success, uid, err := s.auth.RegisterNewUser(ctx, in.GetEmail(), in.GetUsername(), in.GetPassword(), in.GetPhone())
	if err != nil {
		if errors.Is(err, errors2.ErrUserExists) {
			return nil, status.Error(codes.AlreadyExists, "user already exists")
		}

		return nil, status.Error(codes.Internal, "failed to register user")
	}

	return &ssov1.RegisterResponse{Success: success, UserId: uid}, nil
}
