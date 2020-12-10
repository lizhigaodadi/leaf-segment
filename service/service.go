package service

import (
	"github.com/EslRain/leaf-segment/dao"
	"sync"
)

type Service struct {
	dao       *dao.Dao
	leafCache sync.Map
}

func NewService(dao *dao.Dao) *Service {
	return &Service{
		dao: dao,
	}
}
