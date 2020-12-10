//+build wireinject

package wire

import (
	"github.com/EslRain/leaf-segment/common"
	"github.com/EslRain/leaf-segment/dao"
	"github.com/EslRain/leaf-segment/handler"
	"github.com/EslRain/leaf-segment/service"
	"github.com/google/wire"
)

func InitHandler() *handler.Handler {
	wire.Build(
		common.NewMysqlClient,
		common.NewHttpClient,

		dao.NewDao,
		service.NewService,
		handler.NewHandler,
	)
	return &handler.Handler{}
}
