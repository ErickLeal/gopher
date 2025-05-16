package main

import (
	"net/http"
	"strconv"

	"github.com/ErickLeal/gopher/internal/store"
	"github.com/go-chi/chi/v5"
)

type CreatePostRequest struct {
	Title   string   `json:"title" validate:"required,max=100"`
	Content string   `json:"content" validate:"required,max=1000"`
	Tags    []string `json:"tags" validate:"required"`
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
	postId, err := strconv.ParseInt(chi.URLParam(r, "postID"), 10, 64)
	if err != nil {
		app.writeInternalServerErrorResponse(w, r, err)
		return
	}

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

	if err := writeJson(w, http.StatusOK, post); err != nil {
		app.writeInternalServerErrorResponse(w, r, err)
		return
	}
}
