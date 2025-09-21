package router

import (
    "github.com/AhmadNizar/cata-dtc/cmd/api/http/handler"
    "github.com/gin-gonic/gin"
)

func NewRouter(apiHandler *handler.ApiHandler) *gin.Engine {
    router := gin.Default()

    v1 := router.Group("/api/v1")

    v1.POST("/sync", apiHandler.Sync)
    v1.GET("/items", apiHandler.GetItems)

    v1.GET("/health", func(c *gin.Context) {
        c.JSON(200, gin.H{"status": "ok"})
    })

    return router
}