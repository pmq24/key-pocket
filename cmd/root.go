package cmd

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	Dir string
	Profile string
	Alg string
	Key string
)

func Execute() error {
  return rootCmd.Execute()
}

var rootCmd = &cobra.Command{
  Use:   "kp",
  Short: "Encrypt secret files and store them in the repo itself.",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		err := initViper(Dir, Profile)
		if err != nil {
			slog.Error(fmt.Sprintf("Failed load config: %v", err))
			os.Exit(1)
		}

		if cmd.Name() != "new" {
			err = readKey()
			if err != nil {
				slog.Error(fmt.Sprintf("Cannot read key: %v", err))
				os.Exit(1)
			}
		}
	},
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&Dir, "dir", "d", "./", "Set the directory the command will be executed in")
	rootCmd.PersistentFlags().StringVarP(&Profile, "profile", "p", "dev", "Set the profile")
}

func initViper(dir string, profile string) error {
	viper.AddConfigPath(dir)
	viper.SetConfigName("kp." + profile)
	viper.SetConfigType("yaml")
	
	return viper.ReadInConfig()
}

func readKey() error {
	keyFilePath := fmt.Sprintf("%s/kp_key.%s", Dir, Profile)
	err := godotenv.Load(keyFilePath)
	if err != nil {
		return err
	}

	Alg = os.Getenv(KeyAlg)
	if (Alg != "aes-256") {
		return fmt.Errorf("Invalid alg %s", Alg)
	}

	Key = os.Getenv(KeyVal)
	return nil
}
