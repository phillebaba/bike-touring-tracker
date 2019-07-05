package http

import (
	"github.com/gin-contrib/multitemplate"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/thinkerou/favicon"

	"github.com/phillebaba/bike-tracker/pkg/domain"
)

func Run(serviceContext domain.ServiceContext) {
	router := gin.Default()

	store := cookie.NewStore([]byte("secret"))
	router.Use(sessions.Sessions("mysession", store))

	router.Static("/static", "web/static")
	router.Use(favicon.New("web/static/favicon.ico"))

	templates := multitemplate.New()
	templates.AddFromFiles("index", "web/templates/base.html", "web/templates/index.html")
	templates.AddFromFiles("admin", "web/templates/base.html", "web/templates/admin.html")
	templates.AddFromFiles("checkin", "web/templates/base.html", "web/templates/checkin.html")
	templates.AddFromFiles("login", "web/templates/base.html", "web/templates/login.html")
	router.HTMLRender = templates

	// Home
	homeHandler := HomeHandler{serviceContext.TripService}
	router.GET("/", homeHandler.Index)

	// Admin
	adminHandler := AdminHandler{serviceContext.CheckinService}
	adminRoutes := router.Group("/admin")
	{
		adminRoutes.GET("/", EnsureAuthenticated(), adminHandler.Index)

		adminRoutes.GET("/checkin", EnsureAuthenticated(), adminHandler.Checkin)
		adminRoutes.POST("/checkin", EnsureAuthenticated(), adminHandler.SubmitCheckin)
		adminRoutes.POST("/checkin/delete/:id", EnsureAuthenticated(), adminHandler.DeleteCheckin)

		adminRoutes.GET("/login", EnsureNotAuthenticated(), adminHandler.LoginView)
		adminRoutes.POST("/login", EnsureNotAuthenticated(), adminHandler.Login)
	}

	router.Run()
}
