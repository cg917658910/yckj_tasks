package common

import (
	"task-system-go/internal/config"
	"task-system-go/internal/service"
	"task-system-go/internal/util"

	"github.com/gin-gonic/gin"
)

type UploadHandler struct {
	Cfg config.Config
}

func (h UploadHandler) Image(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		util.JSONError(c, "请选择文件", 1)
		return
	}
	data, err := service.SaveImage(h.Cfg, file)
	if err != nil {
		util.JSONError(c, err.Error(), 1)
		return
	}
	util.JSONSuccess(c, data)
}
