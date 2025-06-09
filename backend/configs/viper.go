package configs

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

func findCOnfigFile() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	for {
		configPath := filepath.Join(cwd, ".env")
		if _, err := os.Stat(configPath); err == nil {
			return configPath, nil
		}

		parentDir := filepath.Dir(cwd)
		if parentDir == cwd {

			break
		}
		cwd = parentDir
	}

	return "", fmt.Errorf("config file not found")
}

func NewViperConfig() error {
	configPath, err := findCOnfigFile()
	if err != nil {
		return err
	}

	viper.SetConfigFile(configPath)
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := viper.ReadInConfig(); err != nil {
		return err
	}
	return nil
}
