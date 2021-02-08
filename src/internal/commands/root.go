package commands

import (
	"fmt"
	"os"
	"runtime"
	"strconv"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

var (
	// AWSProxyVersion ...
	AWSProxyVersion string

	verbose       bool
	version       bool
	sshPassphrase bool
)

// rootCommand registration
var rootCommand = &cobra.Command{
	Use: "aws_proxy",
	Run: func(cmd *cobra.Command, args []string) {},
}

// versionCommand registration
var versionCommand = &cobra.Command{
	Use:     "version",
	Short:   "Show current version if AWS Proxy",
	Example: "aws_proxy version",
	Run: func(cmd *cobra.Command, args []string) {
		color.Green("AWS Proxy v%s %s/%s\n", AWSProxyVersion, runtime.GOOS, runtime.GOARCH)
	},
}

var configurationCommand = &cobra.Command{
	Use:     "config",
	Short:   "Show current configuration",
	Example: "aws_proxy config",
	Run: func(cmd *cobra.Command, args []string) {
		c := viper.AllSettings()
		bs, err := yaml.Marshal(c)
		if err != nil {
			color.Red("Unable to marshal config to YAML: %v", err)
			os.Exit(1)
		}

		fmt.Printf("Config file: %s\n\n", viper.ConfigFileUsed())
		fmt.Println(string(bs))
	},
}

func init() {
	rootCommand.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
	rootCommand.PersistentFlags().BoolVarP(&sshPassphrase, "ssh-passphrase", "p", false, "enable ssh key passphrase")
	rootCommand.AddCommand(versionCommand)
	rootCommand.AddCommand(configurationCommand)
}

// Execute commands
func Execute() {
	if err := rootCommand.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "#1526981224 error: %s\n", err)
		rootCommand.Usage()
		os.Exit(1)
	}
}

func askForUserInput(question string) (int, error) {
	if question == "" {
		question = "Enter id (default 0):"
	}
	item := "0"
	green := color.New(color.FgGreen).SprintFunc()
	fmt.Printf("\n%v ", green(question))

	fmt.Scanln(&item)
	itemConverted, err := strconv.Atoi(item)
	if err != nil {
		return 0, err
	}

	return itemConverted, nil
}
