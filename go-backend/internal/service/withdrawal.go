package service

import (
	"errors"
	"time"

	"task-system-go/internal/db"
	"task-system-go/internal/model"

	"gorm.io/gorm"
)

const (
	WithdrawPending  = 1
	WithdrawPaid     = 2
	WithdrawRejected = 3
)

func ApplyWithdrawal(userID uint64, amount float64) error {
	return db.DB.Transaction(func(tx *gorm.DB) error {
		rule, err := GetPointsRule()
		if err != nil {
			return errors.New("积分规则未配置")
		}
		if amount < rule.MinWithdrawAmount {
			return errors.New("未达到最低提现金额")
		}

		var profile model.UserProfile
		if err := tx.Where("user_id = ?", userID).First(&profile).Error; err != nil {
			return errors.New("请先绑定微信收款二维码")
		}
		if profile.WechatQRURL == "" {
			return errors.New("请先绑定微信收款二维码")
		}

		pointsCost := int(amount * float64(rule.ExchangeRate))
		if err := FreezeForWithdraw(userID, pointsCost, "提现冻结"); err != nil {
			return err
		}

		withdrawal := model.Withdrawal{UserID: userID, Amount: amount, PointsCost: pointsCost, QRURL: profile.WechatQRURL, Status: WithdrawPending}
		return tx.Create(&withdrawal).Error
	})
}

func ApproveWithdrawal(withdrawalID uint64, adminID uint64) error {
	return db.DB.Transaction(func(tx *gorm.DB) error {
		var w model.Withdrawal
		if err := tx.Where("id = ?", withdrawalID).First(&w).Error; err != nil {
			return errors.New("提现不存在")
		}
		if w.Status != WithdrawPending {
			return errors.New("当前状态不可打款")
		}

		if err := DeductFrozenToWithdrawn(w.UserID, w.PointsCost, "提现完成"); err != nil {
			return err
		}

		now := time.Now()
		w.Status = WithdrawPaid
		w.ReviewedAt = &now
		w.AdminID = &adminID
		if err := tx.Save(&w).Error; err != nil {
			return err
		}

		return tx.Model(&model.UserProfile{}).Where("user_id = ?", w.UserID).UpdateColumn("total_withdrawn", gorm.Expr("total_withdrawn + ?", w.Amount)).Error
	})
}

func RejectWithdrawal(withdrawalID uint64, adminID uint64, reason string) error {
	return db.DB.Transaction(func(tx *gorm.DB) error {
		var w model.Withdrawal
		if err := tx.Where("id = ?", withdrawalID).First(&w).Error; err != nil {
			return errors.New("提现不存在")
		}
		if w.Status != WithdrawPending {
			return errors.New("当前状态不可驳回")
		}

		if err := UnfreezeToAvailable(w.UserID, w.PointsCost, "提现驳回"); err != nil {
			return err
		}

		now := time.Now()
		w.Status = WithdrawRejected
		w.ReviewedAt = &now
		w.AdminID = &adminID
		w.RejectReason = reason
		return tx.Save(&w).Error
	})
}
