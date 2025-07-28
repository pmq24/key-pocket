package create

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/pmq24/kp/log"
	"github.com/pmq24/kp/cfg"
)

var Cmd = &cobra.Command{
  Use: "create <name>",
	Aliases: []string{"c"},
	Args: cobra.ExactArgs(1),
  Short: "Create a new profile",
	Run: func(cmd *cobra.Command, args []string) {
		global := cfg.ReadGlobal()
		execute(global, args)
	},
}

func execute(global cfg.Global, args []string) {
	name := args[0]
	keyFilePath := filepath.Join(global.Dir, fmt.Sprintf(cfg.KeyFileFormat, name))
	log.Verbosef("Creating %s", keyFilePath)

	_, err := os.Stat(keyFilePath)
	keyFileAlreayExists :=err == nil
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
