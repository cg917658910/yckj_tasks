package model

import "time"

type AdminUser struct {
	ID           uint64    `gorm:"column:id;primaryKey" json:"id"`
	Username     string    `gorm:"column:username" json:"username"`
	PasswordHash string    `gorm:"column:password_hash" json:"password_hash"`
	Status       int       `gorm:"column:status" json:"status"`
	CreatedAt    time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt    time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (AdminUser) TableName() string { return "admin_users" }

type User struct {
	ID           uint64    `gorm:"column:id;primaryKey" json:"id"`
	Username     string    `gorm:"column:username" json:"username"`
	PasswordHash string    `gorm:"column:password_hash" json:"password_hash"`
	Status       int       `gorm:"column:status" json:"status"`
	CreatedAt    time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt    time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (User) TableName() string { return "users" }

type UserProfile struct {
	ID            uint64    `gorm:"column:id;primaryKey" json:"id"`
	UserID        uint64    `gorm:"column:user_id" json:"user_id"`
	WechatQRURL   string    `gorm:"column:wechat_qr_url" json:"wechat_qr_url"`
	TotalWithdraw float64   `gorm:"column:total_withdrawn" json:"total_withdrawn"`
	CreatedAt     time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt     time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (UserProfile) TableName() string { return "user_profiles" }

type PointsAccount struct {
	ID              uint64    `gorm:"column:id;primaryKey" json:"id"`
	UserID          uint64    `gorm:"column:user_id" json:"user_id"`
	AvailablePoints int       `gorm:"column:available_points" json:"available_points"`
	FrozenPoints    int       `gorm:"column:frozen_points" json:"frozen_points"`
	WithdrawnPoints int       `gorm:"column:withdrawn_points" json:"withdrawn_points"`
	CreatedAt       time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt       time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (PointsAccount) TableName() string { return "points_accounts" }

type PointsRule struct {
	ID                  uint64    `gorm:"column:id;primaryKey" json:"id"`
	ExchangeRate        int       `gorm:"column:exchange_rate" json:"exchange_rate"`
	MinWithdrawAmount   float64   `gorm:"column:min_withdraw_amount" json:"min_withdraw_amount"`
	RegisterBonusPoints int       `gorm:"column:register_bonus_points" json:"register_bonus_points"`
	UpdatedAt           time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (PointsRule) TableName() string { return "points_rules" }

type PointsLog struct {
	ID           uint64    `gorm:"column:id;primaryKey" json:"id"`
	UserID       uint64    `gorm:"column:user_id" json:"user_id"`
	ChangePoints int       `gorm:"column:change_points" json:"change_points"`
	Type         string    `gorm:"column:type" json:"type"`
	RefID        *uint64   `gorm:"column:ref_id" json:"ref_id"`
	Remark       string    `gorm:"column:remark" json:"remark"`
	CreatedAt    time.Time `gorm:"column:created_at" json:"created_at"`
}

func (PointsLog) TableName() string { return "points_logs" }

type Task struct {
	ID           uint64    `gorm:"column:id;primaryKey" json:"id"`
	Title        string    `gorm:"column:title" json:"title"`
	Summary      string    `gorm:"column:summary" json:"summary"`
	Detail       string    `gorm:"column:detail" json:"detail"`
	DocURL       string    `gorm:"column:doc_url" json:"doc_url"`
	RewardPoints int       `gorm:"column:reward_points" json:"reward_points"`
	Status       int       `gorm:"column:status" json:"status"`
	CreatedAt    time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt    time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (Task) TableName() string { return "tasks" }

type TaskClaim struct {
	ID                uint64     `gorm:"column:id;primaryKey" json:"id"`
	TaskID            uint64     `gorm:"column:task_id" json:"task_id"`
	UserID            uint64     `gorm:"column:user_id" json:"user_id"`
	Status            int        `gorm:"column:status" json:"status"`
	ClaimedAt         time.Time  `gorm:"column:claimed_at" json:"claimed_at"`
	SubmittedAt       *time.Time `gorm:"column:submitted_at" json:"submitted_at"`
	ReviewedAt        *time.Time `gorm:"column:reviewed_at" json:"reviewed_at"`
	ReviewResult      *int       `gorm:"column:review_result" json:"review_result"`
	RejectReason      string     `gorm:"column:reject_reason" json:"reject_reason"`
	RewardPointsFinal *int       `gorm:"column:reward_points_final" json:"reward_points_final"`
}

func (TaskClaim) TableName() string { return "task_claims" }

type TaskSubmission struct {
	ID        uint64    `gorm:"column:id;primaryKey" json:"id"`
	ClaimID   uint64    `gorm:"column:claim_id" json:"claim_id"`
	Remark    string    `gorm:"column:remark" json:"remark"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
}

func (TaskSubmission) TableName() string { return "task_submissions" }

type TaskSubmissionImage struct {
	ID           uint64    `gorm:"column:id;primaryKey" json:"id"`
	SubmissionID uint64    `gorm:"column:submission_id" json:"submission_id"`
	ImageURL     string    `gorm:"column:image_url" json:"image_url"`
	CreatedAt    time.Time `gorm:"column:created_at" json:"created_at"`
}

func (TaskSubmissionImage) TableName() string { return "task_submission_images" }

type Withdrawal struct {
	ID           uint64     `gorm:"column:id;primaryKey" json:"id"`
	UserID       uint64     `gorm:"column:user_id" json:"user_id"`
	Amount       float64    `gorm:"column:amount" json:"amount"`
	PointsCost   int        `gorm:"column:points_cost" json:"points_cost"`
	QRURL        string     `gorm:"column:qr_url" json:"qr_url"`
	Status       int        `gorm:"column:status" json:"status"`
	RejectReason string     `gorm:"column:reject_reason" json:"reject_reason"`
	CreatedAt    time.Time  `gorm:"column:created_at" json:"created_at"`
	ReviewedAt   *time.Time `gorm:"column:reviewed_at" json:"reviewed_at"`
	AdminID      *uint64    `gorm:"column:admin_id" json:"admin_id"`
}

func (Withdrawal) TableName() string { return "withdrawals" }
