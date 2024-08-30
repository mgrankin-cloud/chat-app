package user

import (
	"Messenger-android/messenger/domain/models"
	"Messenger-android/messenger/storage"
	"context"
	"errors"
	"fmt"
	"log/slog"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
)

type User struct {
	log             *slog.Logger
	usrProvider     UsrProvider
	usrChanger      UsrChanger
	usrRemover      UsrRemover
	passwordChanger PasswordChanger
	appProvider     AppProvider
}

type UsrProvider interface {
	User(ctx context.Context, email, username, phone string) (models.User, error)
}

type UsrChanger interface {
	UpdateUser(ctx context.Context, userID int64, email, username, phone string, photo []byte) error
}

type UsrRemover interface {
	DeleteUser(ctx context.Context, userId int64) error
}

type PasswordChanger interface {
	ChangePassword(ctx context.Context, userID int64, newPassword string) error
}

type AppProvider interface {
	App(ctx context.Context, appID int) (models.App, error)
}

func New(
	log *slog.Logger,
	usrProvider UsrProvider,
	usrChanger UsrChanger,
	usrRemover UsrRemover,
	passwordChanger PasswordChanger,
	appProvider AppProvider,
) *User {
	return &User{
		log:             log,
		usrProvider:     usrProvider,
		usrChanger:      usrChanger,
		usrRemover:      usrRemover,
		passwordChanger: passwordChanger,
		appProvider:     appProvider,
	}
}

func (u *User) User(ctx context.Context, email, username, phone string) (models.User, error) {
	const op = "User.GetUser"

	var user models.User

	log := u.log.With(
		slog.String("op", op),
		slog.Int64("user_id", user.ID),
		slog.String("email", user.Email),
		slog.String("username", user.Username),
		slog.String("password", string(user.PassHash)),
		slog.String("photo", string(user.Photo)),
	)

	log.Info("attempting to get user")

	user, err := u.usrProvider.User(ctx, email, username, phone)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			u.log.Warn("user not found", err)
			return user, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		}

		u.log.Error("failed to get user", err)
		return user, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

func (u *User) UpdateUser(ctx context.Context, userID int64, email, username, phone string, photo []byte) error {
	const op = "User.UpdateUser"

	log := u.log.With(
		slog.String("op", op),
		slog.String("email", email),
		slog.String("username", username),
		slog.String("phone", phone),
		slog.String("photo", string(photo)),
	)

	log.Info("attempting to update user")

	err := u.usrChanger.UpdateUser(ctx, userID, email, username, phone, photo)
	if err != nil {
		log.Error("Failed to update user", slog.String("error", err.Error()))
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (u *User) DeleteUser(ctx context.Context, userID int64) error {
	const op = "User.DeleteUser"

	log := u.log.With(
		slog.String("op", op),
		slog.Int64("user_id", userID),
	)

	log.Info("attempting to delete user")

	err := u.usrRemover.DeleteUser(ctx, userID)
	if err != nil {
		log.Error("Failed to delete user", slog.String("error", err.Error()))
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (u *User) ChangePassword(ctx context.Context, userID int64, newPassword string) error {
	const op = "User.ChangePassword"

	log := u.log.With(
		slog.String("op", op),
		slog.Int64("user_id", userID),
		slog.String("password", newPassword),
	)

	log.Info("attempting to change password")

	err := u.passwordChanger.ChangePassword(ctx, userID, newPassword)
	if err != nil {
		log.Error("Failed to change password", slog.String("error", err.Error()))
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (u *User) GetAppSettings(ctx context.Context, appID int) (models.App, error) {
	const op = "User.GetAppSettings"

	log := u.log.With(
		slog.String("op", op),
		slog.Int("app_id", appID),
	)

	log.Info("attempting to get app settings")

	app, err := u.appProvider.App(ctx, appID)
	if err != nil {
		log.Error("Failed to get app settings", slog.String("error", err.Error()))
		return models.App{}, fmt.Errorf("%s: %w", op, err)
	}

	return app, nil
}
