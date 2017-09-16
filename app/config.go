package app

import (
	"errors"
	"github.com/BurntSushi/toml"
	"github.com/domac/dfc/util"
)

//应用配置信息
type AppConfig struct {
	Name                 string
	Http_address         string
	Local_store_path     string
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
	Peer                 []map[string]string
}

//加载配置
func LoadConfig(filepath string) (*AppConfig, error) {
	if filepath == "" {
		return nil, errors.New("配置文件路径为空")
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
