package user

import (
	"strconv"
	"task-system-go/internal/db"
	"task-system-go/internal/model"
	"task-system-go/internal/service"
	"task-system-go/internal/util"

	"github.com/gin-gonic/gin"
)

type TaskHandler struct{}

func (h TaskHandler) List(c *gin.Context) {
	keyword := c.Query("keyword")
	query := db.DB.Table("tasks t").
		Select("t.*, c.id as claimed").
		Joins("left join task_claims c on t.id = c.task_id").
		Where("t.status = ?", service.TaskStatusOnline)
	if keyword != "" {
		query = query.Where("t.title LIKE ?", "%"+keyword+"%")
	}
	var list []map[string]interface{}
	_ = query.Order("t.id desc").Find(&list).Error
	util.JSONSuccess(c, gin.H{"list": list})
}

func (h TaskHandler) Detail(c *gin.Context) {
	id := c.Param("id")
	var task model.Task
	if err := db.DB.Where("id = ?", id).First(&task).Error; err != nil {
		util.JSONError(c, "任务不存在", 1)
		return
	}
	var claim model.TaskClaim
	claimed := 0
	if err := db.DB.Where("task_id = ?", id).First(&claim).Error; err == nil {
		claimed = 1
	}
	data := gin.H{
		"id":            task.ID,
		"title":         task.Title,
		"summary":       task.Summary,
		"detail":        task.Detail,
		"doc_url":       task.DocURL,
		"reward_points": task.RewardPoints,
		"claimed":       claimed,
	}
	util.JSONSuccess(c, data)
}

func (h TaskHandler) Claim(c *gin.Context) {
	userID, _ := c.Get("user_id")
	uid := uint64(0)
	if v, ok := userID.(float64); ok {
		uid = uint64(v)
	}
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if err := service.ClaimTask(uid, id); err != nil {
		util.JSONError(c, err.Error(), 1)
		return
	}
	util.JSONSuccess(c, gin.H{})
}
