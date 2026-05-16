package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/SaenkoDmitry/training-tg-bot/internal/application/dto"
	"github.com/golang-jwt/jwt/v4"
)

func (s *serviceImpl) YandexRedirectHandler(w http.ResponseWriter, r *http.Request) {
	origin := r.URL.Query().Get("origin")
	state := r.URL.Query().Get("state")

	if origin == "" || state == "" {
		http.Error(w, "missing params", http.StatusBadRequest)
		return
	}

	if !s.isAllowedOrigin(origin) {
		http.Error(w, "invalid origin", http.StatusForbidden)
		return
	}

	clientID := os.Getenv("YANDEX_CLIENT_ID")
	redirectURI := origin + "/auth-yandex"

	authURL := fmt.Sprintf(
		"https://oauth.yandex.ru/authorize?response_type=code&client_id=%s&redirect_uri=%s&state=%s&scope=%s",
		url.QueryEscape(clientID),
		url.QueryEscape(redirectURI),
		url.QueryEscape(state),
		url.QueryEscape("login:birthday login:info"),
	)

	http.Redirect(w, r, authURL, http.StatusFound)
}

func (s *serviceImpl) YandexLoginHandler(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Code string `json:"code"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	tokenResp, err := exchangeCodeForToken(body.Code)
	if err != nil {
		fmt.Println("yandex token exchange failed error:", err.Error())
		http.Error(w, "token exchange failed", http.StatusInternalServerError)
		return
	}

	profile, err := getYandexProfile(tokenResp.AccessToken)
	if err != nil {
		fmt.Println("get yandex profile error:", err.Error())
		http.Error(w, "profile fetch failed", http.StatusInternalServerError)
		return
	}

	user, err := s.container.GetOrCreateUserByYandexUC.Execute(profile)
	if err != nil {
		fmt.Println("create yandex user error:", err.Error())
		http.Error(w, "user error", http.StatusInternalServerError)
		return
	}

	claims := jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(7 * 24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, _ := token.SignedString(jwtSecret)

	json.NewEncoder(w).Encode(map[string]string{
		"token": signed,
	})
}

func exchangeCodeForToken(code string) (*TokenResponse, error) {
	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("code", code)
	data.Set("client_id", os.Getenv("YANDEX_CLIENT_ID"))
	data.Set("client_secret", os.Getenv("YANDEX_CLIENT_SECRET"))

	resp, err := http.PostForm("https://oauth.yandex.ru/token", data)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result TokenResponse
	if err = json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

type TokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}

func getYandexProfile(token string) (*dto.YandexProfile, error) {
	req, _ := http.NewRequest("GET", "https://login.yandex.ru/info", nil)
	req.Header.Set("Authorization", "OAuth "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var profile dto.YandexProfile
	if err := json.NewDecoder(resp.Body).Decode(&profile); err != nil {
		return nil, err
	}

	return &profile, nil
}
