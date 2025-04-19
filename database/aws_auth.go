package database

import (
	"TestHeroBackendGo/config"
	"context"
	"fmt"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/aws"
	aws_config "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/rds"
)

// GenerateAuthToken generates an IAM authentication token for RDS
func GenerateAuthToken(cfg *config.Config) (string, error) {
	if cfg.UseIAMAuth != "true" {
		return "", nil
	}

	// Load AWS configuration
	awsCfg, err := aws_config.LoadDefaultConfig(context.TODO(),
		aws_config.WithRegion(cfg.AWSRegion),
	)
	if err != nil {
		return "", fmt.Errorf("unable to load AWS SDK config: %v", err)
	}

	// Create RDS client
	rdsClient := rds.NewFromConfig(awsCfg)

	// Generate the authentication token
	port, err := strconv.ParseInt(cfg.DBPort, 10, 32)
	if err != nil {
		return "", fmt.Errorf("invalid port number: %v", err)
	}

	// Use the correct method for generating auth token
	authToken, err := rdsClient.(context.TODO(), &rds.GenerateDBAuthTokenInput{
		DBHostname: aws.String(cfg.DBHost),
		Port:       aws.Int32(int32(port)),
		DBUsername: aws.String(cfg.DBUser),
	})
	if err != nil {
		return "", fmt.Errorf("failed to generate auth token: %v", err)
	}

	return authToken, nil
}
