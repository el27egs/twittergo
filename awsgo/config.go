package awsgo

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/starlingapps/twittergo/models"
	"log"
)

var (
	Cfg aws.Config
	err error
)

func LoadDefaultConfig() {
	Cfg, err = config.LoadDefaultConfig(context.Background(), config.WithRegion("us-east-1"))
	if err != nil {
		log.Fatalf("Unable to load SDK awsgo, %v", err)
	}
}

func GetSecret(secretId string) (*models.Settings, error) {
	fmt.Printf("> Pido Secreto %s\n", secretId)
	var settings models.Settings
	svc := secretsmanager.NewFromConfig(Cfg)

	data, err := svc.GetSecretValue(context.Background(), &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(secretId),
	})
	if err != nil {
		return nil, fmt.Errorf("GetSecret error %s", err)
	}
	json.Unmarshal([]byte(*data.SecretString), &settings)
	fmt.Printf("> Lectura de secreto OK %s\n", secretId)
	return &settings, nil
}
