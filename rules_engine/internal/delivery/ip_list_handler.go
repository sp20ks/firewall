package delivery

import (
	"encoding/json"
	"net/http"
	"rules-engine/internal/entity"
	"rules-engine/internal/usecase"
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
		JSONResponse[any](w, http.StatusBadRequest, nil, err)
		return
	}

	if req.IP == "" || req.CreatorID == "" || req.ListType == "" {
		JSONResponse[any](w, http.StatusBadRequest, nil, errMissingFields())
		return
	}

	ipList, err := h.ipListUseCase.Create(req.IP, req.ListType, req.CreatorID)
	if err != nil {
		JSONResponse[any](w, http.StatusBadRequest, nil, err)
		return
	}

	JSONResponse(w, http.StatusOK, ipList, nil)
}

func (h *IPListHandler) HandleUpdateIPList(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		JSONResponse[any](w, http.StatusBadRequest, nil, errMissingID())
		return
	}

	var req IPListRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		JSONResponse[any](w, http.StatusBadRequest, nil, err)
		return
	}

	if req.IP == "" && req.ListType == "" {
		JSONResponse[any](w, http.StatusBadRequest, nil, errMissingFields())
		return
	}

	ipList, err := h.ipListUseCase.Update(id, req.IP, req.ListType)
	if err != nil {
		JSONResponse[any](w, http.StatusBadRequest, nil, err)
		return
	}

	JSONResponse(w, http.StatusOK, ipList, nil)
}

func (h *IPListHandler) HandleGetIPLists(w http.ResponseWriter, r *http.Request) {
	lists, err := h.ipListUseCase.Get()
	if err != nil {
		JSONResponse[any](w, http.StatusInternalServerError, nil, err)
		return
	}

	JSONResponse(w, http.StatusOK, IPListResponse{IPLists: lists}, nil)
}
