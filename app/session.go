package app

import (
	"errors"
	"fmt"
	"github.com/domac/husky"
	"github.com/domac/husky/pb"
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

	log.Printf("session tcp server : %s\n", self.tcpAddr)

	simpleServer := husky.NewServer(self.tcpAddr, self.hconfig,
		//消息接收回调函数
		func(remoteClient *husky.HClient, p *husky.Packet) {
			if p.Data == nil || p.Header.ContentType != husky.PB_BYTES_MESSAGE {
				return
			}
			bm := &pb.BytesMessage{}
			husky.UnmarshalPbMessage(p.Data, bm)
			//接收响应
			if bm.GetHeader().GetFunctionType() == "cache_req" {

				key := bm.GetBody()
				//从in-memory缓存中尝试获取
				val, err := DefaultCacheServer.Cache().Get(key)

				if err != nil {
					resp := husky.NewPbBytesPacket(1, "cache_resp", []byte{})
					remoteClient.Write(*resp)
				} else {
					log.Println("respone to client start")
					resp := husky.NewPbBytesPacket(1, "cache_resp", val)
					remoteClient.Write(*resp)
					log.Println("respone to client end")
				}

			}
		})
	simpleServer.ListenAndServer()
}

//集群会话信息
type SessionPeers struct {
	ParentWrr  RR //轮询策略
	SiblingWrr RR //轮询策略
}

func NewSessionPeers(peerInfos []*PeerInfo) (*SessionPeers, error) {

	if peerInfos == nil || len(peerInfos) == 0 {
		return nil, errors.New("create session peers fail! peer info was null")
	}

	parentWrr := NewWeightedRR(RR_NGINX)
	siblingWrr := NewWeightedRR(RR_NGINX)

	for _, p := range peerInfos {
		if p.Addr == "" || p.Peer_type == "" {
			continue
		}

		switch p.Peer_type {
		case "parent":
			parentWrr.Add(p, p.Weight)
		case "sibling":
			siblingWrr.Add(p, p.Weight)
		}
	}

	return &SessionPeers{
		ParentWrr:  parentWrr,
		SiblingWrr: siblingWrr,
	}, nil
}

//创建集群节点连接session
func CreatePeerSession(p *PeerInfo) (*husky.HClient, error) {

	if p.Peer_type != "parent" && p.Peer_type != "sibling" {
		return nil, errors.New("create Peers fail, wrong peer type")
	}

	if p.Tcp_port == "" {
		return nil, errors.New("create Peers fail, wrong peer type")
	}

	tcp_addr := fmt.Sprintf("%s:%s", p.Addr, p.Tcp_port)
	log.Printf("create a %s session to %s\n", p.Peer_type, tcp_addr)
	conn, err := husky.Dial(tcp_addr)
	if err != nil {
		log.Println(">>>>ERROR:", err.Error())
		return nil, err
	}
	simpleClient := husky.NewClient(conn, nil, nil)
	if simpleClient != nil {
		simpleClient.Start()
		return simpleClient, nil
	}
	return nil, errors.New("create Peers fail")
}

//重置
func (self *SessionPeers) reset() {
	self.ParentWrr.Reset()
	self.SiblingWrr.Reset()
}

//清楚会话所有信息
func (self *SessionPeers) Remove() {
	self.reset()

	self.ParentWrr.RemoveAll()
	self.SiblingWrr.RemoveAll()
}
