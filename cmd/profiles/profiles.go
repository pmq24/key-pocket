package profiles

import (
	"github.com/spf13/cobra"

	"kp/cmd/profiles/create"
)

var Cmd = &cobra.Command{
	Use: "profiles",
	Aliases: []string{"p"},
	Short: "Manage profiles",
}

func init() {
	Cmd.AddCommand(create.Cmd)
}
