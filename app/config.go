package app

//应用配置信息
type AppConfig struct {
	http_address     string
	local_store_path string
	filter_regx      []string
	peers            []map[string]string
}

func LoadConfig(configFile string) *AppConfig {

	return nil
}
