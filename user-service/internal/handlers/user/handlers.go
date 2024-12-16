package user

import (
	"context"
	"errors"
	"fmt"

	ssov2 "github.com/mgrankin-cloud/messenger/contract/gen/go/user"
	"github.com/mgrankin-cloud/messenger/internal/domain/models"
	"github.com/mgrankin-cloud/messenger/internal/logger/validator"
	"github.com/mgrankin-cloud/messenger/pkg/storage"
	errors2 "github.com/mgrankin-cloud/messenger/pkg/storage"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type serverAPI struct {
	ssov2.UnimplementedUserServer
	user    User
	storage *storage.Storage
}

type User interface {
	GetUser(
		ctx context.Context,
		userID int64,
	) (models.User, error)
	UpdateUser(
		ctx context.Context,
		userID int64,
		email string,
		username string,
		phone string,
		photo []byte,
		active bool,
	) (success bool, err error)
	ChangePassword(
		ctx context.Context,
		userID int64,
		newPassword string,
	) (success bool, err error)
	DeleteUser(
		ctx context.Context,
		userID int64,
	) (success bool, err error)
	UpdateUserActiveStatus(
		ctx context.Context,
		userID int64,
		active bool,
	) (success bool, err error)
}

func RegisterUserService(gRPCServer *grpc.Server, user User) {
	ssov2.RegisterUserServer(gRPCServer, &serverAPI{user: user})
}

func (s *serverAPI) GetUser(
	ctx context.Context,
	request *ssov2.GetUserRequest,
) (*ssov2.GetUserResponse, error) {

	if request.GetUserId() == 0 {
		return nil, status.Error(codes.InvalidArgument, "user id is required")
	}

	user, err := s.user.GetUser(ctx, request.GetUserId())
	if err != nil {
		if errors.Is(err, errors2.ErrUserNotFound) {
			return nil, status.Error(codes.NotFound, "user not found")
		}
		return nil, status.Error(codes.Internal, "failed to get user")
	}

	return &ssov2.GetUserResponse{
		Email:    user.Email,
		Username: user.Username,
		Phone:    user.Phone,
		Photo:    user.Photo,
		Active:   user.Active,
	}, nil
}

func (s *serverAPI) UpdateUser(
	ctx context.Context,
	request *ssov2.UpdateUserInfoRequest,
) (*ssov2.UpdateUserInfoResponse, error) {
	if request.GetUserId() == 0 {
		return nil, status.Error(codes.InvalidArgument, "user id is required")
	}

	if err := validator.ValidateField(request.Email, "email", validator.IsValidEmail); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if err := validator.ValidateField(request.Username, "username", validator.IsValidUsername); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if err := validator.ValidateField(request.Phone, "phone", validator.IsValidPhoneNumber); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	success, err := s.user.UpdateUser(ctx, request.GetUserId(), request.GetEmail(), request.GetUsername(), request.GetPhone(), request.GetPhoto(), request.GetActive())
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			return nil, status.Errorf(codes.NotFound, "user not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to update user: %v", err)
	}

	return &ssov2.UpdateUserInfoResponse{
		Success: success,
	}, nil
}

func (s *serverAPI) ChangePassword(
	ctx context.Context,
	request *ssov2.ChangeUserPasswordRequest,
) (*ssov2.ChangeUserPasswordResponse, error) {
	if request.GetUserId() == 0 {
		return nil, status.Error(codes.InvalidArgument, "user id is required")
	}

	if err := validator.ValidateField(request.NewPassword, "password", validator.IsValidPassword); err != nil {
		fmt.Println(err)
	}

	success, err := s.user.ChangePassword(ctx, request.GetUserId(), request.GetNewPassword())
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			return nil, status.Errorf(codes.NotFound, "user not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to change password: %v", err)
	}

	return &ssov2.ChangeUserPasswordResponse{
		Success: success,
	}, nil
}

func (s *serverAPI) DeleteUser(
	ctx context.Context,
	request *ssov2.DeleteUserRequest,
) (*ssov2.DeleteUserResponse, error) {
	if request.GetUserId() == 0 {
		return nil, status.Error(codes.InvalidArgument, "user id is required")
	}

	success, err := s.user.DeleteUser(ctx, request.GetUserId())
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			return nil, status.Errorf(codes.NotFound, "user not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to delete user: %v", err)
	}

	return &ssov2.DeleteUserResponse{
		Success: success,
	}, nil
}

func (s *serverAPI) DeactivateUser(ctx context.Context, request *ssov2.GetActiveUserRequest) (*ssov2.GetActiveUserResponse, error) {
	if request.GetUserId() == 0 {
		return nil, status.Error(codes.InvalidArgument, "user id is required")
	}

	success,  err := s.user.UpdateUserActiveStatus(ctx, request.GetUserId(), false)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to deactivate user: %v", err)
	}

	return &ssov2.GetActiveUserResponse{
		Success: success,
		Active:  false,
	}, nil
}
