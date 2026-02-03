package user

import (
	"task-system-go/internal/db"
	"task-system-go/internal/service"
	"task-system-go/internal/util"

	"github.com/gin-gonic/gin"
)

type WithdrawalHandler struct{}

func (h WithdrawalHandler) Apply(c *gin.Context) {
	userID, _ := c.Get("user_id")
	uid := uint64(0)
	if v, ok := userID.(float64); ok {
		uid = uint64(v)
	}
	var req struct {
		Amount float64 `json:"amount"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || req.Amount <= 0 {
		util.JSONError(c, "金额不合法", 1)
		return
	}
	if err := service.ApplyWithdrawal(uid, req.Amount); err != nil {
		util.JSONError(c, err.Error(), 1)
		return
	}
	util.JSONSuccess(c, gin.H{})
}

func (h WithdrawalHandler) List(c *gin.Context) {
	userID, _ := c.Get("user_id")
	uid := uint64(0)
	if v, ok := userID.(float64); ok {
		uid = uint64(v)
	}
	list := []map[string]interface{}{}
	_ = db.DB.Table("withdrawals").Where("user_id = ?", uid).Order("id desc").Find(&list).Error
	util.JSONSuccess(c, gin.H{"list": list})
}
