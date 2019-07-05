package http

import (
	"net/http"

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
