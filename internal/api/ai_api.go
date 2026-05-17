package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/SaenkoDmitry/training-tg-bot/internal/application/dto"
	"github.com/SaenkoDmitry/training-tg-bot/internal/middlewares"
)

func (s *serviceImpl) GetAIProgramContext(w http.ResponseWriter, r *http.Request) {
	claims, ok := middlewares.FromContext(r.Context())
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	req := dto.AIProgramPromptRequest{}
	if rawProgramID := r.URL.Query().Get("program_id"); rawProgramID != "" {
		programID, err := strconv.ParseInt(rawProgramID, 10, 64)
		if err != nil {
			http.Error(w, "bad program_id", http.StatusBadRequest)
			return
		}
		req.ProgramID = &programID
	}

	ctx, err := s.container.BuildProgramPromptUC.BuildContext(claims.UserID, req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ctx)
}

func (s *serviceImpl) BuildAIProgramPrompt(w http.ResponseWriter, r *http.Request) {
	claims, ok := middlewares.FromContext(r.Context())
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	var req dto.AIProgramPromptRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	result, err := s.container.BuildProgramPromptUC.Execute(claims.UserID, req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func (s *serviceImpl) CreateProgramFromAI(w http.ResponseWriter, r *http.Request) {
	claims, ok := middlewares.FromContext(r.Context())
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	var req dto.AIApplyProgramRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	result, err := s.container.ApplyAIProgramUC.Execute(claims.UserID, req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}
