package assets

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strconv"
)

// ServerConfig - The server configuration used by the console
var ServerConfig = &serverConfig{}

// ServerConfig - A Server configuration file for connecting to it.
type serverConfig struct {
	User          string `json:"user"` // This value is actually ignored for the most part (cert CN is used instead)
	LHost         string `json:"lhost"`
	LPort         int    `json:"lport"`
	CACertificate string `json:"ca_certificate"`
	PrivateKey    string `json:"private_key"`
	Certificate   string `json:"certificate"`
	IsDefault     bool   `json:"is_default"`
}

// LoadServerConfig - Get Server address and certificates
func LoadServerConfig() error {

	// If skip checks if the console has builtin server config
	if HasBuiltinConfig() {
		ReadConfigVars()
		return nil
	}

	// Read config in file and fill values in builtin-config file
	ServerConfig = getDefaultServerConfig()
	// ReadConfigFile(GetConfigDir())

	return nil
}

// GetConfigs - Return all Server configs available
func GetConfigs() (configs map[string]*serverConfig) {

	configDir := GetConfigDir()
	configFiles, err := ioutil.ReadDir(configDir)
	if err != nil {
		log.Printf("No configs found %v", err)
		return map[string]*serverConfig{}
	}

	configs = map[string]*serverConfig{}
	for _, confFile := range configFiles {
		confFilePath := path.Join(configDir, confFile.Name())
		// log.Printf("Parsing config %s", confFilePath)

		conf, err := ReadConfigFile(confFilePath)
		if err != nil {
			continue
		}
		configs[conf.LHost] = conf
	}

	return
}

// ReadConfigVars - Loads all configuration from binary
func ReadConfigVars() (config *serverConfig, err error) {

	conf := ServerConfig
	conf.LHost = ServerLHost
	conf.LPort, _ = strconv.Atoi(ServerLPort)
	conf.User = ServerUser
	conf.CACertificate = ServerCACertificate
	conf.Certificate = ServerCertificate
	conf.PrivateKey = ServerPrivateKey

	return
}

// ReadConfigFile - Loads all configuration from file
func ReadConfigFile(confFilePath string) (config *serverConfig, err error) {

	// Read file
	confFile, err := os.Open(confFilePath)
	defer confFile.Close()
	if err != nil {
		log.Printf("Open failed %v", err)
		return nil, err
	}
	data, err := ioutil.ReadAll(confFile)
	if err != nil {
		log.Printf("Read failed %v", err)
		return nil, err
	}
	conf := &serverConfig{}
	err = json.Unmarshal(data, conf)
	if err != nil {
		log.Printf("Parse failed %v", err)
		return nil, err
	}

	// Set builtin values
	ServerLHost = conf.LHost
	ServerLPort = strconv.Itoa(int(conf.LPort))
	ServerUser = conf.User
	ServerCACertificate = conf.CACertificate
	ServerCertificate = conf.Certificate
	ServerPrivateKey = conf.PrivateKey

	return conf, nil
}

// SaveConfig - Save a configuration to disk
func SaveConfig(config *serverConfig) (err error) {
	if config.LHost == "" || config.User == "" {
		return errors.New("Empty config")
	}
	configDir := GetConfigDir()
	filename := fmt.Sprintf("%s_%s.cfg", filepath.Base(config.User), filepath.Base(config.LHost))
	saveTo, _ := filepath.Abs(path.Join(configDir, filename))
	configJSON, _ := json.Marshal(config)
	err = ioutil.WriteFile(saveTo, configJSON, 0600)
	if err != nil {
		log.Printf("Failed to write config to: %s (%v)", saveTo, err)
		return err
	}
	log.Printf("Saved new client config to: %s", saveTo)
	return nil
}

// HasBuiltinConfig - Check if console was compiled with default server configuration
func HasBuiltinConfig() bool {
	if HasBuiltinServer != "" {
		fmt.Println("Server config: Found compile-time server values")
		return true
	}
	return false
}

func getDefaultServerConfig() *serverConfig {
	configs := GetConfigs()
	if len(configs) == 0 {
		// fmt.Printf(Warnf+"No config files found at %s or -config\n", assets.GetConfigDir())
		return nil
	}

	var config *serverConfig
	for _, conf := range configs {
		if conf.IsDefault {
			config = conf
		}
	}

	return config
}
