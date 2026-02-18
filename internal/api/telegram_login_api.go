package api

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"github.com/SaenkoDmitry/training-tg-bot/internal/repository/users"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
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
	ID           int64  `json:"id"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Username     string `json:"username"`
	LanguageCode string `json:"language_code"`
	PhotoURL     string `json:"photo_url"`
	AuthDate     int64  `json:"auth_date"`
	Hash         string `json:"hash"`
}

func (s *serviceImpl) TelegramLoginHandler(w http.ResponseWriter, r *http.Request) {

	var tgUser TelegramUser

	if err := json.NewDecoder(r.Body).Decode(&tgUser); err != nil {
		http.Error(w, "bad request", 400)
		return
	}

	if !verifyTelegram(tgUser, botToken) {
		http.Error(w, "invalid telegram hash", http.StatusUnauthorized)
		return
	}

	claims := jwt.MapClaims{
		"id":        tgUser.ID,
		"name":      tgUser.FirstName,
		"photo_url": tgUser.PhotoURL,
		"exp":       time.Now().Add(7 * 24 * time.Hour).Unix(),
	}

	if _, err := s.container.GetUserUC.Execute(tgUser.ID); err != nil && errors.Is(err, users.NotFoundUserErr) {
		s.container.CreateUserUC.Execute(tgUser.ID, &tgbotapi.User{
			ID:        tgUser.ID,
			FirstName: tgUser.FirstName,
			LastName:  tgUser.LastName,
			UserName:  tgUser.Username,
		})
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, _ := token.SignedString(jwtSecret)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"token": signed,
	})
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
