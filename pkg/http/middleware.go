package http

import (
	"bytes"
	"net/http"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func EnsureAuthenticated() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		user := session.Get("user")
		if user == nil {
			c.Redirect(http.StatusMovedPermanently, "/admin/login")
			c.Abort()
		} else {
			c.Next()
		}
	}
}

func EnsureNotAuthenticated() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		user := session.Get("user")
		if user != nil {
			c.Redirect(http.StatusMovedPermanently, "/admin")
			c.Abort()
		} else {
			c.Next()
		}
	}
}

func FaviconMiddleware(file []byte) gin.HandlerFunc {
	reader := bytes.NewReader(file)

	return func(c *gin.Context) {
		if c.Request.RequestURI != "/favicon.ico" {
			return
		}

		if c.Request.Method != "GET" && c.Request.Method != "HEAD" {
			status := http.StatusOK
			if c.Request.Method != "OPTIONS" {
				status = http.StatusMethodNotAllowed
			}
			c.Header("Allow", "GET,HEAD,OPTIONS")
			c.AbortWithStatus(status)
			return
		}

		c.Header("Content-Type", "image/x-icon")
		http.ServeContent(c.Writer, c.Request, "favicon.ico", time.Time{}, reader)
		return
	}
}
