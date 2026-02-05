package router

import (
	"task-system-go/internal/config"
	"task-system-go/internal/http/handler/admin"
	"task-system-go/internal/http/handler/common"
	"task-system-go/internal/http/handler/user"
	"task-system-go/internal/http/middleware"
	"task-system-go/internal/util"

	"github.com/gin-gonic/gin"
)

func New(cfg config.Config) *gin.Engine {
	r := gin.Default()
	r.Use(middleware.CORS(cfg))

	jwtMgr := util.JWTManager{Secret: cfg.JWTSecret, Issuer: cfg.JWTIssuer, TTL: cfg.JWTTTL}

	adminAuth := admin.AuthHandler{JWT: jwtMgr}
	userAuth := user.AuthHandler{JWT: jwtMgr}
	upload := common.UploadHandler{Cfg: cfg}

	// static files
	r.Static("/static", "./static")

	r.POST("/admin/auth/login", adminAuth.Login)
	r.POST("/admin/auth/logout", middleware.AdminAuth(jwtMgr), adminAuth.Logout)
	// admin routes
	adminGroup := r.Group("/admin", middleware.AdminAuth(jwtMgr))
	{
		adminGroup.GET("/tasks", admin.TaskHandler{}.List)
		adminGroup.POST("/tasks", admin.TaskHandler{}.Create)
		adminGroup.PUT("/tasks/:id", admin.TaskHandler{}.Update)
		adminGroup.PUT("/tasks/:id/off", admin.TaskHandler{}.Off)

		adminGroup.GET("/claims", admin.ClaimHandler{}.List)
		adminGroup.POST("/claims/:id/approve", admin.ClaimHandler{}.Approve)
		adminGroup.POST("/claims/:id/reject", admin.ClaimHandler{}.Reject)

		adminGroup.GET("/points/rules", admin.PointsHandler{}.Rules)
		adminGroup.PUT("/points/rules", admin.PointsHandler{}.UpdateRules)
		adminGroup.GET("/points/logs", admin.PointsHandler{}.Logs)

		adminGroup.GET("/withdrawals", admin.WithdrawalHandler{}.List)
		adminGroup.POST("/withdrawals/:id/pay", admin.WithdrawalHandler{}.Pay)
		adminGroup.POST("/withdrawals/:id/reject", admin.WithdrawalHandler{}.Reject)

		adminGroup.GET("/users", admin.UserHandler{}.List)
		adminGroup.PUT("/users/:id/status", admin.UserHandler{}.Status)
		adminGroup.POST("/users/:id/points", admin.UserHandler{}.AdjustPoints)
		adminGroup.GET("/users/:id/tasks", admin.UserHandler{}.Tasks)
		adminGroup.GET("/users/:id/withdrawals", admin.UserHandler{}.Withdrawals)

		adminGroup.POST("/upload/image", upload.Image)
	}

	r.POST("/user/auth/register", userAuth.Register)
	r.POST("/user/auth/login", userAuth.Login)
	r.POST("/user/auth/logout", middleware.UserAuth(jwtMgr), userAuth.Logout)

	userGroup := r.Group("/user", middleware.UserAuth(jwtMgr))
	{
		userGroup.GET("/tasks", user.TaskHandler{}.List)
		userGroup.GET("/tasks/:id", user.TaskHandler{}.Detail)
		userGroup.POST("/tasks/:id/claim", user.TaskHandler{}.Claim)

		userGroup.GET("/claims/current", user.ClaimHandler{}.Current)
		userGroup.GET("/claims/history", user.ClaimHandler{}.History)
		userGroup.POST("/claims/:id/submit", user.ClaimHandler{}.Submit)

		userGroup.GET("/profile", user.ProfileHandler{}.Info)
		userGroup.PUT("/profile/password", user.ProfileHandler{}.ChangePassword)
		userGroup.PUT("/profile/wechat-qr", user.ProfileHandler{}.UpdateWechatQR)

		userGroup.POST("/withdrawals", user.WithdrawalHandler{}.Apply)
		userGroup.GET("/withdrawals", user.WithdrawalHandler{}.List)

		userGroup.GET("/points/logs", user.PointsHandler{}.Logs)

		userGroup.POST("/upload/image", upload.Image)
	}

	return r
}
