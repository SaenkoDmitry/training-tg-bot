package api

import (
	"fmt"
	"log"
	"net/http"

	"github.com/SaenkoDmitry/training-tg-bot/internal/middlewares"
)

func (s *serviceImpl) DownloadExcelWorkoutsStats(w http.ResponseWriter, r *http.Request) {
	claims, ok := middlewares.FromContext(r.Context())
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	resp, err := s.container.ExportWorkoutsToExcelUC.Execute(claims.UserID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Устанавливаем заголовки для скачивания Excel
	w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	w.Header().Set("Content-Disposition", "attachment; filename=workouts.xlsx")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", resp.Len()))

	// Записываем данные в ResponseWriter
	if _, err = w.Write(resp.Bytes()); err != nil {
		log.Printf("error writing response: %v", err)
	}
}
