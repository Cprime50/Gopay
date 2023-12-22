package model

import (
	"context"

	"github.com/google/uuid"
)

type UserService interface {
	DeleteImage(ctx context.Context, uid uuid.UUID) error
	Get(ctx context.Context, uid uuid.UUID) (*User, error)
	Signup(ctx context.Context, u *User) error
	Signin(ctx context.Context, u *User) error
	UpdateDetails(ctx context.Context)
}
