package delivery

import (
	"encoding/json"
	"net/http"
	"rules-engine/internal/entity"
	"rules-engine/internal/logger"
	"rules-engine/internal/usecase"

	"go.uber.org/zap"
)

type AnalyzerHandler struct {
	analyzer *usecase.AnalyzerUseCase
}

func NewAnalyzerHandler(analyzer *usecase.AnalyzerUseCase) *AnalyzerHandler {
	return &AnalyzerHandler{analyzer: analyzer}
}

func (h *AnalyzerHandler) HandleAnalyzeRequest(w http.ResponseWriter, r *http.Request) {
	var req entity.Request
	l := logger.Logger()
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		l.Info("failed parse body", zap.Error(err))
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	result, err := h.analyzer.AnalyzeRequest(&req)
	if err != nil {
		l.Info("error analyzing request", zap.Error(err))
		http.Error(w, "Error analyzing request", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(result); err != nil {
		l.Info("error encoding response", zap.Error(err))
		http.Error(w, "Error while encoding response", http.StatusInternalServerError)
	}
}
