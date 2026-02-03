package models

import "time"

type Measurement struct {
	ID        int64 `gorm:"primaryKey;autoIncrement"`
	UserID    int64
	CreatedAt time.Time

	Shoulders int // Плечи
	Chest     int // Грудь
	HandLeft  int // Рука левая
	HandRight int // Рука правая
	Waist     int // Талия
	Buttocks  int // Ягодицы
	HipLeft   int // Бедро левое
	HipRight  int // Бедро правое
	CalfLeft  int // Икра левая
	CalfRight int // Икра правая
	Weight    int // Вес
}

func (*Measurement) TableName() string {
	return "measurements"
}
