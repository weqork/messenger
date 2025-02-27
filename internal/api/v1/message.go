package apiv1

import (
	"errors"
	"github.com/HalvaPovidlo/messenger/internal/pkg/message"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type msgBody struct {
	Text string `json:"text"`
}

func (h *handler) PersonalMessage(c echo.Context) error {
	to, err := uuid.Parse(c.Param("to"))
	if err != nil {
		return c.String(http.StatusBadRequest, "")
	}

	from, err := getUserID(c)
	if err != nil {
		return c.String(http.StatusInternalServerError, "")
	}

	var body msgBody
	if err := c.Bind(&body); err != nil {
		return err
	}

	if err := h.messages.Message(from, to, body.Text); err != nil {
		return c.String(http.StatusInternalServerError, "")
	}
	return c.String(http.StatusOK, "")

}

type historyOut struct {
	History []message.Message `json:"history"`
}

func (h *handler) PersonalHistory(c echo.Context) error {
	to, err := uuid.Parse(c.Param("to"))
	if err != nil {
		return c.String(http.StatusBadRequest, "")
	}

	from, err := getUserID(c)
	if err != nil {
		return c.String(http.StatusInternalServerError, "")
	}
	var history historyOut
	history.History, err = h.messages.History(from, to)
	if err != nil {
		return c.String(http.StatusNotFound, err.Error())
	}

	return c.JSON(http.StatusOK, history)
}

func getUserID(c echo.Context) (uuid.UUID, error) {
	v, ok := c.Get(userIDKey).(uuid.UUID)
	if ok {
		return v, nil
	}
	return uuid.Nil, errors.New("username is empty")
}
