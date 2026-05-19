package ai

import (
	"fmt"
	"strings"
	"time"

	"github.com/SaenkoDmitry/training-tg-bot/internal/application/dto"
	"github.com/SaenkoDmitry/training-tg-bot/internal/constants"
	"github.com/SaenkoDmitry/training-tg-bot/internal/models"
	"github.com/SaenkoDmitry/training-tg-bot/internal/repository/exercisetypes"
	"gorm.io/gorm"
)

type ApplyProgramUseCase struct {
	db                *gorm.DB
	exerciseTypesRepo exercisetypes.Repo
}

func NewApplyProgramUseCase(db *gorm.DB, exerciseTypesRepo exercisetypes.Repo) *ApplyProgramUseCase {
	return &ApplyProgramUseCase{
		db:                db,
		exerciseTypesRepo: exerciseTypesRepo,
	}
}

func (uc *ApplyProgramUseCase) Execute(userID int64, req dto.AIApplyProgramRequest) (*dto.AIApplyProgramResponse, error) {
	catalog, err := uc.exerciseTypesRepo.GetAll()
	if err != nil {
		return nil, err
	}
	exerciseTypes := make(map[int64]models.ExerciseType, len(catalog))
	for _, exerciseType := range catalog {
		exerciseTypes[exerciseType.ID] = exerciseType
	}

	if err = validateGeneratedProgram(req.Program, exerciseTypes); err != nil {
		return nil, err
	}

	programName := strings.TrimSpace(req.Program.Name)
	if programName == "" {
		programName = "AI программа"
	}

	result := &dto.AIApplyProgramResponse{Name: programName, DaysCount: len(req.Program.Days)}
	err = uc.db.Transaction(func(tx *gorm.DB) error {
		program := models.WorkoutProgram{
			UserID:          userID,
			Name:            programName,
			CreatedAt:       time.Now(),
			Summary:         new(req.Summary),
			Warnings:        req.Warnings,
			ValidationNotes: req.Notes,
		}
		if err := tx.Create(&program).Error; err != nil {
			return err
		}
		result.ProgramID = program.ID

		for dayIndex, day := range req.Program.Days {
			preset, err := buildPreset(day, exerciseTypes)
			if err != nil {
				return err
			}

			dayType := models.WorkoutDayType{
				WorkoutProgramID: program.ID,
				Name:             emptyFallback(strings.TrimSpace(day.Name), fmt.Sprintf("День %d", dayIndex+1)),
				Preset:           preset,
				CreatedAt:        time.Now().Add(time.Duration(dayIndex) * time.Millisecond),
			}
			if err := tx.Create(&dayType).Error; err != nil {
				return err
			}

			for _, exercise := range day.Exercises {
				if strings.TrimSpace(exercise.ProgressionRule) == "" && strings.TrimSpace(exercise.Reason) == "" {
					continue
				}
				exerciseTypeID := exercise.ExerciseTypeID
				rule := models.WorkoutProgramProgressionRule{
					WorkoutProgramID: program.ID,
					WorkoutDayTypeID: &dayType.ID,
					ExerciseTypeID:   &exerciseTypeID,
					Rule:             strings.TrimSpace(exercise.ProgressionRule),
					Reason:           strings.TrimSpace(exercise.Reason),
					Source:           "ai",
					CreatedAt:        time.Now(),
					UpdatedAt:        time.Now(),
				}
				if rule.Rule == "" {
					rule.Rule = "Правило прогрессии не указано"
				}
				if err := tx.Create(&rule).Error; err != nil {
					return err
				}
				result.RulesCount++
			}
		}

		if req.Activate {
			if err := tx.Model(&models.User{}).Where("id = ?", userID).Update("active_program_id", program.ID).Error; err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return result, nil
}

func validateGeneratedProgram(program dto.AIGeneratedProgram, exerciseTypes map[int64]models.ExerciseType) error {
	if len(program.Days) == 0 {
		return fmt.Errorf("program must contain at least one day")
	}
	if len(program.Days) > 7 {
		return fmt.Errorf("program must contain no more than 7 days")
	}
	for dayIndex, day := range program.Days {
		if len(day.Exercises) == 0 {
			return fmt.Errorf("day %d must contain at least one exercise", dayIndex+1)
		}
		if len(day.Exercises) > 12 {
			return fmt.Errorf("day %d must contain no more than 12 exercises", dayIndex+1)
		}
		for exerciseIndex, exercise := range day.Exercises {
			exerciseType, ok := exerciseTypes[exercise.ExerciseTypeID]
			if !ok {
				return fmt.Errorf("day %d exercise %d has unknown exercise_type_id %d", dayIndex+1, exerciseIndex+1, exercise.ExerciseTypeID)
			}
			if len(exercise.Sets) == 0 {
				return fmt.Errorf("day %d exercise %d must contain at least one set", dayIndex+1, exerciseIndex+1)
			}
			if len(exercise.Sets) > 10 {
				return fmt.Errorf("day %d exercise %d must contain no more than 10 sets", dayIndex+1, exerciseIndex+1)
			}
			for setIndex, set := range exercise.Sets {
				if err := validateSetByUnits(set, exerciseType); err != nil {
					return fmt.Errorf("day %d exercise %d set %d: %w", dayIndex+1, exerciseIndex+1, setIndex+1, err)
				}
			}
		}
	}
	return nil
}

func validateSetByUnits(set dto.AIGeneratedProgramSet, exerciseType models.ExerciseType) error {
	switch {
	case exerciseType.ContainsReps() && exerciseType.ContainsWeight():
		if set.Reps <= 0 || set.Weight <= 0 {
			return fmt.Errorf("%s requires reps and weight", exerciseType.Name)
		}
	case exerciseType.ContainsReps():
		if set.Reps <= 0 {
			return fmt.Errorf("%s requires reps", exerciseType.Name)
		}
	case exerciseType.ContainsMinutes():
		if set.Minutes <= 0 {
			return fmt.Errorf("%s requires minutes", exerciseType.Name)
		}
	case exerciseType.ContainsMeters():
		if set.Meters <= 0 {
			return fmt.Errorf("%s requires meters", exerciseType.Name)
		}
	default:
		return fmt.Errorf("%s has unsupported units %q", exerciseType.Name, exerciseType.Units)
	}
	return nil
}

func buildPreset(day dto.AIGeneratedProgramDay, exerciseTypes map[int64]models.ExerciseType) (string, error) {
	exercises := make([]string, 0, len(day.Exercises))
	for _, exercise := range day.Exercises {
		exerciseType := exerciseTypes[exercise.ExerciseTypeID]
		sets := make([]string, 0, len(exercise.Sets))
		for _, set := range exercise.Sets {
			sets = append(sets, formatGeneratedSet(set, exerciseType))
		}
		exercises = append(exercises, fmt.Sprintf("%d:[%s]", exercise.ExerciseTypeID, strings.Join(sets, ",")))
	}
	return strings.Join(exercises, ";"), nil
}

func formatGeneratedSet(set dto.AIGeneratedProgramSet, exerciseType models.ExerciseType) string {
	if exerciseType.ContainsReps() && exerciseType.ContainsWeight() {
		return fmt.Sprintf("%d*%s", set.Reps, formatPresetWeight(set.Weight))
	}
	if exerciseType.ContainsReps() {
		return fmt.Sprintf("%d*%d", set.Reps, int(constants.DefaultWeight))
	}
	if exerciseType.ContainsMinutes() {
		return fmt.Sprintf("%d", set.Minutes)
	}
	if exerciseType.ContainsMeters() {
		return fmt.Sprintf("%dm", set.Meters)
	}
	return ""
}

func formatPresetWeight(weight float32) string {
	if weight == float32(int(weight)) {
		return fmt.Sprintf("%.0f", weight)
	}
	return fmt.Sprintf("%.1f", weight)
}
