package handler

import (
	"github.com/EslRain/leaf-segment/service"
	"github.com/gin-gonic/gin"
	"log"
	"os"
)

type Handler struct {
	engine  *gin.Engine
	service *service.Service
}

func NewHandler(engine *gin.Engine, service *service.Service) *Handler {
	return &Handler{
		engine:  engine,
		service: service,
	}
}

func (h *Handler) Run() {

	app_host := os.Getenv("APP_HOST")
	h.RegRouter()
	err := h.engine.Run(app_host)
	if err != nil {
		log.Fatalln("server start failed")
	}
}

func (h *Handler) RegRouter() {
	//leaf
	r := h.engine.Group("api/leaf")
	{
		r.GET("", h.GetID)
		r.POST("", h.CreateLeaf)
	}
}
