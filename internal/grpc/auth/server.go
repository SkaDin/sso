package auth

import (
	"context"
	"errors"
	ssov1 "github.com/SkaDin/protos/gen/go/sso"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/mail"
	"sso/internal/services/auth"
	"sso/internal/storage"
	"strings"
)

type Auth interface {
	Login(ctx context.Context,
		email string,
		password string,
		appID int,
	) (token string, err error)

	RegisterNewUser(ctx context.Context,
		email string,
		password string,
	) (userID int64, err error)

	IsAdmin(ctx context.Context, userID int64) (bool, error)
}

const (
	emptyValue   = 0
	invalidEmail = "email@example.com"
)

type serverAPI struct {
	ssov1.UnimplementedAuthServer
	auth Auth
}

func Register(gRPC *grpc.Server, auth Auth) {
	ssov1.RegisterAuthServer(gRPC, &serverAPI{auth: auth})
}

func (s *serverAPI) Login(
	ctx context.Context,
	req *ssov1.LoginRequest,
) (*ssov1.LoginResponse, error) {
	if err := validateLogin(req); err != nil {
		return nil, err
	}
	token, err := s.auth.Login(ctx, req.GetEmail(), req.GetPassword(), int(req.GetAppId()))
	if err != nil {
		if errors.Is(err, auth.ErrInvalidCredentials) {
			return nil, status.Error(codes.InvalidArgument, "invalid email or password")
		}

		return nil, status.Error(codes.Internal, "failed to login")
	}

	return &ssov1.LoginResponse{Token: token}, nil
}

func (s *serverAPI) Register(
	ctx context.Context,
	req *ssov1.RegisterRequest,
) (*ssov1.RegisterResponse, error) {

	if err := validateRegister(req); err != nil {
		return nil, err
	}

	userID, err := s.auth.RegisterNewUser(ctx, req.GetEmail(), req.GetPassword())
	if err != nil {
		if isDuplicateEmailError(err) {
			return nil, status.Error(codes.AlreadyExists, "user already exists")
		}

		return nil, status.Error(codes.Internal, "failed to register user")
	}
	return &ssov1.RegisterResponse{UserId: userID}, nil
}

func (s *serverAPI) IsAdmin(
	ctx context.Context,
	req *ssov1.IsAdminRequest,
) (*ssov1.IsAdminResponse, error) {
	if err := validateIsAdmin(req); err != nil {
		return nil, err
	}

	isAdmin, err := s.auth.IsAdmin(ctx, req.GetUserId())
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			return nil, status.Error(codes.NotFound, "user not found")
		}

		return nil, status.Error(codes.Internal, "failed to check admin status")
	}
	return &ssov1.IsAdminResponse{
		IsAdmin: isAdmin,
	}, nil
}

func validateLogin(req *ssov1.LoginRequest) error {
	if req.GetEmail() == "" && req.Password == "" {
		return status.Error(codes.InvalidArgument, "invalid email or password")
	}
	if req.GetEmail() == "" || req.GetEmail() == invalidEmail {
		return status.Error(codes.InvalidArgument, "email is required")
	}
	if _, err := mail.ParseAddress(req.GetEmail()); err != nil {
		return status.Error(codes.InvalidArgument, "this not email")
	}
	if req.GetPassword() == "" {
		return status.Error(codes.InvalidArgument, "password is required")
	}
	if req.GetAppId() == emptyValue {
		return status.Error(codes.InvalidArgument, "app_id is required")
	}
	return nil
}

func validateRegister(req *ssov1.RegisterRequest) error {
	if _, err := mail.ParseAddress(req.GetEmail()); err != nil {
		return status.Error(codes.InvalidArgument, "email is required")
	}
	if req.GetEmail() == "" || req.GetEmail() == invalidEmail {
		return status.Error(codes.InvalidArgument, "this not valid email")
	}

	if req.GetPassword() == "" {
		return status.Error(codes.InvalidArgument, "invalid password")
	}
	return nil
}

func validateIsAdmin(req *ssov1.IsAdminRequest) error {
	if req.GetUserId() == emptyValue {
		return status.Error(codes.InvalidArgument, "user_id is required")
	}
	return nil
}

func isDuplicateEmailError(err error) bool {
	return errors.Is(err, storage.ErrUserExists) || (strings.Contains(
		err.Error(),
		"duplicate key value",
	) && strings.Contains(err.Error(), "users_email_key"))
}
