package projectclient

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/gorilla/mux"
)

type Handler struct {
	client *Client
	logger *slog.Logger
}

func NewHandler(client *Client, logger *slog.Logger) *Handler {
	return &Handler{
		client: client,
		logger: logger,
	}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/api/projects", h.GetAllProjects).Methods("GET")
}

func (h *Handler) GetAllProjects(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("fetching all projects from project-service")

	projects, err := h.client.GetAllProjects(r.Context())
	if err != nil {
		h.logger.Error("failed to fetch projects", "error", err)
		respondWithError(w, http.StatusInternalServerError, "Failed to fetch projects")
		return
	}

	respondWithJSON(w, http.StatusOK, projects)
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
