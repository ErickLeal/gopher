package main

import (
	"context"
	"net/http"
	"strconv"

	"github.com/ErickLeal/gopher/internal/store"
	"github.com/go-chi/chi/v5"
)

type postCtxKey string

const postCtx postCtxKey = "post"

type CreatePostRequest struct {
	Title   string   `json:"title" validate:"required,max=100"`
	Content string   `json:"content" validate:"required,max=1000"`
	Tags    []string `json:"tags" validate:"required"`
}

type UpdatePostRequest struct {
	Title   *string `json:"title" validate:"omitempty,max=100"`
	Content *string `json:"content" validate:"omitempty,max=1000"`
}

func (app *application) createPostHandler(w http.ResponseWriter, r *http.Request) {
	var postRequest CreatePostRequest
	if err := readJson(w, r, &postRequest); err != nil {
		app.writeBadRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(postRequest); err != nil {
		app.writeBadRequestResponse(w, r, err)
		return
	}

	post := &store.PostModel{
		Title:   postRequest.Title,
		Content: postRequest.Content,
		Tags:    postRequest.Tags,
		UserID:  1,
	}

	ctx := r.Context()

	if err := app.store.Posts.Create(ctx, post); err != nil {
		app.writeInternalServerErrorResponse(w, r, err)
		return
	}

	if err := writeJson(w, http.StatusCreated, post); err != nil {
		app.writeInternalServerErrorResponse(w, r, err)
		return
	}
}

func (app *application) getPostHandler(w http.ResponseWriter, r *http.Request) {
	post := getPostFromCtx(r)

	comments, err := app.store.Comments.GetByPostID(r.Context(), post.ID)
	if err != nil {
		app.writeInternalServerErrorResponse(w, r, err)
		return
	}

	post.Comments = comments

	if err := writeJson(w, http.StatusOK, post); err != nil {
		app.writeInternalServerErrorResponse(w, r, err)
		return
	}
}

func (app *application) deletePostHandler(w http.ResponseWriter, r *http.Request) {
	postId, err := strconv.ParseInt(chi.URLParam(r, "postID"), 10, 64)
	if err != nil {
		app.writeInternalServerErrorResponse(w, r, err)
		return
	}

	err = app.store.Posts.Delete(r.Context(), postId)
	if err != nil {
		switch err {
		case store.ErrResourceNotFound:
			app.writeBadRequestResponse(w, r, err)
		default:
			app.writeInternalServerErrorResponse(w, r, err)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (app *application) updatePostHandler(w http.ResponseWriter, r *http.Request) {
	post := getPostFromCtx(r)

	var payload UpdatePostRequest
	if err := readJson(w, r, &payload); err != nil {
		app.writeBadRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.writeBadRequestResponse(w, r, err)
		return
	}

	if payload.Title != nil {
		post.Title = *payload.Title
	}
	if payload.Content != nil {
		post.Content = *payload.Content
	}

	err := app.store.Posts.Update(r.Context(), post)
	if err != nil {
		switch err {
		case store.ErrResourceNotFound:
			app.writeBadRequestResponse(w, r, err)
		default:
			app.writeInternalServerErrorResponse(w, r, err)
		}
		return
	}

	if err := writeJson(w, http.StatusCreated, post); err != nil {
		app.writeInternalServerErrorResponse(w, r, err)
		return
	}
}

func (app *application) postsContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		postId, err := strconv.ParseInt(chi.URLParam(r, "postID"), 10, 64)

		if err != nil {
			app.writeInternalServerErrorResponse(w, r, err)
			return
		}

		ctx := r.Context()

		post, err := app.store.Posts.GetById(r.Context(), postId)
		if err != nil {
			switch err {
			case store.ErrResourceNotFound:
				app.writeBadRequestResponse(w, r, err)
			default:
				app.writeInternalServerErrorResponse(w, r, err)
			}
			return
		}

		ctx = context.WithValue(ctx, postCtx, post)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getPostFromCtx(r *http.Request) *store.PostModel {
	post, _ := r.Context().Value(postCtx).(*store.PostModel)
	return post
}
