package api

import (
	"context"
	"github.com/gnasnik/titan-quest/core/dao"
	"github.com/gnasnik/titan-quest/core/generated/model"
	"time"
)

func completeConnectWalletMission(ctx context.Context, address string) error {
	mission, err := dao.GetMissionById(ctx, MissionIdConnectWallet)
	if err != nil {
		log.Errorf("GetMissionById: %v", err)
		return err
	}

	ums, err := dao.GetUserMissionByMissionId(ctx, address, mission.ID, dao.QueryOption{})
	if err != nil {
		log.Errorf("GetUserMissionByMissionId: %v", err)
		return err
	}

	if len(ums) == 0 {
		return dao.AddUserMission(ctx, &model.UserMission{
			Username:  address,
			MissionID: mission.ID,
			Type:      mission.Type,
			Credit:    mission.Credit,
			Content:   address,
			CreatedAt: time.Now(),
		})
	}

	return nil
}
