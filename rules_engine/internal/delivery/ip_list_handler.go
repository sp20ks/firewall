package delivery

import (
	"encoding/json"
	"net/http"
	"rules-engine/internal/entity"
	"rules-engine/internal/logger"
	"rules-engine/internal/usecase"

	"go.uber.org/zap"
)

type IPListHandler struct {
	ipListUseCase *usecase.IPListUseCase
}

func NewIPListHandler(ipListUseCase *usecase.IPListUseCase) *IPListHandler {
	return &IPListHandler{ipListUseCase: ipListUseCase}
}

type IPListRequest struct {
	IP        string `json:"ip"`
	ListType  string `json:"list_type"`
	CreatorID string `json:"creator_id"`
}

type IPListResponse struct {
	IPLists []entity.IPList `json:"ip_lists"`
}

func (h *IPListHandler) HandleCreateIPList(w http.ResponseWriter, r *http.Request) {
	var req IPListRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.IP == "" || req.CreatorID == "" || req.ListType == "" {
		http.Error(w, "All fields (ip, creator_id, list_type) must be provided", http.StatusBadRequest)
		return
	}

	err := h.ipListUseCase.Create(req.IP, req.ListType, req.CreatorID)
	if err != nil {
		logger.Logger().Info("failed to create ip list", zap.Error(err))
		http.Error(w, "Error while creating ip list", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Ip list is created"))
}

func (h *IPListHandler) HandleUpdateIPList(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "id must be provided", http.StatusBadRequest)
		return
	}

	var req IPListRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.IP == "" && req.ListType == "" {
		http.Error(w, "At least one field must be provided for update", http.StatusBadRequest)
		return
	}

	err := h.ipListUseCase.Update(id, req.IP, req.ListType)
	if err != nil {
		logger.Logger().Info("failed to update ip list", zap.Error(err))
		http.Error(w, "Error while updating ip list", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Ip list is updated"))
}

func (h *IPListHandler) HandleGetIPLists(w http.ResponseWriter, r *http.Request) {
	lists, err := h.ipListUseCase.Get()
	if err != nil {
		logger.Logger().Info("failed to get ip lists", zap.Error(err))
		http.Error(w, "Error while fetching ip lists", http.StatusInternalServerError)
		return
	}

	if len(lists) == 0 {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("No ip lists found"))
		return
	}

	response := IPListResponse{IPLists: lists}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		logger.Logger().Info("failed to encode ip lists", zap.Error(err))
		http.Error(w, "Error while encoding ip lists", http.StatusInternalServerError)
		return
	}
}
