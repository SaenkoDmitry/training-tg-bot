package api

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"github.com/golang-jwt/jwt/v4"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

var (
	botToken  = os.Getenv("TELEGRAM_TOKEN")
	jwtSecret = []byte(os.Getenv("JWT_SECRET"))
)

type TelegramUser struct {
	ID        int64  `json:"id"`
	FirstName string `json:"first_name"`
	Username  string `json:"username"`
	PhotoURL  string `json:"photo_url"`
	AuthDate  int64  `json:"auth_date"`
	Hash      string `json:"hash"`
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var tgUser TelegramUser
	if err := json.NewDecoder(r.Body).Decode(&tgUser); err != nil {
		http.Error(w, "bad request: "+err.Error(), http.StatusBadRequest)
		return
	}

	if !verifyTelegram(tgUser, botToken) {
		http.Error(w, "invalid telegram hash", http.StatusUnauthorized)
		return
	}

	// Создаём JWT
	claims := jwt.MapClaims{
		"id":        tgUser.ID,
		"name":      tgUser.FirstName,
		"photo_url": tgUser.PhotoURL,
		"exp":       time.Now().Add(7 * 24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(jwtSecret)
	if err != nil {
		http.Error(w, "internal server error", 500)
		return
	}

	// Ставим httpOnly cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    signed,
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // true на проде
		SameSite: http.SameSiteLaxMode,
	})

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"ok": true}`))
}

func verifyTelegram(user TelegramUser, botToken string) bool {
	// Создаём map[string]string для HMAC
	data := map[string]string{
		"id":         strconv.FormatInt(user.ID, 10),
		"first_name": user.FirstName,
		"auth_date":  strconv.FormatInt(user.AuthDate, 10),
	}
	if user.Username != "" {
		data["username"] = user.Username
	}
	if user.PhotoURL != "" {
		data["photo_url"] = user.PhotoURL
	}

	// сортируем ключи
	keys := make([]string, 0, len(data))
	for k := range data {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// формируем check string
	var parts []string
	for _, k := range keys {
		parts = append(parts, k+"="+data[k])
	}
	checkString := strings.Join(parts, "\n")

	// secret = sha256(botToken)
	secret := sha256.Sum256([]byte(botToken))

	mac := hmac.New(sha256.New, secret[:])
	mac.Write([]byte(checkString))
	expected := hex.EncodeToString(mac.Sum(nil))

	return expected == user.Hash
}
