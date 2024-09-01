package user

import (
	"Messenger-android/messenger/logger/validator"
	"Messenger-android/messenger/storage"
	ssov2 "Messenger-android/messenger/user-service/contract/gen/go/user"
	"context"
	"errors"
	"fmt"
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
		request *ssov2.GetUserRequest,
	) (*ssov2.GetUserResponse, error)
	UpdateUser(
		ctx context.Context,
		email string,
		username string,
		phone string,
		photo []byte,
	) (success bool, userID int64, err error)
	ChangePassword(
		ctx context.Context,
		password string,
		userID int64,
		newPassword string,
	) (success bool, err error)
	DeleteUser(
		ctx context.Context,
		request *ssov2.DeleteUserRequest,
	) (success bool, userID int64, err error)
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

	user, err := s.storage.GetUserByID(ctx, request.GetUserId())
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "user not found")
	}

	return &ssov2.GetUserResponse{
		UserId:   user.ID,
		Email:    user.Email,
		Username: user.Username,
		Phone:    user.Phone,
	}, nil
}

func (s *serverAPI) UpdateUser(
	ctx context.Context,
	request *ssov2.UpdateUserInfoRequest,
) (*ssov2.UpdateUserInfoResponse, error) {
	userID := request.GetUserId()
	email := request.GetEmail()
	username := request.GetUsername()
	phone := request.GetPhone()
	photo := request.GetPhoto()

	if userID == 0 {
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

	err := s.storage.UpdateUser(ctx, userID, email, username, phone, photo)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			return nil, status.Errorf(codes.NotFound, "user not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to update user: %v", err)
	}

	return &ssov2.UpdateUserInfoResponse{
		Success: true,
		UserId:  userID,
	}, nil
}

func (s *serverAPI) ChangePassword(
	ctx context.Context,
	request *ssov2.ChangeUserPasswordRequest,
) (*ssov2.ChangeUserPasswordResponse, error) {

	userID := request.GetUserId()
	newPassword := request.GetNewPassword()

	if userID == 0 {
		return nil, status.Error(codes.InvalidArgument, "user id is required")
	}

	if err := validator.ValidateField(request.NewPassword, "password", validator.IsValidPassword); err != nil {
		fmt.Println(err)
	}

	err := s.storage.ChangePassword(ctx, userID, newPassword)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			return nil, status.Errorf(codes.NotFound, "user not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to change password: %v", err)
	}

	return &ssov2.ChangeUserPasswordResponse{
		Success: true,
	}, nil
}

func (s *serverAPI) DeleteUser(
	ctx context.Context,
	request *ssov2.DeleteUserRequest,
) (*ssov2.DeleteUserResponse, error) {
	userID := request.GetUserId()

	if userID == 0 {
		return nil, status.Error(codes.InvalidArgument, "user id is required")
	}

	err := s.storage.DeleteUser(ctx, userID)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			return nil, status.Errorf(codes.NotFound, "user not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to delete user: %v", err)
	}

	return &ssov2.DeleteUserResponse{
		Success: true,
		UserId:  userID,
	}, nil
}
