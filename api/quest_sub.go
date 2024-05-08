package api

import (
	"context"
	"time"

	"github.com/gnasnik/titan-quest/core/dao"
	"github.com/gnasnik/titan-quest/core/generated/model"
)

func completeMission(ctx context.Context, username string, missionID int64) error {
	mission, err := dao.GetMissionById(ctx, missionID)
	if err != nil {
		log.Errorf("GetMissionById: %v", err)
		return err
	}

	ums, err := dao.GetUserMissionByMissionId(ctx, username, mission.ID, dao.QueryOption{})
	if err != nil {
		log.Errorf("GetUserMissionByMissionId: %v", err)
		return err
	}

	if len(ums) == 0 {
		return dao.AddUserMissionAndInviteLog(ctx, &model.UserMission{
			Username:  username,
			MissionID: mission.ID,
			Type:      mission.Type,
			Credit:    mission.Credit,
			Content:   username,
			CreatedAt: time.Now(),
		})
	}

	return nil
}

func getMission(ctx context.Context, username string, missionID int64) (bool, error) {
	mission, err := dao.GetMissionById(ctx, missionID)
	if err != nil {
		log.Errorf("GetMissionById: %v", err)
		return false, err
	}

	ums, err := dao.GetUserMissionByMissionId(ctx, username, mission.ID, dao.QueryOption{})
	if err != nil {
		log.Errorf("GetUserMissionByMissionId: %v", err)
		return false, err
	}

	return len(ums) > 0, nil
}
