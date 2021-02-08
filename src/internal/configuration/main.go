package configuration

import (
	"fmt"
	"log"
	"os"

	"github.com/fatih/color"
	"github.com/spf13/viper"
)

// Config exported
type Config struct {
	Aws struct {
		Region  string
		Profile string
	}
	Bastion struct {
		PublicIP string `mapstructure:"public_ip"`
		SSH      struct {
			IdentityFile string `mapstructure:"identity_file"`
			User         string `mapstructure:"user"`
			Port         int    `mapstructure:"port"`
		}
		TagFilters []struct {
			Name  string
			Value string
		} `mapstructure:"tag_filters"`
	}
	Proxy struct {
		LocalPorts struct {
			RDS int64
		} `mapstructure:"local_ports"`
	}
}

var configuration Config

// Initialize configuration
func Initialize() {
	viper.SetConfigName(".aws_proxy")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("$HOME/")
	viper.AutomaticEnv()

	viper.SetDefault("bastion.tag_filters", [1]map[string]string{{"name": "tag:Name", "value": "bastion"}})
	viper.SetDefault("bastion.ssh.user", "ec2-user")
	viper.SetDefault("bastion.ssh.port", 22)
	viper.SetDefault("bastion.ssh.identity_file", "~/.ssh/rd_rsa")
	viper.SetDefault("aws.region", "eu-central-1")
	viper.SetDefault("aws.profile", "")
	viper.SetDefault("proxy.local_ports.rds", 6546)

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			log.Fatalf("Error reading config file, %s", err)
		}
		viper.SafeWriteConfig()
		color.Yellow("Could not find config file '~/.aws_proxy'. Will create it now with default values. Check 'aws_proxy config'")
		fmt.Println()
		Initialize()
		return
	}

	err := viper.Unmarshal(&configuration)
	if err != nil {
		color.Red("%v", err)
		os.Exit(1)
	}
}

// Get configuration
func Get() Config {
	return configuration
}
