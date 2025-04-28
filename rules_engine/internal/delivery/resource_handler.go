package delivery

import (
	"encoding/json"
	"net/http"
	"rules-engine/internal/entity"
	"rules-engine/internal/logger"
	"rules-engine/internal/usecase"

	"go.uber.org/zap"
)

type ResourceHandler struct {
	resourceUseCase *usecase.ResourceUseCase
}

func NewResourceHandler(resourceUseCase *usecase.ResourceUseCase) *ResourceHandler {
	return &ResourceHandler{resourceUseCase: resourceUseCase}
}

type ResourceRequest struct {
	Name       string `json:"name"`
	HTTPMethod string `json:"http_method"`
	URL        string `json:"url"`
	Host       string `json:"host"`
	CreatorID  string `json:"creator_id"`
	IsActive   *bool  `json:"is_active"`
}

type UpdateIPListReferenceRequest struct {
	IPListID string `json:"ip_list_id"`
}

type UpdateRuleReferenceRequest struct {
	RuleID string `json:"rule_id"`
}

type ResourcesResponse struct {
	Resources []entity.Resource `json:"resources"`
}

func (h *ResourceHandler) HandleCreateResource(w http.ResponseWriter, r *http.Request) {
	var req ResourceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Name == "" || req.HTTPMethod == "" || req.URL == "" || req.Host == "" || req.CreatorID == "" {
		http.Error(w, "All fields (name, http_method, url, creator_id) must be provided", http.StatusBadRequest)
		return
	}

	err := h.resourceUseCase.Create(req.Name, req.HTTPMethod, req.URL, req.Host, req.CreatorID, req.IsActive)
	if err != nil {
		logger.Logger().Info("failed to create resource", zap.Error(err))
		http.Error(w, "Error while creating resource", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Resource is created"))
}

func (h *ResourceHandler) HandleUpdateResource(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "id must be provided", http.StatusBadRequest)
		return
	}

	var req ResourceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Name == "" && req.HTTPMethod == "" && req.URL == "" && req.Host == "" && req.IsActive == nil {
		http.Error(w, "At least one field must be provided for update", http.StatusBadRequest)
		return
	}

	err := h.resourceUseCase.Update(id, req.Name, req.HTTPMethod, req.URL, req.Host, req.IsActive)
	if err != nil {
		logger.Logger().Info("failed to update resource", zap.Error(err))
		http.Error(w, "Error while updating resource", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Resource is updated"))
}

func (h *ResourceHandler) HandleGetActiveResources(w http.ResponseWriter, r *http.Request) {
	resources, err := h.resourceUseCase.Get()
	if err != nil {
		http.Error(w, "Error while fetching active resources", http.StatusInternalServerError)
		return
	}

	if len(resources) == 0 {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("No active resources found"))
		return
	}

	response := ResourcesResponse{Resources: resources}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		logger.Logger().Info("failed to encode resources", zap.Error(err))
		http.Error(w, "Error while encoding resources", http.StatusInternalServerError)
		return
	}
}

func (h *ResourceHandler) HandleAttachIPList(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "resource id must be provided", http.StatusBadRequest)
		return
	}

	var req UpdateIPListReferenceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.IPListID == "" {
		http.Error(w, "ip_list_id must be provided", http.StatusBadRequest)
		return
	}

	err := h.resourceUseCase.AttachIPList(id, req.IPListID)
	if err != nil {
		logger.Logger().Info("failed to attach IP list", zap.Error(err))
		http.Error(w, "Error while attaching IP list", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("IP list attached to resource"))
}

func (h *ResourceHandler) HandleDetachIPList(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "resource id must be provided", http.StatusBadRequest)
		return
	}

	var req UpdateIPListReferenceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.IPListID == "" {
		http.Error(w, "ip_list_id must be provided", http.StatusBadRequest)
		return
	}

	err := h.resourceUseCase.DetachIPList(id, req.IPListID)
	if err != nil {
		logger.Logger().Info("failed to detach IP list", zap.Error(err))
		http.Error(w, "Error while detaching IP list", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("IP list detached from resource"))
}

func (h *ResourceHandler) HandleAttachRule(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "resource id must be provided", http.StatusBadRequest)
		return
	}

	var req UpdateRuleReferenceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.RuleID == "" {
		http.Error(w, "rule_id must be provided", http.StatusBadRequest)
		return
	}

	err := h.resourceUseCase.AttachRule(id, req.RuleID)
	if err != nil {
		logger.Logger().Info("failed to attach rule", zap.Error(err))
		http.Error(w, "Error while attaching rule", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Rule attached to resource"))
}

func (h *ResourceHandler) HandleDetachRule(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "resource id must be provided", http.StatusBadRequest)
		return
	}

	var req UpdateRuleReferenceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.RuleID == "" {
		http.Error(w, "rule_id must be provided", http.StatusBadRequest)
		return
	}

	err := h.resourceUseCase.DetachRule(id, req.RuleID)
	if err != nil {
		logger.Logger().Info("failed to detach rule", zap.Error(err))
		http.Error(w, "Error while detaching rule", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Rule detached from resource"))
}
