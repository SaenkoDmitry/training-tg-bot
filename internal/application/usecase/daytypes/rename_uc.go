package daytypes

import (
	"github.com/SaenkoDmitry/training-tg-bot/internal/repository/daytypes"
)

type RenameUseCase struct {
	dayTypesRepo daytypes.Repo
}

func NewRenameUseCase(
	dayTypesRepo daytypes.Repo,
) *RenameUseCase {
	return &RenameUseCase{
		dayTypesRepo: dayTypesRepo,
	}
}

func (uc *RenameUseCase) Name() string {
	return "Переименовать тренировочный день"
}

func (uc *RenameUseCase) Execute(dayTypeID int64, newName string) error {
	dayType, err := uc.dayTypesRepo.Get(dayTypeID)
	if err != nil {
		return err
	}
	dayType.Name = newName
	return uc.dayTypesRepo.Save(&dayType)
}
