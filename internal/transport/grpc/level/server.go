package level

import (
	levelv1 "chinese-game-backend/api/gen/level"
	"chinese-game-backend/internal/service"
	"context"
)

type LevelServer struct {
	levelv1.UnimplementedLevelServiceServer
	levelService *service.LevelService
}

func NewLevelServer(levelService *service.LevelService) *LevelServer {
	return &LevelServer{
		levelService: levelService,
	}
}

func (this *LevelServer) GetLevels(ctx context.Context, req *levelv1.GetLevelsRequest) (*levelv1.GetLevelsResponse, error) {
	levels, err := this.levelService.GetAll(false)
	if err != nil {
		return nil, err
	}

	pbLevels := make([]*levelv1.Level, len(levels))
	for i, level := range levels {
		pbLevels[i] = &levelv1.Level{
			Id:            int64(level.ID),
			Title:         level.Title,
			Color:         level.Color,
			Icon:          level.Icon,
			Order:         int32(level.OrderIndex),
			BackgroundSrc: level.BackgroundSrc,
			PlanetImgSrc:  level.PlanetImgSrc,
			CreatedAt:     level.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}
	}

	return &levelv1.GetLevelsResponse{
		Levels: pbLevels,
	}, nil
}

func (this *LevelServer) GetLevelByID(ctx context.Context, req *levelv1.GetLevelByIDRequest) (*levelv1.GetLevelByIDResponse, error) {
	level, err := this.levelService.GetByID(int(req.Id), int(req.Id))
	if err != nil {
		return nil, err
	}
	pbLevel := &levelv1.Level{
		Id:            int64(level.ID),
		Title:         level.Title,
		Color:         level.Color,
		Icon:          level.Icon,
		Order:         int32(level.OrderIndex),
		BackgroundSrc: level.BackgroundSrc,
		PlanetImgSrc:  level.PlanetImgSrc,
		CreatedAt:     level.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	pbSteps := make([]*levelv1.Step, len(level.Steps))
	for i, step := range level.Steps {
		contentStr := string(step.Content)
		if step.Content == nil {
			contentStr = ""
		}

		// Конвертируем Dialog (уже []DialogStepItem)
		pbDialog := make([]*levelv1.DialogItem, len(step.Dialog))
		for j, dialogItem := range step.Dialog {
			pbDialog[j] = &levelv1.DialogItem{
				Speaker: dialogItem.Speaker,
				Text:    dialogItem.Text,
				Emotion: dialogItem.Emotion,
				Bg:      dialogItem.Bg,
			}
		}

		pbSteps[i] = &levelv1.Step{
			Id:          int64(step.ID),
			LevelId:     int64(step.LevelID),
			Type:        string(step.Type),
			Title:       step.Title,
			Order:       int32(step.OrderIndex),
			Content:     contentStr,
			Description: step.Description,
			Dialog:      pbDialog,
		}
	}

	return &levelv1.GetLevelByIDResponse{
		Level: pbLevel,
		Steps: pbSteps,
	}, nil
}
