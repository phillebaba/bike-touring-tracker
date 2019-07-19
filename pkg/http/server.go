package http

import (
	"html/template"
	"net/http"

	"github.com/gin-contrib/multitemplate"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/gobuffalo/packr/v2"

	"github.com/phillebaba/bike-touring-tracker/pkg/domain"
)

func Run(adminPassword string, serviceContext domain.ServiceContext) {
	router := gin.Default()
	router.HTMLRender = loadTemplates()

	store := cookie.NewStore([]byte("secret"))
	router.Use(sessions.Sessions("mysession", store))

	staticBox := packr.New("static", "../../web/static")
	router.StaticFS("/static", http.FileSystem(staticBox))

	faviconBox := packr.New("favicon", "../../web/favicon")
	file, _ := faviconBox.Find("favicon.ico")
	router.Use(FaviconMiddleware(file))

	// Home
	homeHandler := HomeHandler{serviceContext.TripService}
	router.GET("/", homeHandler.Index)

	// Admin
	adminHandler := AdminHandler{adminPassword, serviceContext.CheckinService}
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

func loadTemplates() multitemplate.Renderer {
	templateBox := packr.New("template", "../../web/templates")
	baseTemplateString, err := templateBox.FindString("base.html")
	if err != nil {
		panic("Could not get base.html file")
	}
	baseTemplate, _ := template.New("base").Parse(baseTemplateString)

	renderer := multitemplate.NewRenderer()
	for _, page := range []string{"index", "admin", "checkin", "login"} {
		pageFilePath := page + ".html"
		pageTemplateString, err := templateBox.FindString(pageFilePath)
		if err != nil {
			panic("Could not get page file")
		}

		pageTemplate, _ := template.Must(baseTemplate.Clone()).Parse(pageTemplateString)
		renderer.Add(page, pageTemplate)
	}

	return renderer
}
