package main

import (
	"math/rand"
	"strconv"

	"github.com/gin-gonic/gin"
	"o.o/backend/pkg/common/code/gencode"
)

const (
	ServiceCode = "ABCDEF"
	ServiceFee  = 25000
)

func main() {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.POST("/create_fulfillment", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"success": true,
			"message": "ok",
			"data": gin.H{
				"fulfillment_id": strconv.Itoa(rand.Int()),
				"shipping_code":  gencode.GenerateCode(gencode.Alphabet32, 8),
				"sort_code":      "ABC - 123",
				"tracking_url":   "http://localhost",
				"shipping_state": "created",
				"shipping_fee_lines": []gin.H{
					{
						"cost":              25000,
						"shipping_fee_type": "main",
					},
				},
			},
		})
	})

	r.POST("/get_shipping_services", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"success": true,
			"message": "ok",
			"data": []gin.H{{
				"service_code": ServiceCode,
				"name":         "Chuẩn",
				"service_fee":  ServiceFee,
			},
			},
		})
		return
	})
	r.Run(":3000")
}
