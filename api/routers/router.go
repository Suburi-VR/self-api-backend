package routers

import "github.com/gin-gonic/gin"

// Routes routes
func Routes() *gin.Engine {
	r := gin.Default()
	r.POST("/user/create", create)
	r.POST("/user/info", updateInfo)
	r.GET("/user/info", getInfo)
	r.GET("/user/contact", contact)
	r.POST("/call/start", start)
	r.POST("/call/answer", answer)
	r.POST("/call/get", get)
	r.POST("/call/status", status)
	r.POST("/call/end", end)
	r.GET("/call/history", history)

	return r
}