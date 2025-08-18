package middleware

import (
    "time"

    "github.com/gin-gonic/gin"
)

// Logger is a middleware that logs request details
func Logger() gin.HandlerFunc {
    return func(c *gin.Context) {
        // Start time
        startTime := time.Now()

        // Process request
        c.Next()

        // End time
        endTime := time.Now()

        // Request details
        latency := endTime.Sub(startTime)
        clientIP := c.ClientIP()
        method := c.Request.Method
        statusCode := c.Writer.Status()
        path := c.Request.URL.Path

        // Log request
        gin.DefaultWriter.Write([]byte(
            "[GIN] " + endTime.Format("2006/01/02 - 15:04:05") +
                " | " + statusCode + " | " + latency.String() +
                " | " + clientIP + " | " + method + " | " + path + "\n",
        ))
    }
}