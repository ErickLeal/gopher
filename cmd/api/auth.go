package main

import (
	"fmt"
	"net/http"

	"github.com/ErickLeal/gopher/internal/env"
	"github.com/ErickLeal/gopher/internal/mailer"
	"github.com/ErickLeal/gopher/internal/store"
	"github.com/google/uuid"
)

type RegisterUserPayload struct {
	Username string `json:"username" validate:"required,max=100"`
	Email    string `json:"email" validate:"required,email,max=255"`
	Password string `json:"password" validate:"required,min=3,max=72"`
}

type UserWithToken struct {
	*store.UserModel
	Token string `json:"token"`
}

// registerUserHandler godoc
//
//	@Summary		Registers a user
//	@Description	Registers a user
//	@Tags			authentication
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		RegisterUserPayload	true	"User credentials"
//	@Success		201		{object}	UserWithToken		"User registered"
//	@Failure		400		{object}	error
//	@Failure		500		{object}	error
//	@Router			/authentication/user [post]
func (app *application) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	var payload RegisterUserPayload
	if err := readJson(w, r, &payload); err != nil {
		app.writeBadRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.writeBadRequestResponse(w, r, err)
		return
	}

	user := &store.UserModel{
		Username: payload.Username,
		Email:    payload.Email,
	}

	// hash the user password
	if err := user.Password.Set(payload.Password); err != nil {
		app.writeInternalServerErrorResponse(w, r, err)
		return
	}

	ctx := r.Context()

	token := uuid.New().String()

	err := app.store.Users.CreateAndInvite(ctx, user, token, app.config.mail.exp)
	if err != nil {
		switch err {
		case store.ErrDuplicateEmail:
			app.writeBadRequestResponse(w, r, err)
		case store.ErrDuplicateUsername:
			app.writeBadRequestResponse(w, r, err)
		default:
			app.writeInternalServerErrorResponse(w, r, err)
		}
		return
	}

	userWithToken := UserWithToken{
		UserModel: user,
		Token:     token,
	}

	activationURL := fmt.Sprintf("%s/confirm/%s", env.API_URL, token)
	isProduction := env.ENVIRONMENT == "production"
	vars := struct {
		Username      string
		ActivationURL string
	}{
		Username:      user.Username,
		ActivationURL: activationURL,
	}

	app.logger.Infow("sending user invitation email", "email", user.Email, "activationURL", activationURL, "isProduction", isProduction)
	err = app.mailer.Send(mailer.UserInvitationTemplate, user.Username, user.Email, vars, !isProduction)
	if err != nil {
		app.logger.Errorw("failed to send user invitation email", "error", err)

		if err := app.store.Users.Delete(ctx, user.ID); err != nil {
			app.logger.Errorw("failed to delete user after email send failure", "error", err)
		}

		app.writeInternalServerErrorResponse(w, r, err)
		return
	}

	if err := app.writeDataRespose(w, http.StatusCreated, userWithToken); err != nil {
		app.writeInternalServerErrorResponse(w, r, err)
	}
}
