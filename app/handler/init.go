package handler

import (
	"clip/app/entity"
	"clip/config"
	"clip/repo"

	gintemplate "github.com/foolin/gin-template"
	"github.com/gin-gonic/gin"
)

type Presenter struct {
	VideoSrc  string
	VideoName string
	Flashes   string
}

func InitEngine(conf *config.Config) *gin.Engine {
	if conf.App.Debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.New()
	r.Use(gin.Recovery())

	if conf.App.Debug {
		r.Use(gin.Logger())
	}

	// ----------------------   INIT STATIC
	r.Static("vids", "./vids")

	templateConfig := gintemplate.TemplateConfig{
		Root:      "public",
		Extension: ".html",
	}

	r.HTMLRender = gintemplate.New(templateConfig)

	// ----------------------   INIT HANDLER

	clipHandler := clipHandler{
		clip: entity.NewClip(repo.Clip),
	}

	// ----------------------   INIT ROUTE

	groupIndex := r.Group("")
	groupIndex.GET("/", clipHandler.New)
	groupIndex.GET("/s", clipHandler.Watch)
	groupIndex.GET("/s/:slug", clipHandler.Watch)
	groupIndex.POST("/cut", clipHandler.Cut)

	return r
}

func GET(group *gin.RouterGroup, relativePath string, f func(*gin.Context)) {
	group.GET(relativePath, f)
}

func POST(group *gin.RouterGroup, relativePath string, f func(*gin.Context)) {
	group.POST(relativePath, f)
}

func PUT(group *gin.RouterGroup, relativePath string, f func(*gin.Context)) {
	group.PUT(relativePath, f)
}

func DELETE(group *gin.RouterGroup, relativePath string, f func(*gin.Context)) {
	group.DELETE(relativePath, f)
}
