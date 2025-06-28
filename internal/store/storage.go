package store

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

var (
	ErrResourceNotFound  = errors.New("resource not found")
	ErrConflict          = errors.New("resource already exists")
	QueryTimeoutDuration = time.Second * 5
)

type Storage struct {
	Posts interface {
		Create(context.Context, *PostModel) error
		GetById(context.Context, int64) (*PostModel, error)
		Delete(context.Context, int64) error
		Update(context.Context, *PostModel) error
		GetUserFeed(context.Context, int64, PaginatedFeedQuery) ([]PostWithMetadata, error)
	}
	Users interface {
		Create(context.Context, *UserModel) error
		GetByID(context.Context, int64) (*UserModel, error)
	}
	Comments interface {
		Create(context.Context, *CommentModel) error
		GetByPostID(context.Context, int64) ([]CommentModel, error)
	}
	Followers interface {
		Follow(context.Context, int64, int64) error
		Unfollow(context.Context, int64, int64) error
	}
}

func NewStorage(db *sql.DB) Storage {
	return Storage{
		Posts:     &PostsStore{db: db},
		Users:     &UserStore{db: db},
		Comments:  &CommentStore{db: db},
		Followers: &FollowerStore{db: db},
	}
}
