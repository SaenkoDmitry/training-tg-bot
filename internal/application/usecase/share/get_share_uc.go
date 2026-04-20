package share

import (
	"github.com/SaenkoDmitry/training-tg-bot/internal/models"
	"github.com/SaenkoDmitry/training-tg-bot/internal/repository/share"
)

type GetShareUC struct {
	shareRepo share.Repo
}

func NewGetShareUC(shareRepo share.Repo) *GetShareUC {
	return &GetShareUC{shareRepo: shareRepo}
}

func (uc *GetShareUC) Execute(token string) (*models.WorkoutShare, error) {
	shareModel, err := uc.shareRepo.Get(token)
	if err != nil {
		return nil, err
	}
	return &shareModel, nil
}
