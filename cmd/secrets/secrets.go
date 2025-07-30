package secrets

import (
	"github.com/spf13/cobra"

	"kp/cmd/secrets/decrypt"
	"kp/cmd/secrets/encrypt"
)

var Cmd = &cobra.Command{
	Use:     "secrets",
	Aliases: []string{"s"},
	Short:   "Manage secrets",
}

func init() {
	Cmd.AddCommand(encrypt.Cmd)
	Cmd.AddCommand(decrypt.Cmd)
}
