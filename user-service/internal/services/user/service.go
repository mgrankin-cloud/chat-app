package user

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/mgrankin-cloud/messenger/internal/logger/s1"
	errors2 "github.com/mgrankin-cloud/messenger/pkg/storage"

	"github.com/mgrankin-cloud/messenger/internal/domain/models"
	"github.com/mgrankin-cloud/messenger/pkg/storage"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
)

type User struct {
	log         *slog.Logger
	usrProvider UsrProvider
	appProvider AppProvider
}

type UsrProvider interface {
	GetUserByID(ctx context.Context, userID int64) (models.User, error)
	UpdateUser(ctx context.Context, id int64, email, username, phone string, photo []byte, active bool) (bool, error)
	UpdateUserActiveStatus(ctx context.Context, userID int64, active bool) (bool, error)
	DeleteUser(ctx context.Context, userId int64) (bool, error)
	ChangePassword(ctx context.Context, userID int64, newPassword string) (bool, error)
}

type AppProvider interface {
	App(ctx context.Context, appID int64) (models.App, error)
}

func New(
	log *slog.Logger,
	usrProvider UsrProvider,
	appProvider AppProvider,
) *User {
	return &User{
		log:         log,
		usrProvider: usrProvider,
		appProvider: appProvider,
	}
}

func (u *User) GetUser(ctx context.Context, userID int64) (models.User, error) {
	const op = "User.GetUser"

	var user models.User

	log := u.log.With(
		slog.String("op", op),
		slog.String("email", user.Email),
		slog.String("username", user.Username),
		slog.Bool("active", user.Active),
	)

	log.Info("attempting to get user")

	user, err := u.usrProvider.GetUserByID(ctx, userID)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			u.log.Warn("user not found")
			return user, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		}

		u.log.Error("failed to get user")
		return user, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

func (u *User) UpdateUser(ctx context.Context, id int64, email, username, phone string, photo []byte, active bool) (bool, error) {
	const op = "User.UpdateUser"

	log := u.log.With(
		slog.String("op", op),
		slog.String("email", email),
		slog.String("username", username),
		slog.String("phone", phone),
		slog.String("photo", string(photo)),
		slog.Bool("active", active),
	)

	log.Info("attempting to update user")

	success, err := u.usrProvider.UpdateUser(ctx, id, email, username, phone, photo, active)
	if err != nil {
		if errors.Is(err, errors2.ErrUserNotFound) {
			u.log.Warn("user not found", s1.Err(err))
			return false, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		}

		u.log.Error("failed to get user", s1.Err(err))
		return false, fmt.Errorf("%s: %w", op, err)
	}

	return success, nil
}

func (u *User) DeleteUser(ctx context.Context, userID int64) (bool, error) {
	const op = "User.DeleteUser"

	log := u.log.With(
		slog.String("op", op),
		slog.Int64("user_id", userID),
	)

	log.Info("attempting to delete user")

	success, err := u.usrProvider.DeleteUser(ctx, userID)
	if err != nil {
		log.Error("Failed to delete user", slog.String("error", err.Error()))
		return false, fmt.Errorf("%s: %w", op, err)
	}

	return success, nil
}

func (u *User) ChangePassword(ctx context.Context, userID int64, newPassword string) (bool, error) {
	const op = "User.ChangePassword"

	log := u.log.With(
		slog.String("op", op),
		slog.Int64("user_id", userID),
		slog.String("password", newPassword),
	)

	log.Info("attempting to change password")

	success, err := u.usrProvider.ChangePassword(ctx, userID, newPassword)
	if err != nil {
		log.Error("Failed to change password", slog.String("error", err.Error()))
		return false, fmt.Errorf("%s: %w", op, err)
	}

	return success, nil
}

func (u *User) UpdateUserActiveStatus(ctx context.Context, userID int64, active bool) (bool, error) {
	const op = "user.UpdateUserActiveStatus"

	log := u.log.With(
		slog.String("op", op),
		slog.Int64("user_id", userID),
	)

	log.Warn("attempting to update user status")

	success, err := u.usrProvider.UpdateUserActiveStatus(ctx, userID, active)
	if err != nil {
		log.Error("Failed to update user status", slog.String("error", err.Error()))
		return false, fmt.Errorf("%s: %w", op, err)
	}

	return success, nil
}

func (u *User) GetAppSettings(ctx context.Context, appID int64) (models.App, error) {
	const op = "User.GetAppSettings"

	log := u.log.With(
		slog.String("op", op),
		slog.Int64("app_id", appID),
	)

	log.Info("attempting to get app settings")

	app, err := u.appProvider.App(ctx, appID)
	if err != nil {
		log.Error("Failed to get app settings", slog.String("error", err.Error()))
		return models.App{}, fmt.Errorf("%s: %w", op, err)
	}

	return app, nil
}
