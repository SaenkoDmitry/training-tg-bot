package api

import (
	"encoding/json"
	"net/http"

	"github.com/SaenkoDmitry/training-tg-bot/internal/middlewares"
)

func (s *serviceImpl) MeHandler(w http.ResponseWriter, r *http.Request) {
	claims, ok := middlewares.FromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", 401)
		return
	}

	userID := claims.UserID

	user, err := s.container.GetUserByIDUC.Execute(userID) // из БД
	if err != nil {
		http.Error(w, "not found", 404)
		return
	}

	resp := map[string]interface{}{
		"id":         user.ID,
		"first_name": user.FirstName,
		"last_name":  user.LastName,
	}

	json.NewEncoder(w).Encode(resp)
}
