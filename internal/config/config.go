package config

import (
	"gopkg.in/yaml.v2"
	"os"
)

// Config is the struct describes settings from the .yml file
type Config struct {
	Server struct {
		Host string `yaml:"host"`
		Port string `yaml:"port"`
	}

	Database struct {
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		Dbname   string `yaml:"dbname"`
	}
}

// NewConfig returns sets settings to new configs and returns it
func NewConfig(configPath string) (*Config, error) {
	config := &Config{}

	// Open configs file
	file, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Init new YAML decode
	d := yaml.NewDecoder(file)

	// Start YAML decoding from file
	if err = d.Decode(&config); err != nil {
		return nil, err
	}

	return config, nil
}

//func GetConfigs() *Config {
//	cfg, err := NewConfig("../../configs/config2.yml")
//	if err != nil {
//		log.Fatal(err)
//	}
//	return cfg
//}

//func ValidateConfigPath(path string) error {
//	s, err := os.Stat(path)
//	if err != nil {
//		return err
//	}
//	if s.IsDir() {
//		return fmt.Errorf("'%s' is a directory, not a normal file", path)
//	}
//	return nil
//}
//
//func ParseFlags() (string, error) {
//	var configPath string
//
//	// Set up a CLI flag called "-configs" to allow users
//	// to supply the configuration file
//	flag.StringVar(&configPath, "configs", "./config2.yml", "path to configs file")
//
//	// Actually parse the flags
//	flag.Parse()
//
//	if err := ValidateConfigPath(configPath); err != nil {
//		return "", err
//	}
//
//	return configPath, nil
//}
