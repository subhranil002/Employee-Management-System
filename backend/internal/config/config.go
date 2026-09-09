package config

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv string

	HTTPAddr string

	AWSRegion string

	CognitoUserPoolID   string
	CognitoClientID     string
	CognitoClientSecret string

	EmployeeAPIURL string
	AllowedOrigin  string
}

func Load() (Config, error) {
	_ = godotenv.Load()

	cfg := Config{
		AppEnv:    os.Getenv("APP_ENV"),
		AWSRegion: os.Getenv("AWS_REGION"),
	}

	if cfg.AppEnv == "" {
		cfg.AppEnv = "development"
	}
	if cfg.AWSRegion == "" {
		cfg.AWSRegion = "ap-south-1"
	}

	secretARN := os.Getenv("BACKEND_SECRET_ARN")

	// AWS: load configuration from Secrets Manager
	if secretARN != "" {
		ctx := context.Background()
		awsCfg, err := awsconfig.LoadDefaultConfig(
			ctx,
			awsconfig.WithRegion(cfg.AWSRegion),
		)
		if err != nil {
			return Config{}, fmt.Errorf("load aws config for secrets manager: %w", err)
		}

		smClient := secretsmanager.NewFromConfig(awsCfg)
		result, err := smClient.GetSecretValue(
			ctx,
			&secretsmanager.GetSecretValueInput{
				SecretId: aws.String(secretARN),
			},
		)
		if err != nil {
			return Config{}, fmt.Errorf("get secret value: %w", err)
		}

		var secret map[string]string
		if err := json.Unmarshal([]byte(*result.SecretString), &secret); err != nil {
			return Config{}, fmt.Errorf("parse secret string: %w", err)
		}

		// Helper to fetch keys case-insensitively or with different common formats
		getSecret := func(keys ...string) string {
			for _, k := range keys {
				if val, ok := secret[k]; ok {
					return val
				}
			}
			return ""
		}

		cfg.CognitoUserPoolID = getSecret("COGNITO_USER_POOL_ID", "cognitoUserPoolId")
		cfg.CognitoClientID = getSecret("COGNITO_CLIENT_ID", "cognitoClientId")
		cfg.CognitoClientSecret = getSecret("COGNITO_CLIENT_SECRET", "cognitoClientSecret")
		cfg.EmployeeAPIURL = getSecret("EMPLOYEE_API_URL", "employeeApiUrl")
		cfg.AllowedOrigin = getSecret("ALLOWED_ORIGIN", "allowedOrigin")
		cfg.HTTPAddr = getSecret("HTTP_ADDR", "httpAddr")

	} else {
		// Local: load from environment
		cfg.CognitoUserPoolID = os.Getenv("COGNITO_USER_POOL_ID")
		cfg.CognitoClientID = os.Getenv("COGNITO_CLIENT_ID")
		cfg.CognitoClientSecret = os.Getenv("COGNITO_CLIENT_SECRET")
		cfg.EmployeeAPIURL = os.Getenv("EMPLOYEE_API_URL")
		cfg.AllowedOrigin = os.Getenv("ALLOWED_ORIGIN")
		cfg.HTTPAddr = os.Getenv("HTTP_ADDR")
	}

	// Apply defaults and validate
	if cfg.HTTPAddr == "" {
		cfg.HTTPAddr = ":3000"
	}
	if cfg.AllowedOrigin == "" {
		cfg.AllowedOrigin = "*"
	}

	if cfg.CognitoUserPoolID == "" {
		return Config{}, fmt.Errorf("COGNITO_USER_POOL_ID is required")
	}
	if cfg.CognitoClientID == "" {
		return Config{}, fmt.Errorf("COGNITO_CLIENT_ID is required")
	}
	if cfg.CognitoClientSecret == "" {
		return Config{}, fmt.Errorf("COGNITO_CLIENT_SECRET is required")
	}
	if cfg.EmployeeAPIURL == "" {
		return Config{}, fmt.Errorf("EMPLOYEE_API_URL is required")
	}

	return cfg, nil
}

func (c Config) CognitoIssuer() string {
	return fmt.Sprintf(
		"https://cognito-idp.%s.amazonaws.com/%s",
		c.AWSRegion,
		c.CognitoUserPoolID,
	)
}
