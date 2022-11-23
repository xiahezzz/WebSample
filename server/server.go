package server

import (
	"log"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

const IP_DURATION_WINDOW = 10       //检测窗口
const IP_BEAR_TIME = 5 * 1e8        //检测窗口内，平均请求的间隔时间不能小于0.2s
const IP_RECOVER_TIME = time.Minute //限流后的恢复时间

type MyServer struct {
	router  *gin.Engine
	Request int

	ipLimiter *IpLimiter //限流器

	blackList *BlackList
}

type IpLimiter struct {
	//针对ip访问时间的限流
	ipVolume    map[string]*[IP_DURATION_WINDOW]time.Duration //同一ip连续四次访问的时间间隔
	ipLastVisit map[string]time.Time                          //ip访问次数
	ipVisitNum  map[string]int                                //ip上次访问的时间
}

type BlackList struct {
	//已禁IP、禁止原因、释放时间
	freeIpTime map[string]time.Time //ip什么时间可解禁
	reasonIpIn map[string]string    //ip被禁原因
}

func NewBlackList() *BlackList {
	return &BlackList{
		freeIpTime: make(map[string]time.Time),
		reasonIpIn: make(map[string]string),
	}
}

func NewIpLimiter() *IpLimiter {
	return &IpLimiter{
		ipVolume:    make(map[string]*[IP_DURATION_WINDOW]time.Duration),
		ipLastVisit: make(map[string]time.Time),
		ipVisitNum:  make(map[string]int),
	}
}

func NewServer() *MyServer {
	server := &MyServer{
		Request:   0,
		ipLimiter: NewIpLimiter(),
		blackList: NewBlackList(),
	}

	server.SetupRouter()

	return server
}

func (server *MyServer) SetupRouter() {
	router := gin.Default()
	v, _ := syscall.Getenv("PROTEST")
	if v == "p" {
		log.Println("限流保护已开启")
		router.Use(server.checkVolume) //限流器保护
	}

	router.GET("/ping", server.PingRes)
	router.GET("/tcp", server.TcpAliveNum)

	server.router = router
}

func (server *MyServer) Start() {
	server.router.Run()
}
