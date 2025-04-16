package delivery

import (
	"encoding/json"
	"log"
	"net/http"
	"rules-engine/internal/entity"
	"rules-engine/internal/usecase"
)

type AnalizerHandler struct {
	analizer *usecase.AnalizerUseCase
}

func NewAnalizerHandler(analizer *usecase.AnalizerUseCase) *AnalizerHandler {
	return &AnalizerHandler{analizer: analizer}
}

func (h *AnalizerHandler) HandleAnalizeRequest(w http.ResponseWriter, r *http.Request) {
	var req entity.Request

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	result, err := h.analizer.AnalyzeRequest(&req)
	if err != nil {
		log.Printf("Error analyzing request: %v", err)
		http.Error(w, "Error analyzing request", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(result); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Error while encoding response", http.StatusInternalServerError)
	}
}
