package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.GET("/pingTime", func(c *gin.Context) {
		// JSON serializer is available on gin context
		c.JSON(http.StatusOK, gin.H{
			"serverTime": time.Now().UTC(),
		})
	})
	r.Run(":8000") // Listen and serve on 0.0.0.0:8000
}
