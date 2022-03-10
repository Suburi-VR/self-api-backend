package main

import (

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.POST("/user/create", create)
	r.GET("/user/info", getInfo)
	r.POST("/user/info", updateInfo)
	r.POST("/call/start", start)
	r.POST("/call/answer", answer)
	r.POST("/call/get", get)
	r.Run()
}