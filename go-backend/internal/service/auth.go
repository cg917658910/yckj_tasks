package service

import (
	"errors"
	"task-system-go/internal/db"
	"task-system-go/internal/model"
	"task-system-go/internal/util"

	"golang.org/x/crypto/bcrypt"
)

func AdminLogin(jwt util.JWTManager, username, password string) (map[string]interface{}, error) {
	var admin model.AdminUser
	if err := db.DB.Where("username = ?", username).First(&admin).Error; err != nil {
		return nil, errors.New("用户名或密码错误")
	}
	if admin.Status != 1 {
		return nil, errors.New("账号已禁用")
	}
	if bcrypt.CompareHashAndPassword([]byte(admin.PasswordHash), []byte(password)) != nil {
		return nil, errors.New("用户名或密码错误")
	}

	token, err := jwt.Generate(map[string]interface{}{
		"role": "admin",
		"admin_id": admin.ID,
		"username": admin.Username,
	})
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"token": token,
		"admin": map[string]interface{}{
			"id": admin.ID,
			"username": admin.Username,
		},
	}, nil
}

func UserRegister(jwt util.JWTManager, username, password string) (map[string]interface{}, error) {
	var exists model.User
	if err := db.DB.Where("username = ?", username).First(&exists).Error; err == nil {
		return nil, errors.New("用户名已存在")
	}

	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	user := model.User{Username: username, PasswordHash: string(hash), Status: 1}
	if err := db.DB.Create(&user).Error; err != nil {
		return nil, errors.New("注册失败")
	}

	profile := model.UserProfile{UserID: user.ID}
	_ = db.DB.Create(&profile).Error

	rule := model.PointsRule{}
	_ = db.DB.Order("id desc").First(&rule).Error
	bonus := rule.RegisterBonusPoints
	account := model.PointsAccount{UserID: user.ID, AvailablePoints: bonus, FrozenPoints: 0, WithdrawnPoints: 0}
	_ = db.DB.Create(&account).Error

	if bonus > 0 {
		_ = db.DB.Create(&model.PointsLog{
			UserID:       user.ID,
			ChangePoints: bonus,
			Type:         "register_bonus",
			Remark:       "注册赠送积分",
		}).Error
	}

	token, err := jwt.Generate(map[string]interface{}{
		"role": "user",
		"user_id": user.ID,
		"username": user.Username,
	})
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"token": token,
		"user": map[string]interface{}{
			"id": user.ID,
			"username": user.Username,
		},
	}, nil
}

func UserLogin(jwt util.JWTManager, username, password string) (map[string]interface{}, error) {
	var user model.User
	if err := db.DB.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, errors.New("用户名或密码错误")
	}
	if user.Status != 1 {
		return nil, errors.New("账号已禁用")
	}
	if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)) != nil {
		return nil, errors.New("用户名或密码错误")
	}

	token, err := jwt.Generate(map[string]interface{}{
		"role": "user",
		"user_id": user.ID,
		"username": user.Username,
	})
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"token": token,
		"user": map[string]interface{}{
			"id": user.ID,
			"username": user.Username,
		},
	}, nil
}
