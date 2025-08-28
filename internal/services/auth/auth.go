package auth

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"sso-auth/internal/domain/models"
	"sso-auth/internal/lib/jwt"
	"sso-auth/internal/storage"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type Auth struct {
	log *slog.Logger
	userSaver    UserSaver
	userProvider UserProvider
	appProvider  AppProvider
	tokenTTL     time.Duration
}

type UserSaver interface {
	SaveUser(
		ctx context.Context,
		email string,
		passHash []byte,
	) (uid int64, err error)
}
type UserProvider interface {
	User(
		ctx context.Context,
		email string,
	) (models.User, error)

	IsAdmin(ctx context.Context, userID int64) (bool, error)
}

type AppProvider interface {
	App(ctx context.Context, appID int64) (models.App, error)
	
}

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidAppID = errors.New("invalid app id")
	ErrUserExists = errors.New("user with this email already exists")	
)

func New(
	log *slog.Logger, 
	userSaver UserSaver, 
	userProvider UserProvider, 
	appProvider AppProvider, 
	tolenTTL time.Duration,
) *Auth {
	return &Auth{
		log: log,
		userSaver:   userSaver,
		userProvider: userProvider,
		appProvider: appProvider,	
		tokenTTL: tolenTTL,
	}
}

func (a *Auth) Login(
	ctx context.Context,
	email string,
	password string,
	appID int64,
) (string, error) {
	const op = "auth.Login"
	log := a.log.With(slog.String("op", op), slog.String("username", email))

	log.Info("login attempt")

	user, err := a.userProvider.User(ctx, email)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			a.log.Warn("user not found", slog.String("email", email))

			return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		}

		a.log.Error("failed to get user", slog.String("error", err.Error()))
		return "", fmt.Errorf("%s: %w", op, err)
	}

	if err := bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(password)); err != nil {
		a.log.Info("invalid password", slog.String("email", email))
		return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)	
	}

	app, err := a.appProvider.App(ctx, appID)

	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	log.Info("user logged in")

	token, err := jwt.NewToken(user, app, a.tokenTTL)
	if err != nil {
		a.log.Error("failed to generate token", slog.String("error", err.Error()))
		return "", fmt.Errorf("%s: %w", op, err)
	}
	return token, nil

}


func (a *Auth) RegisterNewUser(
	ctx context.Context,
	email string,
	password string,
) (int64, error) {
	const op = "auth.RegisterNewUser"

	log := a.log.With(slog.String("op", op), slog.String("email", email))

	log.Info("registering new user")

	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		log.Error("failed to hash password", slog.String("error", err.Error()))
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	id, err := a.userSaver.SaveUser(ctx, email, passHash)
	if err != nil {
		if errors.Is(err, storage.ErrUserExists) {
			log.Warn("user already exists", slog.String("email", email))

			return 0, fmt.Errorf("%s: %w", op, ErrUserExists)
		}

		log.Error("failed to save user", slog.String("error", err.Error()))
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("user registered", slog.Int64("user_id", id))
	return id, nil	
}

func (a *Auth) IsAdmin(ctx context.Context, userID int64) (bool, error) {
	const op = "auth.IsAdmin"

	log := a.log.With(slog.String("op", op), slog.Int64("user_id", userID))

	log.Info("checking if user is admin")

	isAdmin, err := a.userProvider.IsAdmin(ctx, userID)
	if err != nil {
		if errors.Is(err, storage.ErrAppNotFound) {
			log.Warn("app not found", slog.Int64("user_id", userID))
			return false, fmt.Errorf("%s: %w", op, ErrInvalidAppID)
		}
		return false, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("checked if user is admin", slog.Bool("is_admin", isAdmin))
	return isAdmin, nil
}






