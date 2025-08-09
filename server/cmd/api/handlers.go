package main

import (
	"net/http"

	"github.com/Babatunde50/byfood-assessment/server/internal/response"
)

func (app *application) status(w http.ResponseWriter, r *http.Request) {
	_ = response.JSON(w, http.StatusOK, map[string]string{"Status": "OK"})
}
