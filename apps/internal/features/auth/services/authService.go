package authService

import (
	"context"
	"errors"

	"go-boilerplate/ent"
	"go-boilerplate/ent/user"

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	client *ent.Client
}

func NewAuthService(client *ent.Client) *AuthService {
	return &AuthService{client: client}
}

func (s *AuthService) Login(
	ctx context.Context,
	email string,
	password string,
) (*ent.User, error) {

	u, err := s.client.User.
		Query().
		Where(user.EmailEQ(email)).
		Only(ctx)

	if err != nil {
		return nil, errors.New("email atau password salah")
	}

	if err := bcrypt.CompareHashAndPassword(
		[]byte(u.Password),
		[]byte(password),
	); err != nil {
		return nil, errors.New("email atau password salah")
	}

	return u, nil
}


func (s *AuthService) Register(
	ctx context.Context,
	name string,
	email string,
	password string,
) error {

	exists, err := s.client.User.
		Query().
		Where(user.EmailEQ(email)).
		Exist(ctx)

	if err != nil {
		return err
	}

	if exists {
		return errors.New("email sudah terdaftar")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return err
	}

	_, err = s.client.User.
		Create().
		SetName(name).
		SetEmail(email).
		SetPassword(string(hashedPassword)).
		Save(ctx)

	return err
}