package cfg

import (
	"fmt"

	"github.com/spf13/viper"
)

type AppCfg struct {
	BaseCfg
}

func NewAppCfg(opts NewCfgOpts) (*AppCfg, error) {
	base := NewBaseCfg(opts)

	viper := viper.New()
	viper.SetConfigName(fmt.Sprintf(CfgFileFormat, base.GetProfile()))
	viper.SetConfigType("yml")
	viper.AddConfigPath(base.GetDir())

	err := viper.ReadInConfig()
	if err != nil {
		return nil, &ErrNewAppCfg{Err: err}
	}

	viper.Set("dir", opts.Dir)
	viper.Set("profile", opts.Profile)

	return &AppCfg{BaseCfg: BaseCfg{viper}}, nil
}

type ErrNewAppCfg struct {
	Err error
}

func (e *ErrNewAppCfg) Error() string {
	return fmt.Sprintf("Cannot create AppCfg: %v", e.Err)
}

func (c *AppCfg) GetSecrets() []string {
	return c.Viper.GetStringSlice("secrets")
}
