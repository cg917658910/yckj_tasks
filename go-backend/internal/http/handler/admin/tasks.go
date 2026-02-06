package admin

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
	statusStr := c.Query("status")
	keyword := c.Query("keyword")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	query := db.DB.Table("tasks t").
		Select("t.*, c.id as claimed").
		Joins("left join task_claims c on t.id = c.task_id").
		Order("t.id desc")
	if statusStr != "" {
		status, _ := strconv.Atoi(statusStr)
		query = query.Where("t.status = ?", status)
	}
	if keyword != "" {
		query = query.Where("t.title LIKE ?", "%"+keyword+"%")
	}

	var total int64
	_ = query.Count(&total).Error
	var list []map[string]interface{}
	_ = query.Offset((page-1)*pageSize).Limit(pageSize).Find(&list).Error

	util.JSONSuccess(c, gin.H{"list": list, "total": total, "page": page, "page_size": pageSize})
}

func (h TaskHandler) Create(c *gin.Context) {
	var req model.Task
	if err := c.ShouldBindJSON(&req); err != nil {
		util.JSONError(c, "参数不完整", 1)
		return
	}
	req.Status = service.TaskStatusOnline
	if err := db.DB.Create(&req).Error; err != nil {
		util.JSONError(c, "创建失败", 1)
		return
	}
	util.JSONSuccess(c, req)
}

func (h TaskHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var task model.Task
	if err := db.DB.Where("id = ?", id).First(&task).Error; err != nil {
		util.JSONError(c, "任务不存在", 1)
		return
	}

	var req model.Task
	_ = c.ShouldBindJSON(&req)

	task.Title = req.Title
	task.Summary = req.Summary
	task.Detail = req.Detail
	task.DocURL = req.DocURL
	task.RewardPoints = req.RewardPoints
	if req.Status != 0 {
		task.Status = req.Status
	}

	if err := db.DB.Save(&task).Error; err != nil {
		util.JSONError(c, "更新失败", 1)
		return
	}
	util.JSONSuccess(c, task)
}

func (h TaskHandler) Off(c *gin.Context) {
	id := c.Param("id")
	if err := db.DB.Model(&model.Task{}).Where("id = ?", id).Update("status", service.TaskStatusOffline).Error; err != nil {
		util.JSONError(c, "下架失败", 1)
		return
	}
	util.JSONSuccess(c, gin.H{})
}

func (h TaskHandler) On(c *gin.Context) {
	id := c.Param("id")
	if err := db.DB.Model(&model.Task{}).Where("id = ?", id).Update("status", service.TaskStatusOnline).Error; err != nil {
		util.JSONError(c, "上架失败", 1)
		return
	}
	util.JSONSuccess(c, gin.H{})
}
