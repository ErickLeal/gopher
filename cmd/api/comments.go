package main

import (
	"net/http"

	"github.com/ErickLeal/gopher/internal/store"
)

type CreateCommentRequest struct {
	Content string `json:"content" validate:"required,max=1000"`
}

func (app *application) createCommentHandler(w http.ResponseWriter, r *http.Request) {
	var commentRequest CreateCommentRequest
	if err := readJson(w, r, &commentRequest); err != nil {
		app.writeBadRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(commentRequest); err != nil {
		app.writeBadRequestResponse(w, r, err)
		return
	}

	post := getPostFromCtx(r)

	comment := &store.CommentModel{
		PostID:  post.ID,
		UserID:  post.UserID,
		Content: commentRequest.Content,
	}

	if err := app.store.Comments.Create(r.Context(), comment); err != nil {
		app.writeInternalServerErrorResponse(w, r, err)
		return
	}

	if err := app.writeDataRespose(w, http.StatusCreated, comment); err != nil {
		app.writeInternalServerErrorResponse(w, r, err)
		return
	}
}

func (app *application) getCommentsHandler(w http.ResponseWriter, r *http.Request) {
	post := getPostFromCtx(r)

	comments, err := app.store.Comments.GetByPostID(r.Context(), post.ID)
	if err != nil {
		app.writeInternalServerErrorResponse(w, r, err)
		return
	}

	if err := app.writeDataRespose(w, http.StatusOK, comments); err != nil {
		app.writeInternalServerErrorResponse(w, r, err)
		return
	}
}
