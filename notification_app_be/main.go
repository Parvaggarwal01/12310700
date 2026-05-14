package main

import (
	"net/http"
	"notification_be/utils"

	"github.com/gin-gonic/gin"
)

const ServerAccessToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJNYXBDbGFpbXMiOnsiYXVkIjoiaHR0cDovLzIwLjI0NC41Ni4xNDQvZXZhbHVhdGlvbi1zZXJ2aWNlIiwiZW1haWwiOiJwYXJ2YWdnYXJ3YWwxMzBAZ21haWwuY29tIiwiZXhwIjoxNzc4NzYzOTM0LCJpYXQiOjE3Nzg3NjMwMzQsImlzcyI6IkFmZm9yZCBNZWRpY2FsIFRlY2hub2xvZ2llcyBQcml2YXRlIExpbWl0ZWQiLCJqdGkiOiI0OTIzMDRiNi0yNjYyLTRiOWUtOWI4NC04MzNmZTZiYTI0NDEiLCJsb2NhbGUiOiJlbi1JTiIsIm5hbWUiOiJwYXJ2IGFnZ2Fyd2FsIiwic3ViIjoiNzViOGI3NjQtMmM2MS00Mzg4LWFhOTYtZDNjMTI2ZWMxODA1In0sImVtYWlsIjoicGFydmFnZ2Fyd2FsMTMwQGdtYWlsLmNvbSIsIm5hbWUiOiJwYXJ2IGFnZ2Fyd2FsIiwicm9sbE5vIjoiMTIzMTA3MDAiLCJhY2Nlc3NDb2RlIjoiVFJ2WldxIiwiY2xpZW50SUQiOiI3NWI4Yjc2NC0yYzYxLTQzODgtYWE5Ni1kM2MxMjZlYzE4MDUiLCJjbGllbnRTZWNyZXQiOiJId3hSQVRlaFlydmpiVUROIn0.V51US7AXMsdzkafe3WREw--9vO_HebM7wJQsehX3jiU"

func main() {
	r := gin.Default()

	// CORS Middleware for the React Frontend
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	api := r.Group("/api/v1")
	{
		// Priority Inbox Endpoint
		api.GET("/priority-inbox", func(c *gin.Context) {
			utils.Log("backend", "info", "controller", "Incoming request for priority inbox")

			// 1. Fetch from Test Server
			notifications, err := FetchNotifications(ServerAccessToken)
			if err != nil {
				utils.Log("backend", "error", "api", "Failed to fetch from test server API")
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch notifications"})
				return
			}
			utils.Log("backend", "info", "domain", "Successfully fetched notifications from external API")

			// 2. Process Top 10 Algorithm
			top10 := GetTopPriority(notifications, 10)
			utils.Log("backend", "info", "domain", "Successfully calculated top 10 priority notifications")

			// 3. Return to Frontend
			c.JSON(http.StatusOK, gin.H{
				"count": len(top10),
				"data":  top10,
			})
		})
	}

	utils.Log("backend", "info", "config", "Priority Backend Server started successfully on :8080")
	r.Run(":8080")
}
