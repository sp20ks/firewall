package delivery

import (
	"encoding/json"
	"net/http"
	"rules-engine/internal/entity"
	"rules-engine/internal/logger"
	"rules-engine/internal/usecase"

	"go.uber.org/zap"
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
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Name == "" || req.CreatorID == "" || req.AttackType == "" || req.ActionType == "" {
		http.Error(w, "All fields (name, creator_id, attack_type, action_type) must be provided", http.StatusBadRequest)
		return
	}

	err := h.ruleUseCase.Create(req.Name, req.AttackType, req.ActionType, req.CreatorID, req.IsActive)
	if err != nil {
		logger.Logger().Info("failed to create rule", zap.Error(err))
		http.Error(w, "Error while creating rule", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Rule is created"))
}

func (h *RuleHandler) HandleUpdateRule(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "id must be provided", http.StatusBadRequest)
		return
	}

	var req RuleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Name == "" && req.AttackType == "" && req.ActionType == "" && req.IsActive == nil {
		http.Error(w, "At least one field must be provided for update", http.StatusBadRequest)
		return
	}

	err := h.ruleUseCase.Update(id, req.Name, req.AttackType, req.ActionType, req.IsActive)
	if err != nil {
		logger.Logger().Info("failed to update rule", zap.Error(err))
		http.Error(w, "Error while updating rule", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Rule is updated"))
}

func (h *RuleHandler) HandleGetRules(w http.ResponseWriter, r *http.Request) {
	rules, err := h.ruleUseCase.Get()
	if err != nil {
		logger.Logger().Info("failed to get rules", zap.Error(err))
		http.Error(w, "Error while fetching rules", http.StatusInternalServerError)
		return
	}

	if len(rules) == 0 {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("No rules found"))
		return
	}

	response := RuleResponse{Rules: rules}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		logger.Logger().Info("failed to encode rules", zap.Error(err))
		http.Error(w, "Error while encoding rules", http.StatusInternalServerError)
		return
	}
}
