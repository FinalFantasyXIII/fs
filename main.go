package main

import (
	"encoding/json"
	"fmt"
	"fs/middleware"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"net/http"
)

type Config struct {
	Address   string      `json:"address"`
	Routers   []FileTree  `json:"routers"`
}

type FileTree struct {
	ServerPath	string	`json:"server_path"`
	LocalPath	string 	`json:"local_path"`
}


func LoadConfig(filename string) (*Config, error) {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var config Config
	err = json.Unmarshal(content, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

type SamllFileServer struct {
	config		*Config
	router		*gin.Engine
	c 			chan map[string]string
}

func NewSamllFileServer (conf *Config) *SamllFileServer{
	r := gin.New()
	server := &SamllFileServer{
		config: conf,
		router: r,
		c: make(chan map[string]string,10),
	}
	return server
}

func (server *SamllFileServer) Load(){
	if len(server.config.Routers) < 1 {
		log.Panic("请设置至少一个路由")
	}
	for _ , v := range server.config.Routers{
		server.router.StaticFS(v.ServerPath,http.Dir(v.LocalPath))
	}
}

func (server *SamllFileServer) SetMiddleWare(){
	server.router.Use(gin.Recovery())
	server.router.Use(middleware.LogPrintter(server.c))
}


func (server *SamllFileServer) ReadAccessLog(){
	for t := range server.c{
		fmt.Println(t)
	}
}
func (server *SamllFileServer) Start(){
	go server.ReadAccessLog()
	server.router.Run(server.config.Address)
}




func main(){
	conf ,err := LoadConfig("config.json")
	if err != nil{
		log.Panic(err)
	}

	server := NewSamllFileServer(conf)
	server.SetMiddleWare()
	server.Load()
	server.Start()
}
