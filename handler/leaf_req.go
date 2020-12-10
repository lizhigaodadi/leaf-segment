package handler

import "github.com/EslRain/leaf-segment/model"

type CreateLeafReq struct {
	BizTag      string  `json:"biz_tag"`
	MaxID       *uint64 `json:"max_id"` // 可以不传 默认为1
	Step        *int32  `json:"step"`   // 可以不传 默认为2000
	Description string  `json:"description"`
}

type GetIdReq struct {
	ID uint64
}

func (c *CreateLeafReq) toCreate() *model.Leaf {
	leaf := &model.Leaf{}
	if c.MaxID == nil {
		leaf.MaxID = 1
	} else {
		leaf.MaxID = *c.MaxID
	}
	if c.Step == nil {
		leaf.Step = 2000
	}
	leaf.BizTag = c.BizTag
	return leaf
}
