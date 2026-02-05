package user

import (
	"fmt"
	"task-system-go/internal/db"
	"task-system-go/internal/model"
	"task-system-go/internal/util"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type ProfileHandler struct{}

func (h ProfileHandler) Info(c *gin.Context) {
	userID, _ := c.Get("user_id")
	uid := uint64(0)
	if v, ok := userID.(float64); ok {
		uid = uint64(v)
	}
	fmt.Printf("Fetching profile for user ID: %d\n", uid)
	data := make(map[string]interface{})
	if err := db.DB.Table("users u").
		Select("u.id,u.username,u.created_at,p.available_points,p.frozen_points,p.withdrawn_points,up.wechat_qr_url,up.total_withdrawn").
		Joins("left join points_accounts p on u.id = p.user_id").
		Joins("left join user_profiles up on u.id = up.user_id").
		Where("u.id = ?", uid).Scan(&data).Error; err != nil {
		util.JSONError(c, "用户不存在", 1)
		return
	}

	util.JSONSuccess(c, data)
}

func (h ProfileHandler) ChangePassword(c *gin.Context) {
	userID, _ := c.Get("user_id")
	uid := uint64(0)
	if v, ok := userID.(float64); ok {
		uid = uint64(v)
	}
	var req struct {
		OldPassword string `json:"old_password"`
		NewPassword string `json:"new_password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || req.OldPassword == "" || req.NewPassword == "" {
		util.JSONError(c, "参数不完整", 1)
		return
	}

	var user model.User
	if err := db.DB.Where("id = ?", uid).First(&user).Error; err != nil {
		util.JSONError(c, "用户不存在", 1)
		return
	}
	if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.OldPassword)) != nil {
		util.JSONError(c, "原密码错误", 1)
		return
	}

	hash, _ := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err := db.DB.Model(&user).Update("password_hash", string(hash)).Error; err != nil {
		util.JSONError(c, "更新失败", 1)
		return
	}
	util.JSONSuccess(c, gin.H{})
}

func (h ProfileHandler) UpdateWechatQR(c *gin.Context) {
	userID, _ := c.Get("user_id")
	uid := uint64(0)
	if v, ok := userID.(float64); ok {
		uid = uint64(v)
	}
	var req struct {
		WechatQRURL string `json:"wechat_qr_url"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || req.WechatQRURL == "" {
		util.JSONError(c, "请上传二维码", 1)
		return
	}
	if err := db.DB.Table("user_profiles").Where("user_id = ?", uid).Update("wechat_qr_url", req.WechatQRURL).Error; err != nil {
		util.JSONError(c, "更新失败", 1)
		return
	}
	util.JSONSuccess(c, gin.H{})
}
