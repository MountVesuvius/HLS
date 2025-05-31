package main

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	// "server/db"
)

// find an easy way to have this just be loaded in
// possibly env??
const FILE_SERVICE = "http://file-service:8080/"

func Route(router *gin.Engine) {
    routes := router.Group("/api/v1/")

    routes.GET("", func(context *gin.Context) {
        context.String(http.StatusOK, "This works")
    })
}

func FileTest(router *gin.Engine) {
    routes := router.Group("/file/")

    routes.GET(":id", func(context *gin.Context) {
        id := context.Param("id")
        path := fmt.Sprintf("%s/file/%s", FILE_SERVICE, id)

        log.Println("Getting file: ", id)
        resp, err := http.Get(path)
        if err != nil {
            log.Println("Error contacting file service", err)
            context.String(http.StatusInternalServerError, "Failed to get file")
        }

        defer resp.Body.Close()
        body, _ := io.ReadAll(resp.Body)
        context.String(http.StatusOK, string(body))
    })
}

func main() {
    // db := db.Connect()

    router := gin.Default()
    Route(router)
    FileTest(router)

    err := router.Run()
    if err != nil {
        log.Fatal("Startup Failed", err)
    }
}

