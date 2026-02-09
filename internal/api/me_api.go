package api

import (
	"encoding/json"
	"github.com/golang-jwt/jwt/v4"
	"net/http"
)

func MeHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session")
	if err != nil {
		http.Error(w, "unauthorized", 401)
		return
	}

	token, err := jwt.Parse(cookie.Value, func(t *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil || !token.Valid {
		http.Error(w, "unauthorized", 401)
		return
	}

	claims := token.Claims.(jwt.MapClaims)

	//if claims["photo_url"] != "" {
	//resImg, err := http.Get(fmt.Sprintf("https://api.telegram.org/bot%s/getUserProfilePhotos?user_id=%s&limit=1", botToken, claims["id"]))
	//}

	resp := map[string]interface{}{
		"id":         claims["id"],
		"first_name": claims["name"],
		"photo_url":  claims["photo_url"],
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
