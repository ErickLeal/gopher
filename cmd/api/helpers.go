package main

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
)

var Validate *validator.Validate

func init() {
	Validate = validator.New(validator.WithRequiredStructEnabled())
}

func writeJson(w http.ResponseWriter, status int, data any) error {

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(data)
}

func readJson(w http.ResponseWriter, r *http.Request, data any) error {
	maxBytes := 1_048_576 // 1 MB
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	return decoder.Decode(data)
}

func writeJsonError(w http.ResponseWriter, status int, message string) error {
	type envelope struct {
		Error string `json:"error"`
	}

	return writeJson(w, status, &envelope{Error: message})
}

func (app *application) writeInternalServerErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Errorw("internal error", "method", r.Method, "path", r.URL.Path, "error", err.Error())

	writeJsonError(w, http.StatusInternalServerError, "internal server error")
}

func (app *application) writeBadRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Warnln("bad request", "method", r.Method, "path", r.URL.Path, "- error -> ", err.Error())

	writeJsonError(w, http.StatusBadRequest, err.Error())
}

func (app *application) writeConlfictResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Errorf("conflict response", "method", r.Method, "path", r.URL.Path, "error", err.Error())

	writeJsonError(w, http.StatusConflict, err.Error())
}

func (app *application) writeNotFoudResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Warnf("not found error", "method", r.Method, "path", r.URL.Path, "error", err.Error())

	writeJsonError(w, http.StatusNotFound, "not found")
}

func (app *application) writeDataRespose(w http.ResponseWriter, status int, data any) error {
	type envelope struct {
		Data any `json:"data"`
	}

	return writeJson(w, status, &envelope{Data: data})
}
