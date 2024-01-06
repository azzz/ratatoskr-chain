package cmd

import (
	"errors"
	"log"
	"os"
	"path"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "ratatoskr",
	Short: "A home-made blockchain node.",
	Run:   func(cmd *cobra.Command, args []string) {},
}

var (
	workdir string
	config  Config
	logger  *log.Logger
)

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initWorkdir, initConfig)

	logger = log.Default()

	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.PersistentFlags().StringVar(&workdir, "workdir", "$HOME/.ratatoskr", "node working directory")
}

func initWorkdir() {
	workdir = os.ExpandEnv(workdir)

	logger.Printf("Use workdir: %s", workdir)

	// create workdir
	if _, err := os.Stat(workdir); errors.Is(err, os.ErrNotExist) {
		cobra.CheckErr(os.MkdirAll(workdir, os.ModePerm))
	} else {
		cobra.CheckErr(err)
	}

	// create datadir
	p := path.Join(workdir, "data")
	if _, err := os.Stat(p); errors.Is(err, os.ErrNotExist) {
		cobra.CheckErr(os.MkdirAll(p, os.ModePerm))
	} else {
		cobra.CheckErr(err)
	}

	viper.Set("data-dir", p)
}

func initConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath(workdir)

	// write default config file
	if err := viper.ReadInConfig(); err != nil {
		cobra.CheckErr(viper.SafeWriteConfig())
	}

	viper.Unmarshal(&config)
}
