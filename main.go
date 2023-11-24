package main

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	app.Get("/secret", func(c *fiber.Ctx) error {
		secretName := "Encryptor-B2B"
		region := "ap-southeast-3"

		config, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
		if err != nil {
			log.Fatal(err)
		}

		// Create Secrets Manager client
		svc := secretsmanager.NewFromConfig(config)

		input := &secretsmanager.GetSecretValueInput{
			SecretId:     aws.String(secretName),
			VersionStage: aws.String("AWSCURRENT"), // VersionStage defaults to AWSCURRENT if unspecified
		}

		result, err := svc.GetSecretValue(context.TODO(), input)
		if err != nil {
			// For a list of exceptions thrown, see
			// https://docs.aws.amazon.com/secretsmanager/latest/apireference/API_GetSecretValue.html
			log.Fatal(err.Error())
		}

		// Decrypts secret using the associated KMS key.
		var secretString string = *result.SecretString

		// Your code goes here.
		fmt.Println(secretString)
		return c.SendString(secretString)
	})

	err := app.Listen(":80")
	if err != nil {
		panic(err)
	}
}

// func getSecret(secretName string) (string, error) {
// 	region := "ap-southeast-3" // Replace with your AWS region

// 	// Load AWS credentials and configuration
// 	config, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
// 	if err != nil {
// 		return "", err
// 	}

// 	// Create Secrets Manager client
// 	svc := secretsmanager.NewFromConfig(config)

// 	input := &secretsmanager.GetSecretValueInput{
// 		SecretId:     aws.String(secretName),
// 		VersionStage: aws.String("AWSCURRENT"), // VersionStage defaults to AWSCURRENT if unspecified
// 	}

// 	result, err := svc.GetSecretValue(context.TODO(), input)
// 	if err != nil {
// 		// For a list of exceptions thrown, see
// 		// https://docs.aws.amazon.com/secretsmanager/latest/apireference/API_GetSecretValue.html
// 		return "", err
// 	}

// 	// Decrypts secret using the associated KMS key.
// 	secretString := *result.SecretString

// 	return secretString, nil
// }
