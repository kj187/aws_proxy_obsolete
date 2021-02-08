package main

import (
	"strings"

	"github.com/fatih/color"
	"go.kj187.de/aws_proxy/src/internal/commands"
	"go.kj187.de/aws_proxy/src/internal/configuration"
)

// AWSProxyVersion ...
var AWSProxyVersion string = "1.0.0"

func main() {
	color.Green(`____ _ _ _ ____    ___  ____ ____ _  _ _   _`)
	color.Green(`|__| | | | [__     |__] |__/ |  |  \/   \_/`)
	color.Green(`|  | |_|_| ___]    |    |  \ |__| _/\_   |`)
	color.Green(``)
	color.Green(`The AWS Proxy enables you to reach AWS services which are protected due to a private network`)
	color.Green(strings.Repeat("_", 88))
	color.Green(``)

	configuration.Initialize()
	commands.AWSProxyVersion = AWSProxyVersion
	commands.Execute()

	color.Green(``)
}
