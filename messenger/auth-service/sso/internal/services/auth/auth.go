package auth

import (
	"Messenger-android/messenger/auth-service/sso/internal/domain/models"
	"Messenger-android/messenger/auth-service/sso/internal/lib/jwt"
	"Messenger-android/messenger/auth-service/sso/internal/lib/logger/s1"
	"Messenger-android/messenger/auth-service/sso/internal/lib/otp"
	errors2 "Messenger-android/messenger/auth-service/sso/internal/storage"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/net/context"
	"log/slog"
	"strconv"
	"strings"
	"time"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
)

var redisClient *redis.Client

type Auth struct {
	log         *slog.Logger
	usrSaver    UserSaver
	usrProvider UserProvider
	appProvider AppProvider
	tokenTTL    time.Duration
}

type UserSaver interface {
	SaveUser(
		ctx context.Context,
		email string,
		username string,
		passHash []byte,
		phone string,
	) (uid int64, err error)
}

type UserProvider interface {
	User(ctx context.Context, email string, username string, phone string) (models.User, error)
}

type AppProvider interface {
	App(ctx context.Context, appID int) (models.App, error)
}

func Init() {
	redisClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
}

func New(
	log *slog.Logger,
	userSaver UserSaver,
	userProvider UserProvider,
	appProvider AppProvider,
	tokenTTL time.Duration,
) *Auth {
	return &Auth{
		usrSaver:    userSaver,
		usrProvider: userProvider,
		log:         log,
		appProvider: appProvider,
		tokenTTL:    tokenTTL,
	}
}

func (a *Auth) RegisterNewUser(ctx context.Context, email string, username string, pass string, phone string) (bool, int64, error) {
	const op = "Auth.RegisterNewUser"

	log := a.log.With(
		slog.String("op", op),
		slog.String("email", email))

	log.Info("Registering new user")

	passHash, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		log.Error("Failed to generate password hash", s1.Err(err))
		return false, 0, fmt.Errorf("%s: %w", op, err)
	}

	// Send OTP to user's email
	err = a.sendOTP(email)
	if err != nil {
		log.Error("Failed to send OTP email", s1.Err(err))
		return false, 0, fmt.Errorf("%s: %w", op, err)
	}

	var userOTP string

	if !otp.VerifyOTP(email, userOTP) {
		log.Error("Failed to verify OTP email")
		return false, 0, fmt.Errorf("%s: invalid OTP", op)
	}

	id, err := a.usrSaver.SaveUser(ctx, email, username, passHash, phone)
	if err != nil {
		log.Error("Failed to save new user", s1.Err(err))
		return false, 0, fmt.Errorf("%s: %w", op, err)
	}

	return true, id, nil
}

func (a *Auth) sendOTP(email string) error {
	// Generate OTP
	otpGen := otp.GenerateOTP()
	otp.StoreOTP(email, otpGen, 5*time.Minute)

	// Send OTP to user's email
	return otp.SendEmail(email, "Email Verification OTP", "Your OTP is: "+otpGen)
}

func (a *Auth) Authenticate(
	ctx context.Context,
	identifier string,
	password string,
	appID int,
) (success bool, token string, err error) {
	const op = "Auth.Authenticate"

	var email, username, phone string

	if strings.Contains(identifier, "@") {
		email = identifier
	} else if _, err := strconv.Atoi(identifier); err != nil {
		phone = identifier
	} else {
		username = identifier
	}

	log := a.log.With(
		slog.String("op", op),
		slog.String("identifier", identifier),
	)

	log.Info("attempting to login user")

	// проверяем количество попыток входа
	attempts, err := a.checkLoginAttempts(ctx, identifier)
	if err != nil {
		a.log.Error("failed to check login attempts", s1.Err(err))
		return false, "", fmt.Errorf("%s: %w", op, err)
	}
	var maxLoginAttempts = 7
	if attempts >= maxLoginAttempts {
		a.log.Warn("too many login attempts", s1.Err(err))
		return false, "", fmt.Errorf("%s: %w", op, err)
	}

	// получение пользователя из бд
	user, err := a.usrProvider.User(ctx, email, username, phone)
	if err != nil {
		if errors.Is(err, errors2.ErrUserNotFound) {
			a.log.Warn("user not found", s1.Err(err))
			a.recordLoginAttempt(ctx, identifier, false)
			return false, "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		}

		a.log.Error("failed to get user", s1.Err(err))
		return false, "", fmt.Errorf("%s: %w", op, err)
	}

	if err := bcrypt.CompareHashAndPassword(user.PassHash, []byte(password)); err != nil {
		a.log.Info("invalid password", s1.Err(err))
		a.recordLoginAttempt(ctx, identifier, false)
		return false, "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}

	// Send OTP to user's email
	err = a.sendOTP(email)
	if err != nil {
		log.Error("Failed to send OTP email", s1.Err(err))
		return false, "", fmt.Errorf("%s: %w", op, err)
	}

	var userOTP string

	if !otp.VerifyOTP(email, userOTP) {
		log.Error("Failed to verify OTP email")
		return false, "", fmt.Errorf("%s: invalid OTP", op)
	}

	app, err := a.appProvider.App(ctx, appID)
	if err != nil {
		return false, "", fmt.Errorf("%s: %w", op, err)
	}

	log.Info("user logged successfully")

	successed, token, err := jwt.NewToken(user, app, a.tokenTTL)
	if err != nil {
		a.log.Error("failed to generate token", s1.Err(err))
		return false, "", fmt.Errorf("%s: %w", op, err)
	}

	// запись успешной попытки входа
	a.recordLoginAttempt(ctx, identifier, true)

	return successed, token, nil
}

func (a *Auth) checkLoginAttempts(ctx context.Context, identifier string) (int, error) {
	key := "login_attempts: " + identifier
	attempts, err := redisClient.Get(ctx, key).Int()
	if errors.Is(err, redis.Nil) {
		return 0, nil
	} else if err != nil {
		return 0, err
	}
	return attempts, nil
}

func (a *Auth) recordLoginAttempt(ctx context.Context, identifier string, success bool) {
	key := "login_attempts:" + identifier
	if success {
		redisClient.Del(ctx, key)
	} else {
		redisClient.Incr(ctx, key)
		redisClient.Expire(ctx, key, 5*time.Minute)
	}
}
