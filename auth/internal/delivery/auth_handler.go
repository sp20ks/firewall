package delivery

import (
	"encoding/json"
	"net/http"

	"auth/internal/logger"
	"auth/internal/usecase"

	"go.uber.org/zap"
)

type AuthHandler struct {
	authUseCase *usecase.AuthUseCase
}

func NewAuthHandler(authUseCase *usecase.AuthUseCase) *AuthHandler {
	return &AuthHandler{authUseCase: authUseCase}
}

type UserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type SuccessResponse struct {
	Message string `json:"message"`
}

type TokenResponse struct {
	Token string `json:"token"`
}

func sendJSONResponse(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func (h *AuthHandler) HandleRegisterUser(w http.ResponseWriter, r *http.Request) {
	var req UserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendJSONResponse(w, http.StatusBadRequest, ErrorResponse{"Invalid request body"})
		return
	}

	if req.Username == "" || req.Password == "" {
		sendJSONResponse(w, http.StatusBadRequest, ErrorResponse{"Both 'username' and 'password' must be provided"})
		return
	}

	err := h.authUseCase.CreateUser(req.Username, req.Password)
	if err != nil {
		// TODO: если юзер уже есть, то должна возвращаться читаемая ошибка, а не как сейчас.
		logger.Logger().Info("failed to create user", zap.Error(err))
		sendJSONResponse(w, http.StatusBadRequest, ErrorResponse{"Error while creating user"})
		return
	}

	sendJSONResponse(w, http.StatusOK, SuccessResponse{"User is registered"})
}

func (h *AuthHandler) HandleGetJwtToken(w http.ResponseWriter, r *http.Request) {
	var req UserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Logger().Info("failed to parse body", zap.Error(err))
		sendJSONResponse(w, http.StatusBadRequest, ErrorResponse{"Invalid request body"})
		return
	}

	if req.Username == "" || req.Password == "" {
		sendJSONResponse(w, http.StatusBadRequest, ErrorResponse{"Both 'username' and 'password' must be provided"})
		return
	}

	tokenString, err := h.authUseCase.Authenticate(req.Username, req.Password)
	if err != nil {
		logger.Logger().Info("failed to authenticate", zap.Error(err))
		sendJSONResponse(w, http.StatusUnauthorized, ErrorResponse{"Invalid credentials"})
		return
	}

	sendJSONResponse(w, http.StatusOK, TokenResponse{tokenString})
}

func (h *AuthHandler) VerifyJwtToken(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if token == "" {
		sendJSONResponse(w, http.StatusBadRequest, ErrorResponse{"Missing 'token' parameter"})
		return
	}

	user, err := h.authUseCase.GetUserByToken(token)
	if err != nil {
		logger.Logger().Info("failed to verify token", zap.Error(err))
		sendJSONResponse(w, http.StatusUnauthorized, ErrorResponse{"Invalid token"})
		return
	}

	sendJSONResponse(w, http.StatusOK, user)
}
