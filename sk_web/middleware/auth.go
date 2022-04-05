package middleware

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"net/http"
	. "sk_web/utils"
)

func AdminAuth(ctx *gin.Context) {
	reqAuth := ctx.Request.Header.Get("Authorization")
	s := bytes.Split([]byte(reqAuth), []byte(" "))
	if len(s) != 2 {
		ctx.JSON(http.StatusOK, gin.H{
			"code":    http.StatusUnauthorized,
			"message": "请求头中无Authorization字段或错误格式",
		})
		ctx.Abort()
		return
	}
	token := string(s[1])
	ok, _ := AuthToken(token, AdminSecretKey)
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{
			"code": http.StatusUnauthorized,
			"msg":  "token error",
		})
		ctx.Abort()
		return
	}

	ctx.Next()
}

func UserAuth(ctx *gin.Context) {
	reqAuth := ctx.Request.Header.Get("Authorization")
	s := bytes.Split([]byte(reqAuth), []byte(" "))
	if len(s) != 2 {
		ctx.JSON(http.StatusOK, gin.H{
			"code":    http.StatusUnauthorized,
			"message": "请求头中无Authorization字段或错误格式",
		})
		ctx.Abort()
		return
	}
	token := string(s[1])

	ok, email := AuthToken(token, UserSecretKey)
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{
			"code": http.StatusUnauthorized,
			"msg":  "token error",
		})
		ctx.Abort()
		return
	}
	ctx.Set("userEmail", email)
	ctx.Next()
}
