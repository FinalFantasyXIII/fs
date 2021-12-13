package controller

import (
	"fmt"
	"fs/config"
	"fs/middleware"
	"fs/model"
	"fs/router"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"net/http"
	"sync"
	"time"
)

type SamllFileServer struct {
	config *config.Config
	router *gin.Engine
	db 	   *gorm.DB
	c      chan map[string]string
	mutex  sync.Mutex
}

func NewSamllFileServer(conf *config.Config) *SamllFileServer {
	r := router.NewRouter()
	dns := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",conf.Mysql.User,conf.Mysql.Key,conf.Mysql.Address,conf.Mysql.DB)
	mysqlConn ,err := gorm.Open(sqlite.Open(dns),&gorm.Config{})
	if err != nil {
		log.Panic("DB init error",err)
	}

	server := &SamllFileServer{
		config: conf,
		router: r,
		db: mysqlConn,
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
	acLogs := make([]model.AccessLog,0)
	for _, v := range logs {
		acLogs = append(acLogs,model.AccessLog{
			Method:     v["method"],
			Root:       "",
			Path:       v["url"],
			ClientIp:   v["client_ip"],
			AccessTime: v["access_time"],
		})
	}

	ret := server.db.CreateInBatches(acLogs,len(acLogs))
	if ret.Error != nil{
		fmt.Println("部分log存储失败", ret.Error)
	}
}

func (server *SamllFileServer) Start() {
	go server.ProcessAccessLog()
	server.router.Run(server.config.Address)
}
