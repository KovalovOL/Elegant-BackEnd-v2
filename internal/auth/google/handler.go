package google

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
	c.SetCookie("state", state, 3*60, "/", "", false, true)
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
	ip := c.ClientIP()
	agent := c.Request.UserAgent()
	refToken, _ := c.Cookie("refresh_token")
	resp, err := h.serv.HandleGoogleCallback(ctx, refToken, code, ip, agent)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.SetCookie("access_token", resp.AccessToken, 15*60, "/", "", false, true)
	c.SetCookie("refresh_token", resp.RefreshToken, 30*24*60*60, "/", "", false, true)
	c.JSON(http.StatusOK, resp.User)
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

func (h *Handler) Refresh(c *gin.Context) {
	refToken, err := c.Cookie("refresh_token")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user is not logged in"})
		return
	}
	ctx := c.Request.Context()

	accToken, err := h.serv.Refresh(ctx, refToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.SetCookie("access_token", accToken, 0.25*60*60, "/", "", false, true)
	c.SetCookie("refresh_token", refToken, 30*24*60*60, "/", "", false, true)
	c.JSON(http.StatusOK, gin.H{"message": "refresh successful"})
}


func (h *Handler) Logout(c *gin.Context) {
	refToken, _ := c.Cookie("refresh_token")
	ctx := c.Request.Context()
	err := h.serv.Logout(ctx, refToken)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.SetCookie("access_token", "", -1, "/", "", false, true)
	c.SetCookie("refresh_token", "", -1, "/", "", false, true)
	c.JSON(http.StatusOK, gin.H{"message": "Logged out"})
}