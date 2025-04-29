package delivery

import (
	"encoding/json"
	"net/http"
	"rules-engine/internal/entity"
	"rules-engine/internal/usecase"
)

type AnalyzerHandler struct {
	analyzer *usecase.AnalyzerUseCase
}

func NewAnalyzerHandler(analyzer *usecase.AnalyzerUseCase) *AnalyzerHandler {
	return &AnalyzerHandler{analyzer: analyzer}
}

func (h *AnalyzerHandler) HandleAnalyzeRequest(w http.ResponseWriter, r *http.Request) {
	var req entity.Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		JSONResponse[any](w, http.StatusBadRequest, nil, errMissingFields())
		return
	}

	result, err := h.analyzer.AnalyzeRequest(&req)
	if err != nil {
		JSONResponse[any](w, http.StatusInternalServerError, nil, err)
		return
	}

	JSONResponse[any](w, http.StatusOK, result, nil)
}
