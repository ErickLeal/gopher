package store

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

var (
	ErrResourceNotFound  = errors.New("resource not found")
	QueryTimeoutDuration = time.Second * 5
)

type Storage struct {
	Posts interface {
		Create(context.Context, *PostModel) error
		GetById(context.Context, int64) (*PostModel, error)
		Delete(context.Context, int64) error
		Update(context.Context, *PostModel) error
	}
	Users interface {
		Create(context.Context, *UserModel) error
	}
	Comments interface {
		Create(context.Context, *CommentModel) error
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
