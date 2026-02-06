package service

import (
	"errors"
	"time"

	"task-system-go/internal/db"
	"task-system-go/internal/model"

	"gorm.io/gorm"
)

const (
	TaskStatusDraft   = 0
	TaskStatusOnline  = 1
	TaskStatusOffline = 2

	ClaimStatusClaimed   = 1
	ClaimStatusSubmitted = 2
	ClaimStatusApproved  = 3
	ClaimStatusRejected  = 4
)

func ClaimTask(userID, taskID uint64) error {
	return db.DB.Transaction(func(tx *gorm.DB) error {
		var active model.TaskClaim
		err := tx.Where("user_id = ? AND status IN (1,2)", userID).First(&active).Error
		if err == nil {
			return errors.New("存在未完成任务，无法领取新任务")
		}

		var task model.Task
		if err := tx.Where("id = ?", taskID).First(&task).Error; err != nil {
			return errors.New("任务不可领取")
		}
		if task.Status != TaskStatusOnline {
			return errors.New("任务不可领取")
		}

		var exists model.TaskClaim
		if err := tx.Where("task_id = ?", taskID).First(&exists).Error; err == nil {
			return errors.New("任务已被领取")
		}

		claim := model.TaskClaim{
			TaskID:            taskID,
			UserID:            userID,
			Status:            ClaimStatusClaimed,
			ClaimedAt:         time.Now(),
			SubmittedAt:       nil,
			ReviewedAt:        nil,
			ReviewResult:      nil,
			RejectReason:      "",
			RewardPointsFinal: nil,
		}
		return tx.Create(&claim).Error
	})
}

func SubmitTask(userID, claimID uint64, remark string, images []string) error {
	return db.DB.Transaction(func(tx *gorm.DB) error {
		var claim model.TaskClaim
		if err := tx.Where("id = ?", claimID).First(&claim).Error; err != nil {
			return errors.New("任务不存在")
		}
		if claim.UserID != userID {
			return errors.New("无权限")
		}
		if claim.Status != ClaimStatusClaimed && claim.Status != ClaimStatusRejected {
			return errors.New("当前状态不可提交")
		}

		var submission model.TaskSubmission
		if err := tx.Where("claim_id = ?", claimID).First(&submission).Error; err == nil {
			submission.Remark = remark
			if err := tx.Save(&submission).Error; err != nil {
				return err
			}
		} else {
			submission = model.TaskSubmission{ClaimID: claimID, Remark: remark}
			if err := tx.Create(&submission).Error; err != nil {
				return err
			}
		}

		_ = tx.Where("submission_id = ?", submission.ID).Delete(&model.TaskSubmissionImage{}).Error
		for _, url := range images {
			img := model.TaskSubmissionImage{SubmissionID: submission.ID, ImageURL: url}
			if err := tx.Create(&img).Error; err != nil {
				return err
			}
		}

		now := time.Now()
		claim.Status = ClaimStatusSubmitted
		claim.SubmittedAt = &now
		return tx.Save(&claim).Error
	})
}

func ApproveClaim(claimID uint64, rewardPoints int) error {
	return db.DB.Transaction(func(tx *gorm.DB) error {
		var claim model.TaskClaim
		if err := tx.Where("id = ?", claimID).First(&claim).Error; err != nil {
			return errors.New("任务不存在")
		}
		if claim.Status != ClaimStatusSubmitted {
			return errors.New("当前状态不可审核通过")
		}

		var task model.Task
		_ = tx.Where("id = ?", claim.TaskID).First(&task).Error
		final := rewardPoints
		if final <= 0 {
			final = task.RewardPoints
		}

		one := 1
		now := time.Now()
		claim.Status = ClaimStatusApproved
		claim.ReviewResult = &one
		claim.ReviewedAt = &now
		claim.RewardPointsFinal = &final
		if err := tx.Save(&claim).Error; err != nil {
			return err
		}

		refID := claim.ID
		return AddAvailable(claim.UserID, final, "task_reward", &refID, "任务完成奖励")
	})
}

func RejectClaim(claimID uint64, reason string) error {
	return db.DB.Transaction(func(tx *gorm.DB) error {
		var claim model.TaskClaim
		if err := tx.Where("id = ?", claimID).First(&claim).Error; err != nil {
			return errors.New("任务不存在")
		}
		if claim.Status != ClaimStatusSubmitted {
			return errors.New("当前状态不可驳回")
		}

		zero := 0
		now := time.Now()
		claim.Status = ClaimStatusRejected
		claim.ReviewResult = &zero
		claim.ReviewedAt = &now
		claim.RejectReason = reason
		return tx.Save(&claim).Error
	})
}
