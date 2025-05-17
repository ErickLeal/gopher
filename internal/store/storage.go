package store

import (
	"context"
	"database/sql"
	"errors"
)

var (
	ErrResourceNotFound = errors.New("resource not found")
)

type Storage struct {
	Posts interface {
		Create(context.Context, *PostModel) error
		GetById(context.Context, int64) (*PostModel, error)
	}
	Users interface {
		Create(context.Context, *UserModel) error
	}
	Comments interface {
		GetByPostID(context.Context, int64) ([]CommentModel, error)
	}
}

func NewStorage(db *sql.DB) Storage {
	return Storage{
		Posts:    &PostsStore{db: db},
		Users:    &UserStore{db: db},
		Comments: &CommentStore{db: db},
	}
}
