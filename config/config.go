package config

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type Configuration struct {
	StoragePath        string `required:"true" envconfig:"STORAGE_PATH"`
	BotAPIToken        string `required:"true" envconfig:"BOT_API_TOKEN"`
	ConfigFilename     string `default:"config.json"`
	HabitsHost         string `default:"" envconfig:"HABITS_HOST"`
	HabitsCertsPath    string `default:"/tmp" envconfig:"HABITS_CERTS_PATH"`
	GUIUserStoragePath string `default:"." envconfig:"GUI_USER_STORAGE_PATH"`
}

var Config Configuration

func LoadConfig() error {
	if err := envconfig.Process("", &Config); err != nil {
		return fmt.Errorf("can't load config: %w", err)
	}

	return nil
}
