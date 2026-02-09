package dto

type Measurement struct {
	ID        int64  `json:"id"`
	CreatedAt string `json:"created_at"`

	Shoulders string `json:"shoulders"`
	Chest     string `json:"chest"`
	HandLeft  string `json:"hand_left"`
	HandRight string `json:"hand_right"`
	Waist     string `json:"waist"`
	Buttocks  string `json:"buttocks"`
	HipLeft   string `json:"hip_left"`
	HipRight  string `json:"hip_right"`
	CalfLeft  string `json:"calf_left"`
	CalfRight string `json:"calf_right"`
	Weight    string `json:"weight"`
}

type FindWithOffsetLimitMeasurement struct {
	Items []Measurement `json:"items"`
	Count int           `json:"count"`
}
