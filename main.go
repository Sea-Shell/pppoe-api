// Some comment
package main

import (
	"database/sql"
	"flag"
	"log"

	utils "github.com/Sea-Shell/gogear-api/pkg/utils"
	endpoints "github.com/bateau84/pppoe-api/pkg/api"
	docs "github.com/bateau84/pppoe-api/pkg/docs"
	models "github.com/bateau84/pppoe-api/pkg/models"

	requestid "github.com/gin-contrib/requestid"
	gin "github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	zap "go.uber.org/zap"
	zapcore "go.uber.org/zap/zapcore"
)

const (
	listenPort     = "8081"
	configFile     = "config.yaml"
	googleCredFile = "google-creds.json"
)

func makeLogger(loglevel zapcore.Level) *zap.SugaredLogger {
	customCallerEncoder := func(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(caller.TrimmedPath())
	}
	cfg := zap.Config{
		Level:       zap.NewAtomicLevelAt(loglevel),
		Encoding:    "json",
		OutputPaths: []string{"stdout"},
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:      "time",
			EncodeLevel:  zapcore.LowercaseLevelEncoder,
			EncodeTime:   zapcore.ISO8601TimeEncoder,
			LevelKey:     "level",
			MessageKey:   "message",
			CallerKey:    "caller",
			EncodeCaller: customCallerEncoder,
		},
	}

	return zap.Must(cfg.Build()).Sugar()
}

func logRequestsMiddleware(logger *zap.SugaredLogger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Log the request details
		logger.Infow("Received request",
			"method", c.Request.Method,
			"status", c.Writer.Status(),
			"url", c.FullPath(),
			"url-params", c.Request.URL.Query(),
			"X-API-Key", c.Request.Header.Get("X-API-Key"),
		)

		// Continue processing the request
		c.Set("logger", logger)
		c.Next()
	}
}

func databaseMiddleware(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("db", db)
		c.Next()
	}
}

func configMiddleware(config *models.General) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("conf", config)
		c.Next()
	}
}

//	@title			PPPoE API
//	@version		1.0
//	@description	This is the API of PPPoE
//	@contact.name	API Support
//	@contact.email	support@sea-shell.no
//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

// @securityDefinitions.apikey				APIKey
// @in										header
// @name									X-API-Key
// @securitydefinitions.oauth2.password	OAuth2Application
// @description							OAuth protects our entity endpoints
// @tokenUrl								https://oauth2.googleapis.com/token
// @authorizationurl						https://accounts.google.com/o/oauth2/auth
// @scope.write							Grants read and write access
// @scope.admin							Grants read and write access to administrative information
// @scope.read								Grants read access
func main() {
	configFile := flag.String("config", configFile, "Config file")

	flag.Parse()

	config, err := utils.LoadConfig[models.Config](*configFile)
	if err != nil {
		log.Fatalf("Failed to load config file: %v", err)
	}

	logLevel := utils.GetLogLevel(config.General.LogLevel)
	log := makeLogger(logLevel)
	defer log.Sync()

	log.Debugf("%#v", config)

	db, err := sql.Open("sqlite3", config.Database.File)
	if err != nil {
		log.Error(err)
	}

	log.Infoln("Connected to database")
	defer db.Close()

	docs.SwaggerInfo.Title = "PPPoE API"
	docs.SwaggerInfo.Description = "This is the API of PPPoE."
	docs.SwaggerInfo.Host = config.General.Hostname
	docs.SwaggerInfo.Schemes = config.General.Schemes
	docs.SwaggerInfo.BasePath = "/api/v1"

	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Recovery())

	router.Use(logRequestsMiddleware(log))
	router.Use(databaseMiddleware(db))
	router.Use(configMiddleware(&config.General))
	router.Use(requestid.New())

	// API v1
	swagger := router.Group("/swagger")
	v1 := router.Group("/api/v1")

	// API Groups
	events := v1.Group("/event")

	// The routes
	v1.GET("/health", endpoints.ReturnHealth)

	// Event endpoints
	events.GET("/list", endpoints.ListEvents)
	events.GET("/:event/get", endpoints.GetEvent)
	// events.POST("/:gear/update", endpoints.UpdateEvent)
	// events.DELETE("/:gear/delete", endpoints.DeleteEvent)
	events.PUT("/insert", endpoints.InsertEvent)

	// Swagger API documentation
	swagger.GET("/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Listen to all addresses and port defined
	if err := router.Run("0.0.0.0:" + config.General.ListenPort); err != nil {
		log.Fatal(err)
	}
}
