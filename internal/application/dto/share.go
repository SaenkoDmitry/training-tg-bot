package dto

type CreateShareRequest struct {
	WorkoutID int64 `json:"workout_id"`
}

type ShareResponse struct {
	Token     string `json:"token"`
	ShareURL  string `json:"share_url"`
	CreatedAt string `json:"created_at"`
}

type PublicWorkoutDTO struct {
	Progress *WorkoutProgress  `json:"progress"`
	Stats    *WorkoutStatistic `json:"stats"`
}
