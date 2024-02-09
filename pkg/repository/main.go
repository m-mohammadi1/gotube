package repository

import (
	"context"
	"database/sql"
	"errors"
	"gotube/pkg/model"
	"gotube/pkg/repository/postgres/channelrepo"
	"gotube/pkg/repository/postgres/userrepo"
)

var ErrNotFound error = errors.New("model not found")

type Repository struct {
	UserRepository    UserRepository
	ChannelRepository ChannelRepository
}

func New(db *sql.DB) Repository {
	return Repository{
		UserRepository:    userrepo.New(db),
		ChannelRepository: channelrepo.New(db),
	}
}

type UserRepository interface {
	All(ctx context.Context) ([]*model.User, error)
	Find(ctx context.Context, id int) (*model.User, error)
	Delete(ctx context.Context, id int) error
	Create(ctx context.Context, user model.User) error

	EmailExists(email string) (bool, error)
	FindByField(ctx context.Context, field string, value string) (*model.User, error)
}

type ChannelRepository interface {
	UpdateOrCreate(ctx context.Context, channel model.Channel) (*model.Channel, error)
}
