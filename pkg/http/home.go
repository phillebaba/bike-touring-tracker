package http

import (
	"github.com/gin-gonic/gin"
	"net/http"

	"github.com/phillebaba/bike-tracker/pkg/domain"
)

type HomeHandler struct {
	TripService domain.TripService
}

func (h *HomeHandler) Index(c *gin.Context) {
	trip := h.TripService.List()[0]
	c.HTML(http.StatusOK, "index", gin.H{
		"Trip": trip,
	})
}
