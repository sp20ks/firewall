package delivery

import (
	"encoding/json"
	"net/http"
	"rules-engine/internal/entity"
	"rules-engine/internal/usecase"
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
		JSONResponse[any](w, http.StatusBadRequest, nil, err)
		return
	}

	if req.Name == "" || req.HTTPMethod == "" || req.URL == "" || req.Host == "" || req.CreatorID == "" {
		JSONResponse[any](w, http.StatusBadRequest, nil, errMissingFields())
		return
	}

	resource, err := h.resourceUseCase.Create(req.Name, req.HTTPMethod, req.URL, req.Host, req.CreatorID, req.IsActive)
	if err != nil {
		JSONResponse[any](w, http.StatusBadRequest, nil, err)
		return
	}

	JSONResponse(w, http.StatusOK, resource, nil)
}

func (h *ResourceHandler) HandleUpdateResource(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		JSONResponse[any](w, http.StatusBadRequest, nil, errMissingID())
		return
	}

	var req ResourceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		JSONResponse[any](w, http.StatusBadRequest, nil, err)
		return
	}

	if req.Name == "" && req.HTTPMethod == "" && req.URL == "" && req.Host == "" && req.IsActive == nil {
		JSONResponse[any](w, http.StatusBadRequest, nil, errMissingFields())
		return
	}

	resource, err := h.resourceUseCase.Update(id, req.Name, req.HTTPMethod, req.URL, req.Host, req.IsActive)
	if err != nil {
		JSONResponse[any](w, http.StatusBadRequest, nil, err)
		return
	}

	JSONResponse(w, http.StatusOK, resource, nil)
}

func (h *ResourceHandler) HandleGetActiveResources(w http.ResponseWriter, r *http.Request) {
	resources, err := h.resourceUseCase.Get()
	if err != nil {
		JSONResponse[any](w, http.StatusInternalServerError, nil, err)
		return
	}

	JSONResponse(w, http.StatusOK, ResourcesResponse{Resources: resources}, nil)
}

func (h *ResourceHandler) HandleAttachIPList(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		JSONResponse[any](w, http.StatusBadRequest, nil, errMissingID())
		return
	}

	var req UpdateIPListReferenceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		JSONResponse[any](w, http.StatusBadRequest, nil, err)
		return
	}

	if req.IPListID == "" {
		JSONResponse[any](w, http.StatusBadRequest, nil, errMissingID())
		return
	}

	err := h.resourceUseCase.AttachIPList(id, req.IPListID)
	if err != nil {
		JSONResponse[any](w, http.StatusInternalServerError, nil, err)
		return
	}

	JSONResponse[any](w, http.StatusOK, nil, nil)
}

func (h *ResourceHandler) HandleDetachIPList(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		JSONResponse[any](w, http.StatusBadRequest, nil, errMissingID())
		return
	}

	var req UpdateIPListReferenceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		JSONResponse[any](w, http.StatusBadRequest, nil, err)
		return
	}

	if req.IPListID == "" {
		JSONResponse[any](w, http.StatusBadRequest, nil, errMissingID())
		return
	}

	err := h.resourceUseCase.DetachIPList(id, req.IPListID)
	if err != nil {
		JSONResponse[any](w, http.StatusInternalServerError, nil, err)
		return
	}

	JSONResponse[any](w, http.StatusOK, nil, nil)
}

func (h *ResourceHandler) HandleAttachRule(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		JSONResponse[any](w, http.StatusBadRequest, nil, errMissingID())
		return
	}

	var req UpdateRuleReferenceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		JSONResponse[any](w, http.StatusBadRequest, nil, err)
		return
	}

	if req.RuleID == "" {
		JSONResponse[any](w, http.StatusBadRequest, nil, errMissingID())
		return
	}

	err := h.resourceUseCase.AttachRule(id, req.RuleID)
	if err != nil {
		JSONResponse[any](w, http.StatusInternalServerError, nil, err)
		return
	}

	JSONResponse[any](w, http.StatusOK, nil, nil)
}

func (h *ResourceHandler) HandleDetachRule(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		JSONResponse[any](w, http.StatusBadRequest, nil, errMissingID())
		return
	}

	var req UpdateRuleReferenceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		JSONResponse[any](w, http.StatusBadRequest, nil, err)
		return
	}

	if req.RuleID == "" {
		JSONResponse[any](w, http.StatusBadRequest, nil, errMissingID())
		return
	}

	err := h.resourceUseCase.DetachRule(id, req.RuleID)
	if err != nil {
		JSONResponse[any](w, http.StatusInternalServerError, nil, err)
		return
	}

	JSONResponse[any](w, http.StatusOK, nil, nil)
}
