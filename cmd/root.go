package cmd

import (
	"github.com/spf13/cobra"

	"kp/cfg"
	"kp/cmd/profiles"
	"kp/cmd/secrets"
	"kp/log"
)

func Execute() error {
	return cmd.Execute()
}

var (
	cmd = &cobra.Command{
		Use:   "kp",
		Short: "Stop having to sync secrets between teammates by encrypting and storing them in the repo itself.",
	}
)

func init() {
	cmd.PersistentFlags().StringVarP(&cfg.BaseCfgDir, "dir", "d", "./", "Set the directory the command will be executed in")
	cmd.PersistentFlags().StringVarP(&cfg.BaseCfgProfile, "profile", "p", "dev", "Set the profile")

	cmd.PersistentFlags().BoolVarP(&log.Verbose, "verbose", "v", false, "Set verbose mode")

	cmd.AddCommand(profiles.Cmd)
	cmd.AddCommand(secrets.Cmd)
}
