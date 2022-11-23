package http

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

type TaskService interface {
	StartTask(c echo.Context, task string) error
}

type StartTaskRequest struct {
	Task string `json:"task" validate:"required,oneof=check_live check_vod get_jwks twitch_auth queue_hold_check"`
}

func (h *Handler) StartTask(c echo.Context) error {
	str := new(StartTaskRequest)
	if err := c.Bind(str); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if err := c.Validate(str); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if err := h.Service.TaskService.StartTask(c, str.Task); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.NoContent(http.StatusOK)
}