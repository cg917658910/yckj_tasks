package user

import (
	"strconv"
	"task-system-go/internal/db"
	"task-system-go/internal/service"
	"task-system-go/internal/util"

	"github.com/gin-gonic/gin"
)

type ClaimHandler struct{}

func (h ClaimHandler) Current(c *gin.Context) {
	userID, _ := c.Get("user_id")
	uid := uint64(0)
	if v, ok := userID.(float64); ok {
		uid = uint64(v)
	}

	item := make(map[string]interface{})
	_ = db.DB.Table("task_claims c").
		Select("c.*, t.title as task_title, t.summary, t.detail, t.doc_url, t.reward_points").
		Joins("left join tasks t on c.task_id = t.id").
		Where("c.user_id = ? AND c.status IN (1,2,4)", uid).
		Order("c.id desc").Scan(&item).Error

	util.JSONSuccess(c, gin.H{"current": item})
}

func (h ClaimHandler) History(c *gin.Context) {
	userID, _ := c.Get("user_id")
	uid := uint64(0)
	if v, ok := userID.(float64); ok {
		uid = uint64(v)
	}
	list := []map[string]interface{}{}
	_ = db.DB.Table("task_claims c").
		Select("c.*, t.title as task_title, t.reward_points").
		Joins("left join tasks t on c.task_id = t.id").
		Where("c.user_id = ?", uid).Order("c.id desc").Find(&list).Error
	util.JSONSuccess(c, gin.H{"list": list})
}

func (h ClaimHandler) Submit(c *gin.Context) {
	userID, _ := c.Get("user_id")
	uid := uint64(0)
	if v, ok := userID.(float64); ok {
		uid = uint64(v)
	}
	claimID, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var req struct {
		Images []string `json:"images"`
		Remark string   `json:"remark"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || len(req.Images) == 0 {
		util.JSONError(c, "请上传截图", 1)
		return
	}
	if err := service.SubmitTask(uid, claimID, req.Remark, req.Images); err != nil {
		util.JSONError(c, err.Error(), 1)
		return
	}
	util.JSONSuccess(c, gin.H{})
}
