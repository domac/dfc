package app

import (
	"errors"
	"github.com/domac/husky"
	"log"
	"time"
)

//session server with rate limit
type SessionServer struct {
	hconfig *husky.HConfig
	tcpAddr string
}

func NewSessionServer(cfg *AppConfig) (*SessionServer, error) {

	if cfg.Tcp_address == "" {
		return nil, errors.New("start session server fail , tcp port was null")
	}

	//session config
	hc := husky.NewConfig(cfg.Max_scheduler_num,
		cfg.Read_buffer_size,
		cfg.Write_buffer_size,
		cfg.Write_channel_size,
		cfg.Read_channel_size,
		time.Duration(cfg.Idle_time)*time.Second,
		cfg.Max_seqId,
		cfg.Init_reqs_per_second,
		cfg.Max_reqs_per_second)

	return &SessionServer{
		hconfig: hc,
		tcpAddr: cfg.Tcp_address,
	}, nil
}

func (self *SessionServer) Start() {
	log.Println("open session management")
	simpleServer := husky.NewServer(self.tcpAddr, self.hconfig,
		//消息接收回调函数
		func(remoteClient *husky.HClient, p *husky.Packet) {
			println("receive a message")
		})
	simpleServer.ListenAndServer()
}
