package main

import (
	"log"
	"os"

	"github.com/sasmeka/water_monitoring/docs"
	"github.com/sasmeka/water_monitoring/internal/routers"
	"github.com/sasmeka/water_monitoring/pkg"

	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"

	_ "github.com/sasmeka/water_monitoring/docs"

	_ "github.com/joho/godotenv/autoload"
)

func init() {
	govalidator.SetFieldsRequiredByDefault(true)
	gin.SetMode(os.Getenv("GIN_MODE"))
}

// @title						Water Monitoring Backend
// @version						1.0
// @BasePath					/water-monitoring/auth
// @securityDefinitions.apikey 	BearerAuth
// @in 							header
// @name 						Authorization

//	@tag.name			Auth
// 	@tag.name			Profile

func main() {
	if os.Getenv("GIN_MODE") != "release" {
		docs.SwaggerInfo.BasePath = "/water-monitoring/auth-dev"
	}

	router := routers.Routers(docs.SwaggerInfo.BasePath)
	server := pkg.Server(router)

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

/// nodemon --watch './**/*.{go,yaml,json}' --signal SIGTERM --exec 'go' run ./cmd/main.go
/// swag init -g ./cmd/main.go --md ./docs/md -o ./docs
/// export PATH=$(go env GOPATH)/bin:$PATH
/// go build -o crm-mobile-backend ./cmd/main.go
