package seckill

import (
	"github.com/garyburd/redigo/redis"
	"github.com/gin-gonic/gin"
	"net/http"
	"sk_web/handler/seckill/predis"
	. "sk_web/utils"
)

func FrontSeckillResultHandler(ctx *gin.Context) {
	redisConn := predis.GetRedisConn()
	defer redisConn.Close()
	email, exist := ctx.Get("userEmail")
	if !exist {
		ctx.JSON(200, gin.H{
			"code": http.StatusUnauthorized,
			"msg":  "请先登录",
		})
		return
	}
	m, err := redis.String(redisConn.Do("GET", email))
	if err != nil {
		ctx.JSON(200, gin.H{
			"code": http.StatusOK, //前端设计，必须是200才能显示信息
			"msg":  "系统错误",
		})
		return
	}
	ret := StrToMap(m)
	ret["code"] = http.StatusOK //前端设计，必须是200才能显示信息
	ctx.JSON(http.StatusOK, ret)
	redisConn.Do("DEL", email)
}
