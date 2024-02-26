package handlers

import (
	"net/http"

	"github.com/GeovaneCavalcante/rate-limit-api/pkg/logger"
)

type HealthHandler struct {
}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

func (h *HealthHandler) HealthHandler(w http.ResponseWriter, r *http.Request) {
	logger.Info("[HealthHandler] starting handler")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status": "ok"}`))
}
