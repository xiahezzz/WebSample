package server

import (
	"fmt"
	"net/http"
	"regexp"
	"web/tools"

	"github.com/gin-gonic/gin"
)

func (server *MyServer) PingRes(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Request Success",
		"ip":      c.ClientIP(),
		"ip_limiter": gin.H{
			"ip_volume":     server.ipLimiter.ipVolume,
			"ip_last_visit": server.ipLimiter.ipLastVisit,
			"ip_visit_num":  server.ipLimiter.ipVisitNum,
		},
	})
}

func (server *MyServer) TcpAliveNum(c *gin.Context) {
	//msg := tools.RunCommand("netstat -nat|grep -i \"8080\"|wc -l")
	msg := tools.RunCommand("top -l 3 -o cpu")
	pattern := "sys, (.+)% idle"
	reg := regexp.MustCompile(pattern)
	result := reg.Find([]byte(msg))
	fmt.Println(result)

	c.JSON(http.StatusOK, gin.H{
		"tcp_alive_num": string(result),
	})
}
