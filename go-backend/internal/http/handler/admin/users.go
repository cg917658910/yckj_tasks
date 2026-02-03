package admin

import (
	"strconv"
	"task-system-go/internal/db"
	"task-system-go/internal/service"
	"task-system-go/internal/util"

	"github.com/gin-gonic/gin"
)

type UserHandler struct{}

func (h UserHandler) List(c *gin.Context) {
	statusStr := c.Query("status")
	keyword := c.Query("keyword")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	query := db.DB.Table("users u").
		Select("u.id,u.username,u.status,u.created_at,p.available_points,p.frozen_points,p.withdrawn_points,up.total_withdrawn").
		Joins("left join points_accounts p on u.id = p.user_id").
		Joins("left join user_profiles up on u.id = up.user_id").
		Order("u.id desc")

	if statusStr != "" {
		status, _ := strconv.Atoi(statusStr)
		query = query.Where("u.status = ?", status)
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

func (h UserHandler) Status(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		Status int `json:"status"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		util.JSONError(c, "参数不完整", 1)
		return
	}
	if err := db.DB.Table("users").Where("id = ?", id).Update("status", req.Status).Error; err != nil {
		util.JSONError(c, "更新失败", 1)
		return
	}
	util.JSONSuccess(c, gin.H{})
}

func (h UserHandler) AdjustPoints(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var req struct {
		Points int    `json:"points"`
		Remark string `json:"remark"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || req.Points == 0 {
		util.JSONError(c, "积分调整值不能为0", 1)
		return
	}
	if err := service.AddAvailable(id, req.Points, "manual_adjust", nil, req.Remark); err != nil {
		util.JSONError(c, err.Error(), 1)
		return
	}
	util.JSONSuccess(c, gin.H{})
}

func (h UserHandler) Tasks(c *gin.Context) {
	id := c.Param("id")
	list := []map[string]interface{}{}
	_ = db.DB.Table("task_claims c").
		Select("c.*, t.title as task_title, t.reward_points").
		Joins("left join tasks t on c.task_id = t.id").
		Where("c.user_id = ?", id).Order("c.id desc").Find(&list).Error
	util.JSONSuccess(c, gin.H{"list": list})
}

func (h UserHandler) Withdrawals(c *gin.Context) {
	id := c.Param("id")
	list := []map[string]interface{}{}
	_ = db.DB.Table("withdrawals").Where("user_id = ?", id).Order("id desc").Find(&list).Error
	util.JSONSuccess(c, gin.H{"list": list})
}
