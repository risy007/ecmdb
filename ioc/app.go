package ioc

import (
	"github.com/Duke1616/ecmdb/internal/event/easyflow"
	"github.com/gin-gonic/gin"
)

type App struct {
	Web   *gin.Engine
	Event *easyflow.ProcessEvent
}
