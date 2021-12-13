package controller

import (
	"fmt"
	"fs/config"
	"fs/middleware"
	"fs/router"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"sync"
	"time"
)

type SamllFileServer struct {
	config *config.Config
	router *gin.Engine
	c      chan map[string]string
	mutex  sync.Mutex
}

func NewSamllFileServer(conf *config.Config) *SamllFileServer {
	r := router.NewRouter()
	server := &SamllFileServer{
		config: conf,
		router: r,
		c:      make(chan map[string]string, 1024),
	}
	return server
}

func (server *SamllFileServer) Load() {
	if len(server.config.Routers) < 1 {
		log.Panic("请设置至少一个路由")
	}
	for _, v := range server.config.Routers {
		server.router.StaticFS(v.ServerPath, http.Dir(v.LocalPath))
	}
}

func (server *SamllFileServer) SetMiddleWare() {
	server.router.Use(gin.Recovery())
	server.router.Use(middleware.LogPrintter(server.c))
}

func (server *SamllFileServer) ProcessAccessLog() {
	logs := make([]map[string]string, 0)
	ticker := time.NewTicker(time.Second * 5)
	for {
		select {
		case cc := <-server.c:
			logs = append(logs, cc)
		case <-ticker.C:
			tmp := make([]map[string]string, len(logs))
			copy(tmp, logs)
			logs = logs[0:0]
			server.Store(tmp)
		}
	}
}

func (server *SamllFileServer) Store(logs []map[string]string) {
	for _, v := range logs {
		fmt.Println(v)
	}
}

func (server *SamllFileServer) Start() {
	go server.ProcessAccessLog()
	server.router.Run(server.config.Address)
}
