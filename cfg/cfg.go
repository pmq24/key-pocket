package cfg

import (
	"fmt"

	"github.com/spf13/viper"
)

const (
	EncryptedExt = ".kpenc"
	KeyFileFormat = "kpkey.%s"
	CfgFileFormat = "kpcfg.%s"
)

func ReadConfig() (*Config, error) {
	config := viper.New()

	config.BindPFlag("dir", dir)
	config.BindPFlag("profile", profile)

	config.SetConfigName(fmt.Sprintf(CfgFileFormat, profile.Value.String()))
	config.SetConfigType("yml")
	config.AddConfigPath(dir.Value.String())

	err := config.ReadInConfig()
	if err != nil {
		return nil, ErrReadConfig{Dir: dir.Value.String(), Profile: profile.Value.String(), Err: err}
	}

	return &Config{*config}, nil
}

func (c *Config) GetGlobal() Global {
	return Global{
		Dir: c.GetString("dir"),
		Profile: c.GetString("profile"),
	}
}

type Config struct {
	viper.Viper
}

type ErrReadConfig struct {
	Dir string
	Profile string
	Err error
}

func (e ErrReadConfig) Error() string {
	return fmt.Sprintf("Failed to load %s config in %s: %v", e.Profile, e.Dir, e.Err)
}
