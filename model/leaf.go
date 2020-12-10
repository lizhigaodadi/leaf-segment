package model

import (
	"sync"
	"time"
)

type Leaf struct {
	ID         uint64 `json:"id" form:"id"`                   //主键ID
	BizTag     string `json:"biz_tag" form:"biz_tag"`         //区分业务
	MaxID      uint64 `json:"max_id" form:"max_id"`           //改biz_tag目前所被分配的ID号段的最大值
	Step       int32  `json:"step" form:"step"`               //每次分配ID号段长度
	UpdateTime uint64 `json:"update_time" form:"update_time"` //更新时间
}

func LeafTableName() string {
	return "leaf"
}

type LeafCache struct {
	Key        string
	Step       int32          //记录步长
	CurrentPos int32          //当前使用的segment buffer光标；总共两个buffer缓存区，交替使用
	Buffer     []*LeafSegment //双buffer，作为预缓冲使用
	UpdateTime time.Time
	IsPreload  bool
	Waiting    []chan byte //挂起等待
	Mutex      sync.Mutex
}

type LeafSegment struct {
	Cursor uint64 //当前发放位置
	Max    uint64 //最大值
	Min    uint64 //开始值，最小值
	InitOk bool   //是否初始化
}
