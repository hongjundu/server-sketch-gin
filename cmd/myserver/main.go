package main

import (
	"github.com/hongjundu/go-level-logger"
	"github.com/hongjundu/server-sketch-gin/server"
)

func main() {
	logger.Init(0, "myserver", "/tmp", 100, 3, 30)

	s := server.NewHttpServer()
	if e := s.Run(8000); e != nil {
		logger.Fatalf("%v", e)
	}
}
