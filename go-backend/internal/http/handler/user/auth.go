package user

import (
	"net/http"
	"task-system-go/internal/service"
	"task-system-go/internal/util"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	JWT util.JWTManager
}

func (h AuthHandler) Register(c *gin.Context) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || req.Username == "" || req.Password == "" {
		util.JSONError(c, "参数不完整", 1)
		return
	}

	data, err := service.UserRegister(h.JWT, req.Username, req.Password)
	if err != nil {
		util.JSONError(c, err.Error(), 1)
		return
	}

	c.JSON(http.StatusOK, util.Resp{Code: 0, Message: "ok", Data: data})
}

func (h AuthHandler) Login(c *gin.Context) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || req.Username == "" || req.Password == "" {
		util.JSONError(c, "参数不完整", 1)
		return
	}

	data, err := service.UserLogin(h.JWT, req.Username, req.Password)
	if err != nil {
		util.JSONError(c, err.Error(), 1)
		return
	}

	c.JSON(http.StatusOK, util.Resp{Code: 0, Message: "ok", Data: data})
}

func (h AuthHandler) Logout(c *gin.Context) {
	util.JSONSuccess(c, gin.H{})
}
