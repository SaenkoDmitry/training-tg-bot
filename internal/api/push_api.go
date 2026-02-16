package api

import (
	"encoding/json"
	"fmt"
	"github.com/SaenkoDmitry/training-tg-bot/internal/application/dto"
	"github.com/SaenkoDmitry/training-tg-bot/internal/middlewares"
	"github.com/SaenkoDmitry/training-tg-bot/internal/models"
	"github.com/SherClockHolmes/webpush-go"
	"net/http"
	"os"
)

func (s *serviceImpl) PushSubscribe(w http.ResponseWriter, r *http.Request) {
	fmt.Println("PushSubscribe")
	claims, ok := middlewares.FromContext(r.Context())
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	var sub dto.PushSubscription
	if err := json.NewDecoder(r.Body).Decode(&sub); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err := s.container.CreatePushSubscriptionUC.Execute(claims.ChatID, sub)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func sendPush(sub *models.PushSubscription, payload map[string]string) (int, error) {
	subscription := &webpush.Subscription{
		Endpoint: sub.Endpoint,
		Keys: webpush.Keys{
			Auth:   sub.Auth,
			P256dh: sub.P256dh,
		},
	}
	options := &webpush.Options{
		TTL:             3600,
		VAPIDPrivateKey: os.Getenv("VAPID_PRIVATE_KEY"),
		VAPIDPublicKey:  os.Getenv("VAPID_PUBLIC_KEY"),
		Subscriber:      "https://form-journey.ru",
	}

	payloadJson, _ := json.Marshal(payload)

	resp, err := webpush.SendNotification(payloadJson, subscription, options)
	if err != nil {
		// если ошибка, может быть nil
		if resp != nil {
			return resp.StatusCode, err
		}
		return 0, err
	}
	fmt.Println("Push sent, status:", resp.StatusCode)
	defer resp.Body.Close()

	return resp.StatusCode, nil
}
