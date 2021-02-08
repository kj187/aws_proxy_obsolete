package aws

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/fatih/color"
	"go.kj187.de/aws_proxy/src/internal/configuration"
)

func loadSession(verbose bool) (*session.Session, error) {
	config := configuration.Get()
	sess, err := session.NewSessionWithOptions(session.Options{
		Profile:           config.Aws.Profile,
		SharedConfigState: session.SharedConfigEnable,
		Config: aws.Config{
			Region:                        aws.String(config.Aws.Region),
			CredentialsChainVerboseErrors: aws.Bool(verbose),
		},
	})
	if err != nil {
		return nil, err
	}
	return sess, nil
}

// GetBastionPublicIP exported
func GetBastionPublicIP(verbose bool) (string, error) {
	sess, err := loadSession(verbose)
	if err != nil {
		return "", err
	}

	config := configuration.Get()
	var filters []*ec2.Filter
	if len(config.Bastion.TagFilters) > 0 {
		if verbose {
			color.Blue(">> Try to find a bastion EC2 instance with tag filter:")
		}
		for _, filter := range config.Bastion.TagFilters {
			filters = append(filters, &ec2.Filter{
				Name: aws.String(filter.Name),
				Values: []*string{
					aws.String(filter.Value),
				},
			})
			if verbose {
				color.Blue(">>   Name: %v, Value: %v", filter.Name, filter.Value)
			}
		}
	}

	svc := ec2.New(sess)
	input := &ec2.DescribeInstancesInput{Filters: filters}

	result, err := svc.DescribeInstances(input)
	if err != nil {
		return "", err
	}

	if len(result.Reservations) <= 0 {
		return "", fmt.Errorf("No bastion instance found")
	}

	instance := result.Reservations[0].Instances[0]

	if verbose {
		if len(result.Reservations[0].Instances) > 1 {
			color.Blue(">> Found '%v' bastion instances. Using the first one.", len(result.Reservations[0].Instances))
		}
		color.Blue(">> Using bastion EC2 instance: %v (InstanceId: %v)", aws.StringValue(instance.PublicIpAddress), aws.StringValue(instance.InstanceId))
	}

	return aws.StringValue(instance.PublicIpAddress), nil
}

// RDSInstancesOutput exported
type RDSInstancesOutput struct {
	DBInstanceIdentifier string
	DBName               string
	Engine               string
	EngineVersion        string
	MasterUsername       string
	Host                 string
	Port                 int64
}

// GetRDSInstances exported
func GetRDSInstances(verbose bool) ([]RDSInstancesOutput, error) {
	if verbose {
		color.Blue("Try to find RDS instances")
	}

	sess, err := loadSession(verbose)
	if err != nil {
		return nil, err
	}

	svc := rds.New(sess)
	input := &rds.DescribeDBInstancesInput{}

	result, err := svc.DescribeDBInstances(input)
	if err != nil {
		return nil, err
	}

	var instances []RDSInstancesOutput
	for _, instance := range result.DBInstances {
		instances = append(instances, RDSInstancesOutput{
			DBInstanceIdentifier: *instance.DBInstanceIdentifier,
			DBName:               *instance.DBName,
			Engine:               *instance.Engine,
			EngineVersion:        *instance.EngineVersion,
			MasterUsername:       *instance.MasterUsername,
			Host:                 *instance.Endpoint.Address,
			Port:                 *instance.Endpoint.Port,
		})
	}

	if verbose {
		color.Blue(">> %v RDS instances found", len(instances))
	}

	return instances, nil
}
