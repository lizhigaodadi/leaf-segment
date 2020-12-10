package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) GetID(ctx *gin.Context) {
	bizTag := ctx.Query("biz_tag")
	id, err := h.service.GetID(ctx, bizTag)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"code": 1001, "msg": err.Error()})
		return
	}
	resp := GetIdReq{
		ID: id,
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "msg": "success", "data": resp})
}

func (h *Handler) CreateLeaf(ctx *gin.Context) {
	req := &CreateLeafReq{}
	if err := ctx.ShouldBindJSON(req); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"code": 1000, "msg": "param Invalid"})
		return
	}

	err := h.service.CreateLeaf(ctx, req.toCreate())
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"code": 1001, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "msg": "success", "data": nil})
}
