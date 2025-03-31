package delivery

import (
	"encoding/json"
	"log"
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
		log.Printf("Failed to create resource: %v", err)
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
		log.Printf("Failed to update resource: %v", err)
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
		log.Printf("Error encoding resources: %v", err)
		http.Error(w, "Error while encoding resources", http.StatusInternalServerError)
		return
	}
}
