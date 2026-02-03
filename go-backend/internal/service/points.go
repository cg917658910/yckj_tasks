package service

import (
	"errors"
	"task-system-go/internal/db"
	"task-system-go/internal/model"

	"gorm.io/gorm"
)

func GetPointsRule() (model.PointsRule, error) {
	var rule model.PointsRule
	err := db.DB.Order("id desc").First(&rule).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			rule = model.PointsRule{ExchangeRate: 10, MinWithdrawAmount: 10, RegisterBonusPoints: 10}
			_ = db.DB.Create(&rule).Error
			return rule, nil
		}
		return rule, err
	}
	return rule, nil
}

func UpdatePointsRule(adminID uint64, data model.PointsRule) (model.PointsRule, error) {
	rule, err := GetPointsRule()
	if err != nil {
		return rule, err
	}

	rule.ExchangeRate = data.ExchangeRate
	rule.MinWithdrawAmount = data.MinWithdrawAmount
	rule.RegisterBonusPoints = data.RegisterBonusPoints

	if err := db.DB.Save(&rule).Error; err != nil {
		return rule, err
	}

	_ = db.DB.Table("points_rule_logs").Create(map[string]interface{}{
		"rule_id":   rule.ID,
		"old_value": "{}",
		"new_value": "{}",
		"admin_id":  adminID,
	}).Error

	return rule, nil
}

func AddAvailable(userID uint64, points int, typ string, refID *uint64, remark string) error {
	return db.DB.Transaction(func(tx *gorm.DB) error {
		var acc model.PointsAccount
		err := tx.Where("user_id = ?", userID).First(&acc).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				acc = model.PointsAccount{UserID: userID, AvailablePoints: 0, FrozenPoints: 0, WithdrawnPoints: 0}
				if err := tx.Create(&acc).Error; err != nil {
					return err
				}
			} else {
				return err
			}
		}

		acc.AvailablePoints += points
		if err := tx.Save(&acc).Error; err != nil {
			return err
		}

		log := model.PointsLog{UserID: userID, ChangePoints: points, Type: typ, RefID: refID, Remark: remark}
		return tx.Create(&log).Error
	})
}

func FreezeForWithdraw(userID uint64, points int, remark string) error {
	return db.DB.Transaction(func(tx *gorm.DB) error {
		var acc model.PointsAccount
		if err := tx.Where("user_id = ?", userID).First(&acc).Error; err != nil {
			return err
		}
		if acc.AvailablePoints < points {
			return errors.New("积分不足")
		}
		acc.AvailablePoints -= points
		acc.FrozenPoints += points
		if err := tx.Save(&acc).Error; err != nil {
			return err
		}
		log := model.PointsLog{UserID: userID, ChangePoints: -points, Type: "withdraw_freeze", Remark: remark}
		return tx.Create(&log).Error
	})
}

func UnfreezeToAvailable(userID uint64, points int, remark string) error {
	return db.DB.Transaction(func(tx *gorm.DB) error {
		var acc model.PointsAccount
		if err := tx.Where("user_id = ?", userID).First(&acc).Error; err != nil {
			return err
		}
		if acc.FrozenPoints < points {
			return errors.New("冻结积分不足")
		}
		acc.FrozenPoints -= points
		acc.AvailablePoints += points
		if err := tx.Save(&acc).Error; err != nil {
			return err
		}
		log := model.PointsLog{UserID: userID, ChangePoints: points, Type: "withdraw_unfreeze", Remark: remark}
		return tx.Create(&log).Error
	})
}

func DeductFrozenToWithdrawn(userID uint64, points int, remark string) error {
	return db.DB.Transaction(func(tx *gorm.DB) error {
		var acc model.PointsAccount
		if err := tx.Where("user_id = ?", userID).First(&acc).Error; err != nil {
			return err
		}
		if acc.FrozenPoints < points {
			return errors.New("冻结积分不足")
		}
		acc.FrozenPoints -= points
		acc.WithdrawnPoints += points
		if err := tx.Save(&acc).Error; err != nil {
			return err
		}
		log := model.PointsLog{UserID: userID, ChangePoints: -points, Type: "withdraw_complete", Remark: remark}
		return tx.Create(&log).Error
	})
}
