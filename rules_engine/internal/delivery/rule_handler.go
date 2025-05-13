package delivery

import (
	"encoding/json"
	"net/http"
	"rules-engine/internal/delivery/middleware"
	"rules-engine/internal/entity"
	"rules-engine/internal/usecase"
)

type RuleHandler struct {
	ruleUseCase *usecase.RuleUseCase
}

func NewRuleHandler(ruleUseCase *usecase.RuleUseCase) *RuleHandler {
	return &RuleHandler{ruleUseCase: ruleUseCase}
}

type RuleRequest struct {
	Name       string `json:"name"`
	AttackType string `json:"attack_type"`
	ActionType string `json:"action_type"`
	CreatorID  string `json:"creator_id"`
	IsActive   *bool  `json:"is_active"`
}

type RuleResponse struct {
	Rules []entity.Rule `json:"rules"`
}

func (h *RuleHandler) HandleCreateRule(w http.ResponseWriter, r *http.Request) {
	var req RuleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		JSONResponse[any](w, http.StatusBadRequest, nil, err)
		return
	}

	if user, ok := middleware.GetUserFromContext(r.Context()); ok {
		req.CreatorID = user.ID
	}

	if req.Name == "" || req.CreatorID == "" || req.AttackType == "" || req.ActionType == "" {
		JSONResponse[any](w, http.StatusBadRequest, nil, errMissingFields())
		return
	}

	rule, err := h.ruleUseCase.Create(req.Name, req.AttackType, req.ActionType, req.CreatorID, req.IsActive)
	if err != nil {
		JSONResponse[any](w, http.StatusBadRequest, nil, err)
		return
	}

	JSONResponse(w, http.StatusOK, rule, nil)
}

func (h *RuleHandler) HandleUpdateRule(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		JSONResponse[any](w, http.StatusBadRequest, nil, errMissingID())
		return
	}

	var req RuleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		JSONResponse[any](w, http.StatusBadRequest, nil, err)
		return
	}

	if req.Name == "" && req.AttackType == "" && req.ActionType == "" && req.IsActive == nil {
		JSONResponse[any](w, http.StatusBadRequest, nil, errMissingFields())
		return
	}

	rule, err := h.ruleUseCase.Update(id, req.Name, req.AttackType, req.ActionType, req.IsActive)
	if err != nil {
		JSONResponse[any](w, http.StatusBadRequest, nil, err)
		return
	}

	JSONResponse(w, http.StatusOK, rule, nil)
}

func (h *RuleHandler) HandleGetRules(w http.ResponseWriter, r *http.Request) {
	rules, err := h.ruleUseCase.Get()
	if err != nil {
		JSONResponse[any](w, http.StatusInternalServerError, nil, err)
		return
	}

	JSONResponse(w, http.StatusOK, RuleResponse{Rules: rules}, nil)
}
