package main

import (
        "fmt"
        "github.com/gin-gonic/gin"
        "net/http"
)

var available = true

func main() {
        r := gin.Default()
        r.GET("/ping", ping)
        r.GET("/", ping)
        r.GET("/switch", switchState)
        r.Run()
}

func ping(c *gin.Context) {
        if available {
                c.Writer.Write([]byte("pong"))
                c.Writer.WriteHeader(http.StatusOK)
        } else {
                c.Writer.WriteHeader(http.StatusServiceUnavailable)
        }
}

func switchState(c *gin.Context) {
        available = !available
        c.Writer.Write([]byte(fmt.Sprintf("service is %v", statusToStr())))
        c.Writer.WriteHeader(http.StatusOK)
}

func statusToStr() string {
        if available {
                return "ON"
        } else {
                return "OFF"
        }
}