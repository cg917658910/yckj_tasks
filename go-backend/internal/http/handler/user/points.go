package user

import (
	"strconv"
	"task-system-go/internal/db"
	"task-system-go/internal/util"

	"github.com/gin-gonic/gin"
)

type PointsHandler struct{}

func (h PointsHandler) Logs(c *gin.Context) {
	userID, _ := c.Get("user_id")
	uid := uint64(0)
	if v, ok := userID.(float64); ok {
		uid = uint64(v)
	}
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	query := db.DB.Table("points_logs").Where("user_id = ?", uid).Order("id desc")
	var total int64
	_ = query.Count(&total).Error
	list := []map[string]interface{}{}
	_ = query.Offset((page-1)*pageSize).Limit(pageSize).Find(&list).Error

	util.JSONSuccess(c, gin.H{"list": list, "total": total, "page": page, "page_size": pageSize})
}
