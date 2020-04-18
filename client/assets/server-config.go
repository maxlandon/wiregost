package assets

// ServerConfig - A Server configuration file for connecting to it.
type ServerConfig struct{}

// GetConfigs - Return all Server configs available
func GetConfigs() (configs map[string]*ServerConfig) {

	return
}

// ReadConfig - Loads all configuration from file
func ReadConfig(confPath string) (config *ServerConfig, err error) {
	return
}

// SaveConfig - Save a configuration to disk
func SaveConfig(config *ServerConfig) (err error) {
	return
}
