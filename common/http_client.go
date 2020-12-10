package common

import "github.com/gin-gonic/gin"

func NewHttpClient() *gin.Engine {
	return gin.Default()
}
