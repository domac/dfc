package app

import (
	"github.com/domac/husky"
	"log"
	"time"
)

//会话服务
type SessionServer struct {
	rateLimitNum int
}

func NewSessionServer(rateLimitNum int) *SessionServer {
	return &SessionServer{
		rateLimitNum: rateLimitNum}
}

func (self *SessionServer) Start() {
	log.Println("startup a session serverco")
	cfg := husky.NewConfig(1000, 4*1024, 4*1024, 10000, 10000, 10*time.Second, 160000, -1, self.rateLimitNum)
	simpleServer := husky.NewServer("localhost:10201", cfg, func(remoteClient *husky.HClient, p *husky.Packet) {
		println("receive a message")
	})
	simpleServer.ListenAndServer()
}
