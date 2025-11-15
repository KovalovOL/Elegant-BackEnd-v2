package user

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Handler struct {
	serv *Service
}

func NewHandler(serv *Service) *Handler {
	return &Handler{serv}
}

func (h *Handler) Create(c *gin.Context) {
	var user CreateUser
	if err := c.BindJSON(&user); err != nil {
		c.JSON(400, err.Error())
		return
	}	
	ctx := c.Request.Context()

	id, err := h.serv.Create(ctx, user)
	if err != nil {
		c.JSON(500, err.Error())
		return
	}
	c.JSON(200, id)
}

func (h *Handler) GetByID(c *gin.Context) {
	ctx := c.Request.Context()

	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(400, fmt.Errorf("failed to parse uuid: %w", err))
		return
	}

	user, err := h.serv.GetByID(ctx, id)
	if err != nil {
		c.JSON(500, err.Error())
		return
	}
	c.JSON(200, user)

}
