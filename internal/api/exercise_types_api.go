package api

import (
	"encoding/json"
	"net/http"
)

func (s *serviceImpl) GetExerciseGroups(w http.ResponseWriter, _ *http.Request) {
	result, err := s.container.GetAllGroupsUC.Execute()
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result.Groups)
}

func (s *serviceImpl) GetExerciseTypesByGroup(w http.ResponseWriter, r *http.Request) {
	groupCode := r.PathValue("group")

	result, err := s.container.FindTypesByGroupUC.Execute(groupCode)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result.ExerciseTypes)
}
