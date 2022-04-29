package cmd

import (
	"path"
	"path/filepath"
	"strings"

	"github.com/Clevermind-Think-Mint/buymint-cli-go/internal/logger"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "buymint-cli",
	Short: "buymint-cli - CLI to access and use BuyMint API",
	Long:  `buymint-cli - command line interface to access and use BuyMint API`,
}

// Execute adds all child commands to the root command and sets flags appropriately. This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute(v string, bh string, bd string) {
	rootCmd.Version = v + " (Build: " + bd + ")"
	// Executing...
	if err := rootCmd.Execute(); err != nil {
		logger.Fatal(errors.Wrap(err, "Failed to execute command").Error())
	}
}

func init() {
	rootCmd.PersistentFlags().StringP("config", "c", "config.json", "Configuration file to use")
	viper.BindPFlag("config", rootCmd.PersistentFlags().Lookup("config"))
	rootCmd.PersistentFlags().Bool("pretty", false, "Enable or disable human friendly logs (Pretty but inefficient)")
	viper.BindPFlag("pretty", rootCmd.PersistentFlags().Lookup("pretty"))
	rootCmd.PersistentFlags().Bool("debug", false, "Enable or disable log debug level")
	viper.BindPFlag("debug", rootCmd.PersistentFlags().Lookup("debug"))
	rootCmd.PersistentFlags().Bool("info", false, "Enable or disable log info level")
	viper.BindPFlag("info", rootCmd.PersistentFlags().Lookup("info"))
	rootCmd.PersistentFlags().Bool("warn", false, "Enable or disable log warn level")
	viper.BindPFlag("warn", rootCmd.PersistentFlags().Lookup("warn"))
	rootCmd.PersistentFlags().Bool("error", false, "Enable or disable log error level")
	viper.BindPFlag("error", rootCmd.PersistentFlags().Lookup("error"))
	rootCmd.PersistentFlags().Bool("self-signed", false, `Use this option if you wish to contact BuyMint API with self-signed certificate`)
	viper.BindPFlag("self-signed", rootCmd.PersistentFlags().Lookup("self-signed"))
	rootCmd.PersistentFlags().StringP("token", "t", "", `Authentication token to contact BuyMint API`)
	viper.BindPFlag("token", rootCmd.PersistentFlags().Lookup("token"))

	cobra.OnInitialize(func() {
		// Reading custom config file (if set) and merging content with default config file content
		fileConfig := viper.GetString("config")
		if fileConfig != "config.json" {
			fileExtension := filepath.Ext(fileConfig)
			fileName := strings.TrimSuffix(path.Base(fileConfig), fileExtension)
			fileDir := filepath.Dir(fileConfig)
			viper.SetConfigName(fileName)          // Configuring config file name
			viper.SetConfigType(fileExtension[1:]) // Configuring config file type without the "dot"
			viper.AddConfigPath(fileDir)
			// Ignoring if config file is not present otherwise signaling the real error
			if err := viper.MergeInConfig(); err != nil {
				// Config file was found but another error was produced
				logger.Fatal(errors.Wrap(err, "Fatal error while reading config file").Error())
			}
		}
		// Setting correct log level (assuming by default all the values are false)
		logLevel := logger.PanicLevel
		if viper.GetBool("debug") {
			logLevel = logger.DebugLevel
		} else if viper.GetBool("info") {
			logLevel = logger.InfoLevel
		} else if viper.GetBool("warn") {
			logLevel = logger.WarnLevel
		} else if viper.GetBool("error") {
			logLevel = logger.ErrorLevel
		}
		logger.LogInit(logLevel, viper.GetBool("pretty"))
	})

	// Reading default config file (if exists)
	fileConfig := viper.GetString("config")
	fileExtension := filepath.Ext(fileConfig)
	fileName := strings.TrimSuffix(path.Base(fileConfig), fileExtension)
	fileDir := filepath.Dir(fileConfig)
	viper.SetConfigName(fileName)          // Configuring config file name
	viper.SetConfigType(fileExtension[1:]) // Configuring config file type without the "dot"
	viper.AddConfigPath(fileDir)
	// Ignoring if config file is not present otherwise signaling the real error
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			// Config file was found but another error was produced
			logger.Fatal(errors.Wrap(err, "Fatal error while reading config file").Error())
		}
	}
}
