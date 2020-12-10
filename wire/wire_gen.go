// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//+build !wireinject

package wire

import (
	"github.com/EslRain/leaf-segment/common"
	"github.com/EslRain/leaf-segment/dao"
	"github.com/EslRain/leaf-segment/handler"
	"github.com/EslRain/leaf-segment/service"
)

// Injectors from wire.go:

func InitHandler() *handler.Handler {
	engine := common.NewHttpClient()
	db := common.NewMysqlClient()
	daoDao := dao.NewDao(db)
	serviceService := service.NewService(daoDao)
	handlerHandler := handler.NewHandler(engine, serviceService)
	return handlerHandler
}