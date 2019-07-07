package http

import (
	"math/rand"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"

	"github.com/phillebaba/bike-touring-tracker/pkg/domain"
)

type AdminHandler struct {
	CheckinService domain.CheckinService
}

func (h *AdminHandler) Index(c *gin.Context) {
	checkins := h.CheckinService.List()
	c.HTML(http.StatusOK, "admin", gin.H{
		"Checkins": checkins,
	})
}

func (h *AdminHandler) Checkin(c *gin.Context) {
	c.HTML(http.StatusOK, "checkin", nil)
}

func (h *AdminHandler) SubmitCheckin(c *gin.Context) {
	lat, _ := strconv.ParseFloat(c.PostForm("lat"), 64)
	lng, _ := strconv.ParseFloat(c.PostForm("lng"), 64)
	precision, _ := strconv.ParseInt(c.PostForm("precision"), 10, 32)
	description := c.PostForm("description")

	h.CheckinService.Register(lat, lng, int(precision), description)

	c.Redirect(http.StatusMovedPermanently, "/admin")
	c.Abort()
}

func (h *AdminHandler) DeleteCheckin(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 32)
	h.CheckinService.Delete(int(id))

	c.Redirect(http.StatusMovedPermanently, "/admin")
	c.Abort()
}

func (h *AdminHandler) LoginView(c *gin.Context) {
	c.HTML(http.StatusOK, "login", nil)
}

func (h *AdminHandler) Login(c *gin.Context) {
	session := sessions.Default(c)
	username := c.PostForm("username")
	password := c.PostForm("password")

	if strings.Trim(username, " ") == "" || strings.Trim(password, " ") == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Parameters can't be empty"})
		return
	}

	if username == "admin" && password == "password" {
		session.Set("user", username) //In real world usage you'd set this to the users ID
		err := session.Save()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate session token"})
		} else {
			c.Redirect(http.StatusMovedPermanently, "/admin/login")
			c.Abort()
		}
	} else {
		c.Redirect(http.StatusMovedPermanently, "/admin")
		c.Abort()
	}
}

func (h *AdminHandler) Logout(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get("user")
	if user == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid session token"})
	} else {
		session.Delete("user")
		session.Save()
		c.JSON(http.StatusOK, gin.H{"message": "Successfully logged out"})
	}
}

func generateSessionToken() string {
	return strconv.FormatInt(rand.Int63(), 16)
}
