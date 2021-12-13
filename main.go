package main

import (
	"fs/config"
	"fs/controller"
	"log"
	"os"
	"os/signal"
)

func waitForSignal() os.Signal {
	signalChan := make(chan os.Signal, 1)
	defer close(signalChan)
	signal.Notify(signalChan, os.Kill, os.Interrupt)
	s := <-signalChan
	signal.Stop(signalChan)
	return s
}

func main() {
	conf, err := config.LoadConfig("./config/config.json")
	if err != nil {
		log.Panic(err)
	}

	server := controller.NewSamllFileServer(conf)
	server.SetMiddleWare()
	server.Load()
	server.Start()

	signalNumber := waitForSignal()
	log.Println("signal received, broker closed.", signalNumber)
}
