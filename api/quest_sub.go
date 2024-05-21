package api

import (
	"context"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/gnasnik/titan-quest/config"
	"github.com/gnasnik/titan-quest/core/dao"
	"github.com/gnasnik/titan-quest/core/generated/model"
	"github.com/gnasnik/titan-quest/pkg/opgoogle"
)

var (
	credJSON, tokenJSON []byte
	mu                  = new(sync.Mutex)
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

func getCredAndToken() ([]byte, []byte, error) {
	mu.Lock()
	defer mu.Unlock()

	if credJSON == nil {
		cf, err := os.ReadFile(config.Cfg.GoogleDoc.CredPath)
		if err != nil {
			return nil, nil, fmt.Errorf("read credjson error:%w", err)
		}
		credJSON = cf
	}

	if tokenJSON == nil {
		tf, err := os.ReadFile(config.Cfg.GoogleDoc.TokenPath)
		if err != nil {
			return nil, nil, fmt.Errorf("read tokenjson error:%w", err)
		}
		tokenJSON = tf
	}

	return credJSON, tokenJSON, nil
}

func checkBVComplete(un, speedID string) (bool, error) {
	cj, tj, err := getCredAndToken()
	if err != nil {
		return false, err
	}

	sheetSrv, err := opgoogle.GetSheetService(cj, tj)
	if err != nil {
		return false, fmt.Errorf("get sheet service error:%w", err)
	}
	resp, err := sheetSrv.Spreadsheets.Values.Get(speedID, "sheet1!F2:F").Do()
	if err != nil {
		return false, fmt.Errorf("get body of sheet error:%w", err)
	}

	if len(resp.Values) == 0 {
		return false, nil
	}

	for _, row := range resp.Values {
		if len(row) == 0 {
			continue
		}
		key, ok := row[0].(string)
		if !ok {
			continue
		}

		if strings.EqualFold(key, un) {
			return true, nil
		}
	}

	return false, nil
}

func checkRFComplete(un, speedID string) (bool, error) {
	cj, tj, err := getCredAndToken()
	if err != nil {
		return false, err
	}

	sheetSrv, err := opgoogle.GetSheetService(cj, tj)
	if err != nil {
		return false, fmt.Errorf("get sheet service error:%w", err)
	}
	resp, err := sheetSrv.Spreadsheets.Values.Get(speedID, "sheet1!B2:E").Do()
	if err != nil {
		return false, fmt.Errorf("get body of sheet error:%w", err)
	}

	if len(resp.Values) == 0 {
		return false, nil
	}

	for _, row := range resp.Values {
		if len(row) == 0 {
			continue
		}
		key, ok := row[0].(string)
		if !ok {
			continue
		}

		if strings.EqualFold(key, un) {
			// 校验审核状态
			if status, ok := row[len(row)-1].(string); ok && strings.TrimSpace(status) == "通过" {
				return true, nil
			}
		}
	}

	return false, nil
}
