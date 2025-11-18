package auth

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	serv *Service
} 

func NewHandler(serv *Service) *Handler {
	return &Handler{serv}
}

func (h *Handler) LoginGoogle(c *gin.Context) {
	state, url := h.serv.StartGoogleAuth()
	c.SetCookie("state", state, 7*60, "/", "", false, true)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func (h *Handler) GoogleCallback(c *gin.Context) {
	state, err := c.Cookie("state")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "state not found"})
		return
	}

	if state != c.Query("state") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid state"})
		return
	}
	ctx := c.Request.Context()

	code := c.Query("code")
	resp, err := h.serv.HandleGoogleCallback(ctx, code)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid code"})
	}

	c.SetCookie("access_token", resp.AccessToken, 24*60*60, "/", "", false, true)
	c.JSON(http.StatusOK, resp.AccessToken)
}

func (h *Handler) Me(c *gin.Context) {
	id := c.GetString("user_id")
	name := c.GetString("user_name")
	email := c.GetString("user_email")
	c.JSON(http.StatusOK, gin.H{
		"id":    id,
		"name":  name,
		"email": email,
	})
}


func (h *Handler) Logout(c *gin.Context) {
	c.SetCookie("token", "", -1, "/", "", false, true)
	c.JSON(http.StatusOK, gin.H{"message": "Logged out"})
}