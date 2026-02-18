package api

import (
	"errors"
	"github.com/SaenkoDmitry/training-tg-bot/internal/repository/users"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/golang-jwt/jwt/v4"
	"net/http"
	"strconv"
	"time"
)

func (s *serviceImpl) TelegramCallbackHandler(w http.ResponseWriter, r *http.Request) {
	// парсим query параметры, которые Telegram присылает

	id, _ := strconv.ParseInt(r.URL.Query().Get("id"), 10, 64)
	authDate, _ := strconv.ParseInt(r.URL.Query().Get("auth_date"), 10, 64)
	tgUser := TelegramUser{
		ID:        id,
		FirstName: r.URL.Query().Get("first_name"),
		LastName:  r.URL.Query().Get("last_name"),
		Username:  r.URL.Query().Get("username"),
		PhotoURL:  r.URL.Query().Get("photo_url"),
		AuthDate:  authDate,
		Hash:      r.URL.Query().Get("hash"),
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

	// Создаём юзера, если ещё нет
	if _, err := s.container.GetUserUC.Execute(tgUser.ID); err != nil && errors.Is(err, users.NotFoundUserErr) {
		s.container.CreateUserUC.Execute(tgUser.ID, &tgbotapi.User{
			ID:           tgUser.ID,
			IsBot:        false,
			FirstName:    tgUser.FirstName,
			LastName:     tgUser.LastName,
			UserName:     tgUser.Username,
			LanguageCode: tgUser.LanguageCode,
		})
	}

	// Ставим cookie
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, _ := token.SignedString(jwtSecret)

	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    signed,
		Path:     "/",
		HttpOnly: true,
		Secure:   true, // true на проде
		SameSite: http.SameSiteLaxMode,
	})

	// Редиректим обратно на SPA (например на /profile)
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(`
			<!DOCTYPE html>
			<html>
			<head>
			<meta charset="utf-8" />
			<script>
			  window.location.replace("/");
			</script>
			</head>
			<body></body>
			</html>
`))

}
