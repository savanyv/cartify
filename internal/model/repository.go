package model

import (
	"context"
)

type UserRepository interface {
	Create(ctx context.Context, user *User) error
	FindByEmail(ctx context.Context, email string) (*User, error)
	FindByUsername(ctx context.Context, username string) (*User, error)
	FindByID(ctx context.Context, ID string) (*User, error)
	UpdateTokenVersion(ctx context.Context, ID string, version int) error
	Update(ctx context.Context, user *User) error
}
