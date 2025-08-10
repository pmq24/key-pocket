package create

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"kp/cfg"
	"kp/log"
)

var (
	Config *cfg.BaseCfg
	Cmd    = &cobra.Command{
		Use:     "create <name>",
		Aliases: []string{"c"},
		Args:    cobra.ExactArgs(1),
		Short:   "Create a new profile",
		Run: func(cmd *cobra.Command, args []string) {
			pflags := cmd.Root().PersistentFlags()

			baseCfg := cfg.NewBaseCfg(cfg.NewCfgOpts{
				Dir:     pflags.Lookup("dir").Value.String(),
				Profile: pflags.Lookup("profile").Value.String(),
			})

			execute(baseCfg, args)
		},
	}
)

func execute(config *cfg.BaseCfg, args []string) {
	name := args[0]
	keyFilePath := filepath.Join(config.GetDir(), fmt.Sprintf(cfg.KeyFileFormat, name))

	log.Verbosef("Creating %s", keyFilePath)

	_, err := os.Stat(keyFilePath)
	keyFileAlreayExists := err == nil
	if keyFileAlreayExists {
		log.Errorln("Key already exists")
		return
	}

	key := make([]byte, 32)
	rand.Read(key)
	encodedKey := base64.StdEncoding.EncodeToString(key)

	err = os.WriteFile(keyFilePath, []byte(encodedKey), 0600)
	if err != nil {
		log.Errorf("Failed to write key file: %v", err)
		return
	}

	log.Infoln("Created profile successfully")
}
