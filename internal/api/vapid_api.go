package api

import (
	"encoding/json"
	"net/http"
	"os"
)

func (s *serviceImpl) GetVapidKey(w http.ResponseWriter, r *http.Request) {
	key := os.Getenv("VAPID_PUBLIC_KEY")
	json.NewEncoder(w).Encode(map[string]string{
		"public_key": key,
	})
}
