package admin

import (
	"strconv"
	"task-system-go/internal/db"
	"task-system-go/internal/service"
	"task-system-go/internal/util"

	"github.com/gin-gonic/gin"
)

type WithdrawalHandler struct{}

func (h WithdrawalHandler) List(c *gin.Context) {
	statusStr := c.Query("status")
	keyword := c.Query("keyword")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	query := db.DB.Table("withdrawals w").
		Select("w.*, u.username").
		Joins("left join users u on w.user_id = u.id").Order("w.id desc")

	if statusStr != "" {
		status, _ := strconv.Atoi(statusStr)
		query = query.Where("w.status = ?", status)
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

func (h WithdrawalHandler) Pay(c *gin.Context) {
	adminID, _ := c.Get("admin_id")
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	admin := uint64(0)
	if v, ok := adminID.(float64); ok {
		admin = uint64(v)
	}
	if err := service.ApproveWithdrawal(id, admin); err != nil {
		util.JSONError(c, err.Error(), 1)
		return
	}
	util.JSONSuccess(c, gin.H{})
}

func (h WithdrawalHandler) Reject(c *gin.Context) {
	adminID, _ := c.Get("admin_id")
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	admin := uint64(0)
	if v, ok := adminID.(float64); ok {
		admin = uint64(v)
	}
	var req struct {
		Reason string `json:"reason"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || req.Reason == "" {
		util.JSONError(c, "请填写驳回原因", 1)
		return
	}
	if err := service.RejectWithdrawal(id, admin, req.Reason); err != nil {
		util.JSONError(c, err.Error(), 1)
		return
	}
	util.JSONSuccess(c, gin.H{})
}
