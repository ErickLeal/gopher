package main

import (
	"context"
	"net/http"
	"strconv"

	"github.com/ErickLeal/gopher/internal/store"
	"github.com/go-chi/chi/v5"
)

type userKey string

const userCtx userKey = "user"

type TargetUserRequest struct {
	TargetUserID int64 `json:"target_user_id" validate:"required"`
}

// GetUser godoc
//
//	@Summary		Fetches a user profile
//	@Description	Fetches a user profile by ID
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"User ID"
//	@Success		200	{object}	store.UserModel
//	@Failure		400	{object}	error
//	@Failure		404	{object}	error
//	@Failure		500	{object}	error
//	@Security		ApiKeyAuth
//	@Router			/users/{id} [get]
func (app *application) getUserHandler(w http.ResponseWriter, r *http.Request) {
	user := getUserFromContext(r)

	if err := app.writeDataRespose(w, http.StatusOK, user); err != nil {
		app.writeInternalServerErrorResponse(w, r, err)
	}
}

// FollowUser godoc
//
//	@Summary		Follows a user
//	@Description	Follows a user by ID
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			userID	path		int					true	"User ID"
//	@Param			payload	body		TargetUserRequest	true	"Follow user payload"
//	@Success		204		{string}	string				"User followed"
//	@Failure		400		{object}	error				"User payload missing"
//	@Failure		404		{object}	error				"User not found"
//	@Security		ApiKeyAuth
//	@Router			/users/{userID}/follow [put]
func (app *application) followUserHandler(w http.ResponseWriter, r *http.Request) {
	user := getUserFromContext(r)

	var followRequest TargetUserRequest
	if err := readJson(w, r, &followRequest); err != nil {
		app.writeBadRequestResponse(w, r, err)
		return
	}

	err := app.store.Followers.Follow(r.Context(), user.ID, followRequest.TargetUserID)
	if err != nil {
		switch err {
		case store.ErrConflict:
			app.writeConlfictResponse(w, r, err)
		default:
			app.writeInternalServerErrorResponse(w, r, err)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (app *application) unfollowUserHandler(w http.ResponseWriter, r *http.Request) {
	user := getUserFromContext(r)

	var followRequest TargetUserRequest
	if err := readJson(w, r, &followRequest); err != nil {
		app.writeBadRequestResponse(w, r, err)
		return
	}

	err := app.store.Followers.Unfollow(r.Context(), user.ID, followRequest.TargetUserID)
	if err != nil {
		switch err {
		case store.ErrConflict:
			app.writeConlfictResponse(w, r, err)
		default:
			app.writeInternalServerErrorResponse(w, r, err)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func getUserFromContext(r *http.Request) *store.UserModel {
	user, _ := r.Context().Value(userCtx).(*store.UserModel)
	return user
}

func (app *application) usersContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, err := strconv.ParseInt(chi.URLParam(r, "userID"), 10, 64)

		if err != nil {
			app.writeInternalServerErrorResponse(w, r, err)
			return
		}

		ctx := r.Context()

		user, err := app.store.Users.GetByID(r.Context(), userID)
		if err != nil {
			switch err {
			case store.ErrResourceNotFound:
				app.writeBadRequestResponse(w, r, err)
				return
			default:
				app.writeInternalServerErrorResponse(w, r, err)
				return
			}
		}

		ctx = context.WithValue(ctx, userCtx, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
