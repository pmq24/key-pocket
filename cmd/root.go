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
	dir string
	profile string
	cmd = &cobra.Command{
		Use:   "kp",
		Short: "Stop having to syncing secrets between teammates by encrypting and storing them in the repo itself.",
	}
)

func init() {
	cmd.PersistentFlags().StringVarP(&dir, "dir", "d", "./", "Set the directory the command will be executed in")
	cmd.PersistentFlags().StringVarP(&profile, "profile", "p", "dev", "Set the profile")
	cmd.PersistentFlags().BoolVarP(&log.Verbose, "verbose", "v", false, "Set verbose mode")

	cfg.SetGlobal(cfg.SetGlobalOpts{
    Dir: cmd.PersistentFlags().Lookup("dir"),
    Profile: cmd.PersistentFlags().Lookup("profile"),
  })

	cmd.AddCommand(profiles.Cmd)
	cmd.AddCommand(secrets.Cmd)
}
