package model

import "time"

type AdminUser struct {
	ID           uint64    `gorm:"column:id;primaryKey"`
	Username     string    `gorm:"column:username"`
	PasswordHash string    `gorm:"column:password_hash"`
	Status       int       `gorm:"column:status"`
	CreatedAt    time.Time `gorm:"column:created_at"`
	UpdatedAt    time.Time `gorm:"column:updated_at"`
}

func (AdminUser) TableName() string { return "admin_users" }

type User struct {
	ID           uint64    `gorm:"column:id;primaryKey"`
	Username     string    `gorm:"column:username"`
	PasswordHash string    `gorm:"column:password_hash"`
	Status       int       `gorm:"column:status"`
	CreatedAt    time.Time `gorm:"column:created_at"`
	UpdatedAt    time.Time `gorm:"column:updated_at"`
}

func (User) TableName() string { return "users" }

type UserProfile struct {
	ID            uint64    `gorm:"column:id;primaryKey"`
	UserID        uint64    `gorm:"column:user_id"`
	WechatQRURL   string    `gorm:"column:wechat_qr_url"`
	TotalWithdraw float64   `gorm:"column:total_withdrawn"`
	CreatedAt     time.Time `gorm:"column:created_at"`
	UpdatedAt     time.Time `gorm:"column:updated_at"`
}

func (UserProfile) TableName() string { return "user_profiles" }

type PointsAccount struct {
	ID              uint64    `gorm:"column:id;primaryKey"`
	UserID          uint64    `gorm:"column:user_id"`
	AvailablePoints int       `gorm:"column:available_points"`
	FrozenPoints    int       `gorm:"column:frozen_points"`
	WithdrawnPoints int       `gorm:"column:withdrawn_points"`
	CreatedAt       time.Time `gorm:"column:created_at"`
	UpdatedAt       time.Time `gorm:"column:updated_at"`
}

func (PointsAccount) TableName() string { return "points_accounts" }

type PointsRule struct {
	ID                 uint64    `gorm:"column:id;primaryKey"`
	ExchangeRate       int       `gorm:"column:exchange_rate"`
	MinWithdrawAmount  float64   `gorm:"column:min_withdraw_amount"`
	RegisterBonusPoints int      `gorm:"column:register_bonus_points"`
	UpdatedAt          time.Time `gorm:"column:updated_at"`
}

func (PointsRule) TableName() string { return "points_rules" }

type PointsLog struct {
	ID           uint64    `gorm:"column:id;primaryKey"`
	UserID       uint64    `gorm:"column:user_id"`
	ChangePoints int       `gorm:"column:change_points"`
	Type         string    `gorm:"column:type"`
	RefID        *uint64   `gorm:"column:ref_id"`
	Remark       string    `gorm:"column:remark"`
	CreatedAt    time.Time `gorm:"column:created_at"`
}

func (PointsLog) TableName() string { return "points_logs" }

type Task struct {
	ID          uint64    `gorm:"column:id;primaryKey"`
	Title       string    `gorm:"column:title"`
	Summary     string    `gorm:"column:summary"`
	Detail      string    `gorm:"column:detail"`
	DocURL      string    `gorm:"column:doc_url"`
	RewardPoints int      `gorm:"column:reward_points"`
	Status      int       `gorm:"column:status"`
	CreatedAt   time.Time `gorm:"column:created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at"`
}

func (Task) TableName() string { return "tasks" }

type TaskClaim struct {
	ID              uint64     `gorm:"column:id;primaryKey"`
	TaskID          uint64     `gorm:"column:task_id"`
	UserID          uint64     `gorm:"column:user_id"`
	Status          int        `gorm:"column:status"`
	ClaimedAt       time.Time  `gorm:"column:claimed_at"`
	SubmittedAt     *time.Time `gorm:"column:submitted_at"`
	ReviewedAt      *time.Time `gorm:"column:reviewed_at"`
	ReviewResult    *int       `gorm:"column:review_result"`
	RejectReason    string     `gorm:"column:reject_reason"`
	RewardPointsFinal *int     `gorm:"column:reward_points_final"`
}

func (TaskClaim) TableName() string { return "task_claims" }

type TaskSubmission struct {
	ID        uint64    `gorm:"column:id;primaryKey"`
	ClaimID   uint64    `gorm:"column:claim_id"`
	Remark    string    `gorm:"column:remark"`
	CreatedAt time.Time `gorm:"column:created_at"`
}

func (TaskSubmission) TableName() string { return "task_submissions" }

type TaskSubmissionImage struct {
	ID           uint64    `gorm:"column:id;primaryKey"`
	SubmissionID uint64    `gorm:"column:submission_id"`
	ImageURL     string    `gorm:"column:image_url"`
	CreatedAt    time.Time `gorm:"column:created_at"`
}

func (TaskSubmissionImage) TableName() string { return "task_submission_images" }

type Withdrawal struct {
	ID          uint64     `gorm:"column:id;primaryKey"`
	UserID      uint64     `gorm:"column:user_id"`
	Amount      float64    `gorm:"column:amount"`
	PointsCost  int        `gorm:"column:points_cost"`
	QRURL       string     `gorm:"column:qr_url"`
	Status      int        `gorm:"column:status"`
	RejectReason string    `gorm:"column:reject_reason"`
	CreatedAt   time.Time  `gorm:"column:created_at"`
	ReviewedAt  *time.Time `gorm:"column:reviewed_at"`
	AdminID     *uint64    `gorm:"column:admin_id"`
}

func (Withdrawal) TableName() string { return "withdrawals" }
