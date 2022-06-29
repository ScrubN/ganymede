package http

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/zibbp/ganymede/ent"
	"github.com/zibbp/ganymede/internal/utils"
	"github.com/zibbp/ganymede/internal/vod"
	"net/http"
)

type VodService interface {
	CreateVod(c echo.Context, vod vod.Vod, cID uuid.UUID) (*ent.Vod, error)
	GetVods(c echo.Context) ([]*ent.Vod, error)
	GetVod(c echo.Context, vID uuid.UUID) (*ent.Vod, error)
	DeleteVod(c echo.Context, vID uuid.UUID) error
	UpdateVod(c echo.Context, vID uuid.UUID, vod vod.Vod, cID uuid.UUID) (*ent.Vod, error)
}

type CreateVodRequest struct {
	ChannelID        string            `json:"channel_id" validate:"required"`
	ExtID            string            `json:"ext_id" validate:"min=1"`
	Platform         utils.VodPlatform `json:"platform" validate:"required,oneof=twitch youtube"`
	Type             utils.VodType     `json:"type" validate:"required,oneof=archive live highlight upload clip"`
	Title            string            `json:"title" validate:"required,min=1"`
	Duration         int               `json:"duration" validate:"required"`
	Views            int               `json:"views" validate:"required"`
	Resolution       string            `json:"resolution"`
	Processing       bool              `json:"processing"`
	ThumbnailPath    string            `json:"thumbnail_path"`
	WebThumbnailPath string            `json:"web_thumbnail_path" validate:"required,min=1"`
	VideoPath        string            `json:"video_path" validate:"required,min=1"`
	ChatPath         string            `json:"chat_path"`
	ChatVideoPath    string            `json:"chat_video_path"`
	InfoPath         string            `json:"info_path"`
}

func (h *Handler) CreateVod(c echo.Context) error {
	var req CreateVodRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if err := c.Validate(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	cUUID, err := uuid.Parse(req.ChannelID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	cvrDto := vod.Vod{
		ExtID:            req.ExtID,
		Platform:         req.Platform,
		Type:             req.Type,
		Title:            req.Title,
		Duration:         req.Duration,
		Views:            req.Views,
		Resolution:       req.Resolution,
		Processing:       req.Processing,
		ThumbnailPath:    req.ThumbnailPath,
		WebThumbnailPath: req.WebThumbnailPath,
		VideoPath:        req.VideoPath,
		ChatPath:         req.ChatPath,
		ChatVideoPath:    req.ChatVideoPath,
		InfoPath:         req.InfoPath,
	}

	v, err := h.Service.VodService.CreateVod(c, cvrDto, cUUID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, v)
}

func (h *Handler) GetVods(c echo.Context) error {
	v, err := h.Service.VodService.GetVods(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, v)
}

func (h *Handler) GetVod(c echo.Context) error {
	vID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	v, err := h.Service.VodService.GetVod(c, vID)
	if err != nil {
		if err.Error() == "vod not found" {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, v)
}

func (h *Handler) DeleteVod(c echo.Context) error {
	vID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	err = h.Service.VodService.DeleteVod(c, vID)
	if err != nil {
		if err.Error() == "vod not found" {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.NoContent(http.StatusOK)
}

func (h *Handler) UpdateVod(c echo.Context) error {
	vID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	var req CreateVodRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if err := c.Validate(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	cUUID, err := uuid.Parse(req.ChannelID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	cvrDto := vod.Vod{
		ExtID:            req.ExtID,
		Platform:         req.Platform,
		Type:             req.Type,
		Title:            req.Title,
		Duration:         req.Duration,
		Views:            req.Views,
		Resolution:       req.Resolution,
		Processing:       req.Processing,
		ThumbnailPath:    req.ThumbnailPath,
		WebThumbnailPath: req.WebThumbnailPath,
		VideoPath:        req.VideoPath,
		ChatPath:         req.ChatPath,
		ChatVideoPath:    req.ChatVideoPath,
		InfoPath:         req.InfoPath,
	}

	v, err := h.Service.VodService.UpdateVod(c, vID, cvrDto, cUUID)
	if err != nil {
		if err.Error() == "vod not found" {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, v)
}