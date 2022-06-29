package queue

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/zibbp/ganymede/ent"
	"github.com/zibbp/ganymede/ent/queue"
	"github.com/zibbp/ganymede/internal/database"
	"github.com/zibbp/ganymede/internal/utils"
	"time"
)

type Service struct {
	Store *database.Database
}

func NewService(store *database.Database) *Service {
	return &Service{Store: store}
}

type Queue struct {
	ID                       uuid.UUID        `json:"id"`
	LiveArchive              bool             `json:"live_archive"`
	OnHold                   bool             `json:"on_hold"`
	VideoProcessing          bool             `json:"video_processing"`
	ChatProcessing           bool             `json:"chat_processing"`
	Processing               bool             `json:"processing"`
	TaskVodCreateFolder      utils.TaskStatus `json:"task_vod_create_folder"`
	TaskVodDownloadThumbnail utils.TaskStatus `json:"task_vod_download_thumbnail"`
	TaskVodSaveInfo          utils.TaskStatus `json:"task_vod_save_info"`
	TaskVideoDownload        utils.TaskStatus `json:"task_video_download"`
	TaskVideoMove            utils.TaskStatus `json:"task_video_move"`
	TaskChatDownload         utils.TaskStatus `json:"task_chat_download"`
	TaskChatRender           utils.TaskStatus `json:"task_chat_render"`
	TaskChatMove             utils.TaskStatus `json:"task_chat_move"`
	UpdatedAt                time.Time        `json:"updated_at"`
	CreatedAt                time.Time        `json:"created_at"`
}

func (s *Service) CreateQueueItem(c echo.Context, queueDto Queue, vID uuid.UUID) (*ent.Queue, error) {
	q, err := s.Store.Client.Queue.Create().SetVodID(vID).Save(c.Request().Context())
	if err != nil {
		if _, ok := err.(*ent.ConstraintError); ok {
			return nil, fmt.Errorf("queue item exists for vod or vod does not exist")
		}
		log.Debug().Err(err).Msg("error creating queue")
		return nil, fmt.Errorf("error creating queue: %v", err)
	}
	return q, nil
}

func (s *Service) UpdateQueueItem(c echo.Context, queueDto Queue, qID uuid.UUID) (*ent.Queue, error) {
	q, err := s.Store.Client.Queue.UpdateOneID(qID).SetLiveArchive(queueDto.LiveArchive).SetOnHold(queueDto.OnHold).SetVideoProcessing(queueDto.VideoProcessing).SetChatProcessing(queueDto.ChatProcessing).SetProcessing(queueDto.Processing).SetTaskVodCreateFolder(queueDto.TaskVodCreateFolder).SetTaskVodDownloadThumbnail(queueDto.TaskVodDownloadThumbnail).SetTaskVodSaveInfo(queueDto.TaskVodSaveInfo).SetTaskVideoDownload(queueDto.TaskVideoDownload).SetTaskVideoMove(queueDto.TaskVideoMove).SetTaskChatDownload(queueDto.TaskChatDownload).SetTaskChatRender(queueDto.TaskChatRender).SetTaskChatMove(queueDto.TaskChatMove).Save(c.Request().Context())
	if err != nil {
		return nil, fmt.Errorf("error updating queue: %v", err)
	}
	return q, nil
}

func (s *Service) GetQueueItems(c echo.Context) ([]*ent.Queue, error) {
	q, err := s.Store.Client.Queue.Query().WithVod().All(c.Request().Context())
	if err != nil {
		return nil, fmt.Errorf("error getting queue tasks: %v", err)
	}
	return q, nil
}
func (s *Service) GetQueueItemsFilter(c echo.Context, processing bool) ([]*ent.Queue, error) {
	q, err := s.Store.Client.Queue.Query().Where(queue.Processing(processing)).WithVod().All(c.Request().Context())
	if err != nil {
		return nil, fmt.Errorf("error getting queue tasks: %v", err)
	}
	return q, nil
}

func (s *Service) DeleteQueueItem(c echo.Context, qID uuid.UUID) error {
	err := s.Store.Client.Queue.DeleteOneID(qID).Exec(c.Request().Context())
	if err != nil {
		return fmt.Errorf("error deleting queue: %v", err)
	}
	return nil
}

func (s *Service) GetQueueItem(c echo.Context, qID uuid.UUID) (*ent.Queue, error) {
	q, err := s.Store.Client.Queue.Query().Where(queue.ID(qID)).WithVod().Only(c.Request().Context())
	if err != nil {
		return nil, fmt.Errorf("error getting queue task: %v", err)
	}
	return q, nil
}

func (s *Service) ReadLogFile(c echo.Context, qID uuid.UUID, logType string) ([]byte, error) {
	q, err := s.GetQueueItem(c, qID)
	if err != nil {
		return nil, err
	}
	path := fmt.Sprintf("/logs/%s_%s-%s.log", q.Edges.Vod.ExtID, q.Edges.Vod.ID, logType)
	logLines, err := utils.ReadLastLines(path, "10")
	if err != nil {
		return nil, err
	}
	return []byte(logLines), nil
}