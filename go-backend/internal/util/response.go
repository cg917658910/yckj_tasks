package util

import "github.com/gin-gonic/gin"

type Resp struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func JSONSuccess(c *gin.Context, data interface{}) {
	c.JSON(200, Resp{Code: 0, Message: "ok", Data: data})
}

func JSONError(c *gin.Context, message string, code int) {
	if code == 0 {
		code = 1
	}
	c.JSON(200, Resp{Code: code, Message: message, Data: []interface{}{}})
}
