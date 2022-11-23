package server

import (
	"net/http"
	"strings"
	"time"
	"web/tools"

	"github.com/gin-gonic/gin"
)

func (server *MyServer) opIdVolume(ip string) {

	if server.ipLimiter.ipVisitNum[ip] == 0 {
		server.ipLimiter.ipVolume[ip] = &[IP_DURATION_WINDOW]time.Duration{time.Duration(3600000000000)}
	} else {
		server.ipLimiter.ipVolume[ip][server.ipLimiter.ipVisitNum[ip]%IP_DURATION_WINDOW] = time.Since(server.ipLimiter.ipLastVisit[ip])
	}

	server.ipLimiter.ipLastVisit[ip] = time.Now()
	server.ipLimiter.ipVisitNum[ip]++
}

func (server *MyServer) addFireWallRule(ip string) string {
	msg := tools.RunCommand(strings.Join([]string{"iptables -I INPUT -s ", ip, " -j DROP"}, ""))
	return msg
}

func (server *MyServer) checkVolume(c *gin.Context) {
	ip := c.ClientIP()

	//ip已被禁，就不必更新其访问了
	if _, ok := server.blackList.reasonIpIn[ip]; !ok {
		server.opIdVolume(ip)
		visitSpeed := time.Duration(0)
		for _, v := range server.ipLimiter.ipVolume[ip] {
			visitSpeed += v
		}

		if visitSpeed/IP_DURATION_WINDOW < IP_BEAR_TIME {
			server.blackList.reasonIpIn[ip] = "IP_LIMITER"
			server.blackList.freeIpTime[ip] = server.ipLimiter.ipLastVisit[ip].Add(IP_RECOVER_TIME)
			/*if h, ok := c.Writer.(http.Hijacker); ok {
				c, _, err := h.Hijack()
				if err != nil {
					c.Close()
				}
			}*/
			server.addFireWallRule(ip)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"ip":        c.ClientIP(),
				"msg":       "DDos Attack Detected!",
				"free_time": server.blackList.freeIpTime[ip],
				"reason":    server.blackList.reasonIpIn[ip],
			})

		}
	}
}
