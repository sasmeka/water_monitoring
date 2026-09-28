package routers

import (
	"github.com/sasmeka/water_monitoring/internal/handlers"
	"github.com/sasmeka/water_monitoring/internal/repositories"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func waterMonitoring(g *gin.Engine, base_endpoint string, d *sqlx.DB) {
	route := g.Group(base_endpoint + "/water-monitoring")

	// Dependency Injection
	repo := repositories.New_WaterMonitoring(d)
	handler := handlers.New_WaterMonitoring(repo)

	route.POST("", handler.Create)
	route.POST("/batch", handler.CreateBatch)
	route.GET("", handler.GetAll)
	route.GET("/:id", handler.GetByID)
}
