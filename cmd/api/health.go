package main

import (
	"net/http"

	"github.com/ErickLeal/gopher/internal/env"
)

func (app *application) healthHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"status": "ok",
		"env":    env.ENVIRONMENT,
	}

	if err := writeJson(w, http.StatusOK, data); err != nil {
		app.writeInternalServerErrorResponse(w, r, err)
		return
	}
}
