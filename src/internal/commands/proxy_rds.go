package commands

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"go.kj187.de/aws_proxy/src/internal/aws"
	"go.kj187.de/aws_proxy/src/internal/configuration"
	"go.kj187.de/aws_proxy/src/internal/proxy"
)

var rdsDBName string

func init() {
	proxyRdsCommand.Flags().StringVarP(&rdsDBName, "dbname", "n", "", "RDS database name")
	rootCommand.AddCommand(proxyRdsCommand)
}

// proxyRdsCommand registration
var proxyRdsCommand = &cobra.Command{
	Use:     "proxy:rds",
	Short:   "Open a proxy to an RDS instance",
	Example: "aws_proxy proxy:rds",
	Run: func(cmd *cobra.Command, args []string) {
		instances, err := aws.GetRDSInstances(verbose)
		if err != nil {
			color.Red("%v", err)
			return
		}

		config := configuration.Get()
		var instance aws.RDSInstancesOutput
		if rdsDBName != "" {
			if verbose {
				color.Blue("Try to get RDS instance with db name '%v'", rdsDBName)
			}
			for _, v := range instances {
				if v.DBName == rdsDBName {
					instance = v
					break
				}
			}
			if (aws.RDSInstancesOutput{}) == instance {
				color.Red("Instance not found with DBName '%v'", rdsDBName)
				return
			}
			if verbose {
				color.Blue(">> RDS instance with db name '%v' found (DBInstanceIdentifier: %v)", rdsDBName, instance.DBInstanceIdentifier)
			}
		} else {
			color.Green("Choose a RDS instance:")
			color.Green("")
			for i, instance := range instances {
				color.Green("[%v] %v (%v:%v)", i, instance.DBName, instance.Engine, instance.EngineVersion)
			}

			i, err := askForUserInput("")
			if err != nil {
				color.Red("%v", err)
				return
			}
			if i > len(instances) {
				color.Red("Instance not available with id '%v'", i)
				return
			}
			instance = instances[i]
			fmt.Println()
		}

		p := proxy.Proxy{
			Label:         "RDS",
			Verbose:       verbose,
			SSHPassphrase: sshPassphrase,
			Destination:   fmt.Sprintf("%v:%v", instance.Host, instance.Port),
			LocalPort:     config.Proxy.LocalPorts.RDS,
		}

		p.Start([]string{
			"Use something like:",
			fmt.Sprintf("$ mysql --host 127.0.0.1 --user %v --port %v -p %v", instance.MasterUsername, p.LocalPort, instance.DBName),
		})
	},
}
