package cfg

import (
	"github.com/spf13/viper"
)

var (
	BaseCfgDir     = "./"
	BaseCfgProfile = "dev"
)

type BaseCfg struct {
	Viper *viper.Viper
}

func NewBaseCfg(opts NewCfgOpts) *BaseCfg {
	viper := viper.New()

	viper.Set("dir", opts.Dir)
	viper.Set("profile", opts.Profile)

	return &BaseCfg{viper}
}

func (c *BaseCfg) GetDir() string {
	return c.Viper.GetString("dir")
}

func (c *BaseCfg) GetProfile() string {
	return c.Viper.GetString("profile")
}
