package configuration

import (
	"fan2go-tui/internal/logging"
	"os"
	path2 "path"
	"time"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

type Configuration struct {
	Api       ApiConfig       `json:"api"`
	Profiling ProfilingConfig `json:"profiling"`
	Ui        UiConfig        `json:"ui"`
}

type UiConfig struct {
	UpdateInterval time.Duration `json:"updateInterval"`
}

var (
	defaultUpdateInterval = 500 * time.Millisecond
)

var CurrentConfig Configuration

// InitConfig reads in config file and ENV variables if set.
func InitConfig(cfgFile string) {
	viper.SetConfigName("fan2go-tui")

	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			logging.Error("Path Error: Couldn't detect home directory: %v", err)
			os.Exit(1)
		}

		viper.AddConfigPath(".")
		viper.AddConfigPath(home)
		viper.AddConfigPath(path2.Join(home, ".config"))
		viper.AddConfigPath(path2.Join(home, ".config", "fan2go-tui"))
		viper.AddConfigPath("/etc/fan2go-tui/")
	}

	viper.AutomaticEnv() // read in environment variables that match

	setDefaultValues()
}

func setDefaultValues() {
	viper.SetDefault("Ui", UiConfig{
		UpdateInterval: defaultUpdateInterval,
	})
	viper.SetDefault("Ui.UpdateInterval", defaultUpdateInterval)

	viper.SetDefault("Api", ApiConfig{
		Host: "127.0.0.1",
		Port: 9001,
	})
	viper.SetDefault("Api.Host", "127.0.0.1")
	viper.SetDefault("Api.Port", 9001)

	viper.SetDefault("Profiling", ProfilingConfig{
		Enabled: false,
		Host:    "localhost",
		Port:    6060,
	})
	viper.SetDefault("Profiling.Host", "127.0.0.1")
	viper.SetDefault("Profiling.Port", 6060)
}

// DetectAndReadConfigFile detects the path of the first existing config file
func DetectAndReadConfigFile() string {
	// TODO: no config for now
	_ = readInConfig()
	return GetFilePath()
}

// readInConfig reads and parses the config file
func readInConfig() error {
	return viper.ReadInConfig()
}

// GetFilePath this is only populated _after_ readInConfig()
func GetFilePath() string {
	return viper.ConfigFileUsed()
}

func LoadConfig() {
	// load default configuration values
	err := viper.Unmarshal(&CurrentConfig)
	if err != nil {
		logging.Fatal("unable to decode into struct, %v", err)
	}
}
