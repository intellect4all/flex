package authentication

import "context"

type UserID string

type AuthenticatedUserJWT string

type User struct {
	Id   UserID `json:"email" bson:"_id" validate:"required,email" `
	Role string `json:"role" bson:"role" validate:"required, oneof=customer admin"`
}

type UserCredential struct {
	User
	Password string `json:"password" bson:"password" validate:"required"`
}

type UserDetail struct {
	User
	Phone string `json:"phone"bson:"phone" validate:"required"`
}

func (u *UserDetail) isAdmin() bool {
	return u.Role == "admin"
}

func (u *UserDetail) isCustomer() bool {
	return u.Role == "customer"
}

type Service struct {
	repository Repository
	jwtHelper  JWTHelper
}

func (s *Service) CreateUser(ctx context.Context, userDetail *UserCredential) error {
	return s.repository.CreateNewUser(ctx, userDetail)
}

func (s *Service) AuthenticateUser(ctx context.Context, userCredential *UserCredential) (*AuthenticatedUserJWT, error) {
	userCredentialFromDb, err := s.repository.GetUserCredential(ctx, userCredential.Id)
	if err != nil {
		return nil, err
	}

	if !s.repository.isValidPassword(ctx, userCredential.Password, userCredentialFromDb.Password) {
		return nil, ErrInvalidCredentials
	}

	token, err := s.jwtHelper.GenerateJWT(userCredentialFromDb.Id)
	if err != nil {
		return nil, err
	}

	return token, nil
}

func (s *Service) refreshJWT(ctx context.Context, jwt AuthenticatedUserJWT) (token *AuthenticatedUserJWT, err error) {
	token, err = s.jwtHelper.RenewJWT(string(jwt))
	return
}

func (s *Service) DeleteUser(ctx context.Context, jwt AuthenticatedUserJWT, userID UserID) error {
	userDetail, err := s.jwtHelper.ValidateJWT(jwt)
	if err != nil {
		return err
	}

	if !userDetail.isAdmin() {
		return ErrUnauthorized
	}

	userToDelete, err := s.repository.GetUserDetail(ctx, userID)
	if err != nil {
		return err
	}

	if userToDelete.isAdmin() {
		return ErrUnauthorized
	}

	return s.repository.DeleteUser(ctx, userID)
}

func NewService(repository Repository, jwtHelper JWTHelper) *Service {
	return &Service{repository: repository, jwtHelper: jwtHelper}
}

type Repository interface {
	CreateNewUser(ctx context.Context, userDetail *UserCredential) error
	GetUserCredential(ctx context.Context, email UserID) (*UserCredential, error)
	GetUserDetail(ctx context.Context, email UserID) (*UserDetail, error)
	DeleteUser(ctx context.Context, email UserID) error
	isValidPassword(ctx context.Context, passwordFromRequest, passwordFromDb string) bool
}

type JWTHelper interface {
	GenerateJWT(id UserID) (*AuthenticatedUserJWT, error)
	ValidateJWT(jwt AuthenticatedUserJWT) (UserDetail, error)
	RenewJWT(jwt string) (*AuthenticatedUserJWT, error)
}
