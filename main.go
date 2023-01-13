package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	defer func() {
		if e := recover(); e != nil {
			fmt.Println("---------", e)
		}
	}()
	data := make([]string, 0)
	data = append(data, "hello123")
	r := gin.Default()
	r.GET("/hello", func(c *gin.Context) {
		resp := map[string]interface{}{
			"code": 200,
			"data": data,
		}
		c.JSON(http.StatusOK, gin.H{"data": resp})
	})

	r.Run(":8081")
}
