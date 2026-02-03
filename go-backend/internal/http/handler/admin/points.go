package admin

import (
	"strconv"
	"task-system-go/internal/db"
	"task-system-go/internal/model"
	"task-system-go/internal/service"
	"task-system-go/internal/util"

	"github.com/gin-gonic/gin"
)

type PointsHandler struct{}

func (h PointsHandler) Rules(c *gin.Context) {
	rule, err := service.GetPointsRule()
	if err != nil {
		util.JSONError(c, err.Error(), 1)
		return
	}
	util.JSONSuccess(c, rule)
}

func (h PointsHandler) UpdateRules(c *gin.Context) {
	adminID, _ := c.Get("admin_id")
	var req model.PointsRule
	if err := c.ShouldBindJSON(&req); err != nil {
		util.JSONError(c, "参数不完整", 1)
		return
	}
	id, _ := adminID.(float64)
	rule, err := service.UpdatePointsRule(uint64(id), req)
	if err != nil {
		util.JSONError(c, err.Error(), 1)
		return
	}
	util.JSONSuccess(c, rule)
}

func (h PointsHandler) Logs(c *gin.Context) {
	userID := c.Query("user_id")
	typeStr := c.Query("type")
	keyword := c.Query("keyword")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	query := db.DB.Table("points_logs l").
		Select("l.*, u.username").
		Joins("left join users u on l.user_id = u.id").Order("l.id desc")
	if userID != "" {
		query = query.Where("l.user_id = ?", userID)
	}
	if typeStr != "" {
		query = query.Where("l.type = ?", typeStr)
	}
	if keyword != "" {
		query = query.Where("u.username LIKE ?", "%"+keyword+"%")
	}

	var total int64
	_ = query.Count(&total).Error
	var list []map[string]interface{}
	_ = query.Offset((page-1)*pageSize).Limit(pageSize).Find(&list).Error

	util.JSONSuccess(c, gin.H{"list": list, "total": total, "page": page, "page_size": pageSize})
}
