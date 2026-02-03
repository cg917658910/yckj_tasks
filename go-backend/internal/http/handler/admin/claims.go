package admin

import (
	"strconv"
	"task-system-go/internal/db"
	"task-system-go/internal/service"
	"task-system-go/internal/util"

	"github.com/gin-gonic/gin"
)

type ClaimHandler struct{}

func (h ClaimHandler) List(c *gin.Context) {
	statusStr := c.Query("status")
	keyword := c.Query("keyword")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	query := db.DB.Table("task_claims c").
		Select("c.*, t.title as task_title, t.reward_points, u.username").
		Joins("left join tasks t on c.task_id = t.id").
		Joins("left join users u on c.user_id = u.id").Order("c.id desc")

	if statusStr != "" {
		status, _ := strconv.Atoi(statusStr)
		query = query.Where("c.status = ?", status)
	}
	if keyword != "" {
		query = query.Where("t.title LIKE ? OR u.username LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	var total int64
	_ = query.Count(&total).Error
	var list []map[string]interface{}
	_ = query.Offset((page - 1) * pageSize).Limit(pageSize).Find(&list).Error

	util.JSONSuccess(c, gin.H{"list": list, "total": total, "page": page, "page_size": pageSize})
}

func (h ClaimHandler) Approve(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var req struct {
		RewardPoints int `json:"reward_points"`
	}
	_ = c.ShouldBindJSON(&req)

	if err := service.ApproveClaim(id, req.RewardPoints); err != nil {
		util.JSONError(c, err.Error(), 1)
		return
	}
	util.JSONSuccess(c, gin.H{})
}

func (h ClaimHandler) Reject(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var req struct {
		Reason string `json:"reason"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || req.Reason == "" {
		util.JSONError(c, "请填写驳回原因", 1)
		return
	}
	if err := service.RejectClaim(id, req.Reason); err != nil {
		util.JSONError(c, err.Error(), 1)
		return
	}
	util.JSONSuccess(c, gin.H{})
}
