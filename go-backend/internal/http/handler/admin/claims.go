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

	// append submission images
	if len(list) > 0 {
		claimIDs := make([]uint64, 0, len(list))
		for _, item := range list {
			if v, ok := item["id"]; ok {
				switch t := v.(type) {
				case int64:
					claimIDs = append(claimIDs, uint64(t))
				case uint64:
					claimIDs = append(claimIDs, t)
				case float64:
					claimIDs = append(claimIDs, uint64(t))
				}
			}
		}

		subMap := map[uint64]uint64{}
		if len(claimIDs) > 0 {
			rows := []map[string]interface{}{}
			_ = db.DB.Table("task_submissions").Select("id, claim_id").Where("claim_id IN ?", claimIDs).Find(&rows).Error
			for _, row := range rows {
				var claimID uint64
				var subID uint64
				switch v := row["claim_id"].(type) {
				case int64:
					claimID = uint64(v)
				case uint64:
					claimID = v
				case float64:
					claimID = uint64(v)
				}
				switch v := row["id"].(type) {
				case int64:
					subID = uint64(v)
				case uint64:
					subID = v
				case float64:
					subID = uint64(v)
				}
				if claimID > 0 && subID > 0 {
					subMap[claimID] = subID
				}
			}
		}

		subIDs := make([]uint64, 0, len(subMap))
		for _, v := range subMap {
			subIDs = append(subIDs, v)
		}
		imagesMap := map[uint64][]string{}
		if len(subIDs) > 0 {
			imgRows := []map[string]interface{}{}
			_ = db.DB.Table("task_submission_images").Select("submission_id, image_url").Where("submission_id IN ?", subIDs).Find(&imgRows).Error
			for _, row := range imgRows {
				var subID uint64
				switch v := row["submission_id"].(type) {
				case int64:
					subID = uint64(v)
				case uint64:
					subID = v
				case float64:
					subID = uint64(v)
				}
				if subID == 0 {
					continue
				}
				if url, ok := row["image_url"].(string); ok {
					imagesMap[subID] = append(imagesMap[subID], url)
				}
			}
		}

		for _, item := range list {
			var claimID uint64
			switch v := item["id"].(type) {
			case int64:
				claimID = uint64(v)
			case uint64:
				claimID = v
			case float64:
				claimID = uint64(v)
			}
			subID := subMap[claimID]
			item["images"] = imagesMap[subID]
		}
	}

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
