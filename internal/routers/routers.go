package routers

import (
	"net/http"
	"strings"

	"github.com/sasmeka/water_monitoring/config"
	"github.com/sasmeka/water_monitoring/pkg"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Routers(base_endpoint string) *gin.Engine {
	router := gin.Default()
	router.RedirectTrailingSlash = true
	// router.Use(cors.Default())
	router.Use(cors.New(config.CorsConfig))
	// router.Use(middleware.CORSMiddleware)

	router.Use(func(c *gin.Context) {
		if strings.Contains(c.Request.URL.Path, "/docs/index.html") {
			c.Redirect(http.StatusMovedPermanently, base_endpoint+"/documentation/")
		}
	})
	router.GET(base_endpoint+"/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, func(c *ginSwagger.Config) { c.Title = "Api Water Monitoring Backend" }))
	router.StaticFS(base_endpoint+"/documentation", http.Dir("./rapidoc"))
	router.GET(base_endpoint+base_endpoint+"/documentation/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, base_endpoint+"/documentation/")
	})

	pgdb := pkg.Postgres_Database("water_monitoring")
	// rdb := pkg.Redis("REDISDB")
	waterMonitoring(router, base_endpoint, pgdb)

	return router
}
