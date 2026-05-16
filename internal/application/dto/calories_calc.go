package dto

type CaloriesCalc struct {
	Calories    float64  `json:"calories"`
	DurationMin int      `json:"duration_min"`
	UserWeight  *float64 `json:"user_weight"`
	WeightNote  string   `json:"weight_note"`
}
