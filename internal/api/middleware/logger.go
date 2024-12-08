package middleware

import (
    "github.com/gin-gonic/gin"
    "log"
    "time"
)

func Logger(logger *log.Logger) gin.HandlerFunc {
    return func(c *gin.Context) {
        start := time.Now()

        // Process request
        c.Next()

        // Log details after request is processed
        duration := time.Since(start)
        logger.Printf(
            "Method: %s | Path: %s | Status: %d | Duration: %v",
            c.Request.Method,
            c.Request.URL.Path,
            c.Writer.Status(),
            duration,
        )
    }
}