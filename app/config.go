package app

import (
	"errors"
	"github.com/BurntSushi/toml"
	"github.com/domac/dfc/util"
)

//应用配置
type AppConfig struct {
	Name                 string
	Log_path             string
	Log_level            string
	Http_address         string
	Tcp_address          string
	Local_store_path     string
	Cache_max_size       int
	Cache_ttl            int
	Max_scheduler_num    int
	Read_buffer_size     int
	Write_buffer_size    int
	Write_channel_size   int
	Read_channel_size    int
	Idle_time            int
	Max_seqId            int
	Init_reqs_per_second int
	Max_reqs_per_second  int
	Filter_regx          []string
	Peer                 []*PeerInfo
}

//peer node config
type PeerInfo struct {
	Name      string
	Addr      string
	Peer_type string
	Http_port string
	Tcp_port  string
	Weight    int
}

func (self *AppConfig) GetParentPeerNodes() ([]*PeerInfo, error) {
	return self.getPeerNodes("parent")
}

func (self *AppConfig) GetSublingPeerNodes() ([]*PeerInfo, error) {
	return self.getPeerNodes("sibling")
}

func (self *AppConfig) getPeerNodes(ptype string) ([]*PeerInfo, error) {
	peers := self.Peer
	if ptype != "" && len(peers) > 0 {
		results := []*PeerInfo{}
		for _, p := range peers {
			if p.Peer_type != ptype {
				continue
			}
			results = append(results, p)
		}
		if len(results) > 0 {
			return results, nil
		}
	}
	return nil, errors.New("no peers")
}

func LoadConfig(filepath string) (*AppConfig, error) {
	if filepath == "" {
		return nil, errors.New("the config file dir is empty")
	}
	if err := util.CheckDataFileExist(filepath); err != nil {
		return nil, err
	}
	var cfg *AppConfig
	if filepath != "" {
		_, err := toml.DecodeFile(filepath, &cfg)
		if err != nil {
			return nil, err
		}
	}
	return cfg, nil
}
